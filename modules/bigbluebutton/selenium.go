package bigbluebutton

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

var (
	ErrInvalidMeeting    = errors.New("the meeting does not exist or is not flagged for recording")
	ErrInvalidContext    = errors.New("can't find the selenium driver in the context")
	ErrSeleniumGet       = errors.New("selenium couldn't get url")
	ErrSeleniumTimeout   = errors.New("selenium had a timeout waiting on something")
	ErrSeleniumClickFail = errors.New("selenium failed to click an existing element")
	ErrUnknownRecState   = errors.New("couldn't figure out in which state the recording is")
)

var (
	buttonStart  = "Start recording"
	buttonPause  = "Pause recording"
	buttonResume = "Resume recording"
)

func (bbb *BBB) IsRecordingRunning(wp worker.Pool, meetingId string) (bool, error) {
	result := false
	err := getRecordingButton(
		wp, bbb, meetingId,
		func(_ selenium.WebDriver) error {
			return nil
		},
		func(isRunning bool, _ selenium.WebElement, _ selenium.WebDriver) error {
			result = isRunning
			return nil
		},
	)
	return result, err
}

func (bbb *BBB) PauseRecording(wp worker.Pool, meetingId string) (bool, error) {
	wasRunning := false
	err := getRecordingButton(
		wp, bbb, meetingId,
		func(driver selenium.WebDriver) error {
			return waitForAndClick("//button[@aria-label='Close How would you like to join the audio?'][1]", driver, bbb.Logger)
		},
		func(isRunning bool, button selenium.WebElement, driver selenium.WebDriver) error {
			if isRunning {
				wasRunning = true
				err := button.Click()
				if err != nil {
					return ErrSeleniumClickFail
				}
				return waitForAndClick("//button[@aria-label='Yes'][1]", driver, bbb.Logger)
			} else {
				wasRunning = false
				return nil
			}
		},
	)
	return wasRunning, err
}

// getRecordingButton spawns a task which
// - joins the meeting
// - runs `prepare`
// - waits for the recording button
// - runs `then`
func getRecordingButton(
	wp worker.Pool,
	bbb *BBB,
	meetingId string,
	prepare func(selenium.WebDriver) error,
	then func(bool, selenium.WebElement, selenium.WebDriver) error,
) error {
	task := worker.NewTaskWithContext(func(ctx context.Context) error {
		// Get driver from context
		var driver selenium.WebDriver
		switch d := ctx.Value("selenium").(type) {
		case selenium.WebDriver:
			driver = d
		default:
			return ErrInvalidContext
		}

		// Pre-check meeting
		info, err := bbb.GetMeetingInfo(meetingId)
		if err != nil {
			return err
		}
		if *info.ReturnCode != "SUCCESS" {
			if *info.MessageKey != "notFound" {
				bbb.Logger.Errorf("getMeetingInfo returned: '%s'", *info.MessageKey)
			}
			return ErrInvalidMeeting
		}
		if !*info.Recording {
			return ErrInvalidMeeting
		}

		// Join meeting
		err = driver.Get(bbb.Join(bbb.Config.Username, meetingId))
		if err != nil {
			return ErrSeleniumGet
		}

		// Leave meeting when done
		defer func() {
			err := driver.Get(bbb.BaseUrl.String())
			bbb.Logger.Error(err)
		}()

		// Run prepare
		err = prepare(driver)
		if err != nil {
			return err
		}

		// Get the recording button
		var button selenium.WebElement
		err = driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
			b, err := driver.FindElement(
				selenium.ByXPATH,
				fmt.Sprintf("//div[@aria-label='%s' or @aria-label='%s' or @aria-label='%s'][1]", buttonStart, buttonPause, buttonResume),
			)
			if err != nil {
				return false, nil
			} else {
				button = b
				return true, nil
			}
		}, time.Second*10)
		if err != nil {
			// the above condition function doesn't produce an error, so the only possible one is a timeout
			return ErrSeleniumTimeout
		}

		// Get and match the label
		var isRunning bool
		label, err := button.GetAttribute("aria-label")
		if err != nil {
			bbb.Logger.Error(err)
			label = ""
		}
		switch label {
		case buttonStart:
			isRunning = false
		case buttonPause:
			isRunning = true
		case buttonResume:
			isRunning = false
		default:
			return ErrUnknownRecState
		}

		// Run then
		return then(isRunning, button, driver)
	})
	wp.AddTask(task)
	return task.WaitForResult()
}

// waitForAndClick waits for a button identified via XPath and clicks it once appeared.
func waitForAndClick(xpath string, driver selenium.WebDriver, logger echo.Logger) error {
	err := driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		// Try getting button
		button, err := driver.FindElement(selenium.ByXPATH, xpath)

		// Not found => keep on waiting
		if err != nil {
			return false, nil
		}

		// Try clicking the button
		if err = button.Click(); err != nil {
			return false, err
		}

		// Succeeded
		return true, nil
	}, time.Second*10)

	// Filter and log error
	if err != nil {
		if strings.HasPrefix(err.Error(), "timeout after ") {
			return ErrSeleniumTimeout
		} else {
			logger.Error(err)
			return ErrSeleniumClickFail
		}
	} else {
		return nil
	}
}
