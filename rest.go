package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

var restCommand = &cli.Command{
	Name:    "rest",
	Aliases: []string{"r"},
	Usage:   "きゅうけい〜",
	Action:  rest,
	Flags:   restFlag,
}

var restFlag = []cli.Flag{
	&cli.StringFlag{
		Name:    "project",
		Aliases: []string{"p"},
		Value:   "default",
		Usage:   "project name",
	},
}

func rest(c *cli.Context) error {
	restArg := c.Args().Get(0)
	if restArg == "" {
		return fmt.Errorf("休憩時間を引数に渡してください")
	}
	rest, err := strconv.ParseInt(restArg, 10, 64)
	if err != nil {
		return err
	}
	p := c.String("project")
	return hStore.IncrementRestMin(p, time.Now(), rest)
}
