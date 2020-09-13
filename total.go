package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

var totalCommand = &cli.Command{
	Name:    "total",
	Aliases: []string{"t"},
	Usage:   "合計勤務時間（h）",
	Action:  totalWorkingTime,
	Flags:   totalFlags,
}

var totalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "project",
		Aliases: []string{"p"},
		Value:   "default",
		Usage:   "project name",
	},
}

func totalWorkingTime(c *cli.Context) error {
	unitStr := c.Args().Get(0)
	unit := 1.0
	if unitStr != "" {
		var err error
		unit, err = strconv.ParseFloat(unitStr, 10)
		if err != nil {
			return err
		}
	}
	p := c.String("project")
	h, err := hStore.FindHistory(p)
	if err != nil {
		return err
	}
	total := 0.0
	for _, k := range h.SortedKey() {
		w := h[k]
		st := *w.StartedAt
		fi := *w.FinishedAt
		if w.StartedAt.Hour() < 6 {
			st = w.StartedAt.Add(time.Hour * 24)
		}
		if w.FinishedAt.Hour() < 6 {
			fi = w.FinishedAt.Add(time.Hour * 24)
		}
		total += float64(fi.Unix()-st.Unix()) / 60 / 60
		total -= float64(w.RestMin) / 60
		fmt.Printf("%s: %.2f\n", k, float64(fi.Unix()-st.Unix())/60/60-float64(w.RestMin)/60)
	}
	fmt.Printf("Total Working Time: %.2f\n", total)
	if unit != 1 {
		fmt.Printf("Salary: %.2f\n", total*unit)
	}
	return nil
}
