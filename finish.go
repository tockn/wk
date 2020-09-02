package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

var finishCommand = &cli.Command{
	Name:    "finish",
	Aliases: []string{"f"},
	Usage:   "仕事おわり",
	Action:  finishWorking,
	Flags:   finishFlags,
}

var finishFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "time",
		Aliases: []string{"t"},
		Value:   "",
		Usage:   "勤務終了時間",
	},
	&cli.StringFlag{
		Name:    "project",
		Aliases: []string{"p"},
		Value:   "default",
		Usage:   "project name",
	},
}

func finishWorking(c *cli.Context) error {
	t := time.Now()

	timeFlg := c.String("time")
	if timeFlg != "" {
		var err error
		t, err = parseTimeFlag(timeFlg)
		if err != nil {
			return err
		}
	}
	p := c.String("project")
	return hStore.SaveFinishedAt(p, time.Now(), t)
}
