package main

import (
	"fmt"

	"github.com/hellflame/argparse"

	"github.com/myOmikron/bbb-rec-controller/server"
)

func main() {
	parser := argparse.NewParser(
		"bbb-rec-controller",
		"",
		&argparse.ParserConfig{
			DisableDefaultShowHelp: true,
		},
	)

	configPath := parser.String("", "config-path", &argparse.Option{Default: "/etc/bbb-rec-controller/config.toml"})

	if err := parser.Parse(nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	server.StartServer(*configPath)
}
