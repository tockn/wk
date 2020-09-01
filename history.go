package main

import "time"

type HistoryStore interface {
	SaveStartedAt(label string, date, startedAt time.Time) error
	SaveFinishedAt(label string, date, finishedAt time.Time) error
}

func timeToPtr(t time.Time) *time.Time {
	return &t
}

type History struct {
	StartedAt  *time.Time
	FinishedAt *time.Time
}

var data = make(map[string]*History, 0)
