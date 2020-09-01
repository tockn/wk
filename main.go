package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := &cli.App{
		Name:     "wk",
		Usage:    "超素朴打刻ツール",
		Commands: commands,
	}
	return app.Run(os.Args)
}

var store HistoryStore

var commands = []*cli.Command{
	startCommand,
	{
		Name:    "finish",
		Aliases: []string{"f"},
		Usage:   "仕事おわり",
		Action:  finishWorking,
	},
}

func finishWorking(c *cli.Context) error {
	return nil
}
