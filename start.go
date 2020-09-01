package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var startCommand = &cli.Command{
	Name:    "start",
	Aliases: []string{"s"},
	Usage:   "仕事スタート",
	Action:  startWorking,
	Flags:   startFlags,
}

var startFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "time",
		Value: "",
		Usage: "勤務開始時間",
	},
}

func startWorking(c *cli.Context) error {
	t := time.Now()

	timeFlg := c.String("time")
	if timeFlg != "" {
		var err error
		t, err = parseTimeFlag(timeFlg)
		if err != nil {
			return err
		}
	}
	return store.SaveStartedAt(t)
}

var ErrInvalidTimeFormat = fmt.Errorf("時間指定のフォーマットが不正です")

func parseTimeFlag(s string) (time.Time, error) {
	now := time.Now()
	ss := strings.Split(s, ":")
	if len(ss) != 2 {
		return now, ErrInvalidTimeFormat
	}
	hour, err := strconv.Atoi(ss[0])
	if err != nil {
		return now, ErrInvalidTimeFormat
	}
	min, err := strconv.Atoi(ss[1])
	if err != nil {
		return now, ErrInvalidTimeFormat
	}
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location()), nil
}

/*
2020-09-01,9:30,14:50
*/