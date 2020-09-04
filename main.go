package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

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
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	hStore, err = NewHistoryStore(filepath.Join(home, ".wk"))
	if err != nil {
		return err
	}
	return app.Run(os.Args)
}

var hStore HistoryStore

var commands = []*cli.Command{
	startCommand,
	finishCommand,
	restCommand,
}
