package main

import (
	"fmt"
	"time"
)

type HistoryStore interface {
	SaveStartedAt(projectName string, date, startedAt time.Time) error
	SaveFinishedAt(projectName string, date, finishedAt time.Time) error
}

func timeToPtr(t time.Time) *time.Time {
	return &t
}

type WorkingTime struct {
	StartedAt  *time.Time
	FinishedAt *time.Time
}

type History map[HistoryKey]*WorkingTime
type HistoryKey string

func GetHistoryKey(t time.Time) HistoryKey {
	return HistoryKey(fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day()))
}

func NewHistoryStore() HistoryStore {
	return &historyStore{
		histories: make(map[string]History, 0),
	}
}

type historyStore struct {
	// k = projectName
	histories map[string]History
}

func (h *historyStore) SaveStartedAt(projectName string, date, startedAt time.Time) error {
	c := h.histories[projectName]
	c[GetHistoryKey(date)].StartedAt = timeToPtr(startedAt)
	h.histories[projectName] = c
	return nil
}

func (h *historyStore) SaveFinishedAt(projectName string, date, finishedAt time.Time) error {
	c := h.histories[projectName]
	c[GetHistoryKey(date)].FinishedAt = timeToPtr(finishedAt)
	h.histories[projectName] = c
	return nil
}
