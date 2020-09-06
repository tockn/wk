package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HistoryStore interface {
	SaveStartedAt(projectName string, date, startedAt time.Time) error
	SaveFinishedAt(projectName string, date, finishedAt time.Time) error
	SaveRestMin(projectName string, date time.Time, min int64) error
	FindHistory(projectName string) (History, error)
}

func timeToPtr(t time.Time) *time.Time {
	return &t
}

type WorkingTime struct {
	StartedAt  *time.Time
	FinishedAt *time.Time
	RestMin    int64
}

type History map[HistoryKey]WorkingTime
type HistoryKey string

func GetHistoryKey(t time.Time) HistoryKey {
	if t.Hour() < 6 {
		t = t.Add(-time.Hour * 24)
	}
	return HistoryKey(fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day()))
}

func NewHistoryStore(dir string) (HistoryStore, error) {
	h := &historyStore{
		histories: make(map[string]History, 0),
		dir:       dir,
	}
	return h, h.init()
}

type historyStore struct {
	// k = projectName
	histories map[string]History
	dir       string
}

func (s *historyStore) SaveStartedAt(projectName string, date, startedAt time.Time) error {
	c, ok := s.histories[projectName]
	if !ok {
		s.histories[projectName] = make(History, 0)
		c = s.histories[projectName]
	}
	h := c[GetHistoryKey(date)]
	h.StartedAt = timeToPtr(startedAt)
	c[GetHistoryKey(date)] = h
	s.histories[projectName] = c
	return s.writeCSV()
}

func (s *historyStore) SaveFinishedAt(projectName string, date, finishedAt time.Time) error {
	c, ok := s.histories[projectName]
	if !ok {
		s.histories[projectName] = make(History, 0)
		c = s.histories[projectName]
	}
	h := c[GetHistoryKey(date)]
	h.FinishedAt = timeToPtr(finishedAt)
	c[GetHistoryKey(date)] = h
	s.histories[projectName] = c
	return s.writeCSV()
}

var fileFormat = "wk-%s.csv"

func (s *historyStore) writeCSV() error {
	for p, h := range s.histories {
		path := filepath.Join(s.dir, fmt.Sprintf(fileFormat, p))
		os.Remove(path)
		f, err := os.OpenFile(filepath.Join(s.dir, fmt.Sprintf(fileFormat, p)), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		fmt.Fprint(f, "date,started_at,finished_at,rest_min")

		ks := make([]HistoryKey, 0, len(h))
		for k := range h {
			ks = append(ks, k)
		}
		sort.Slice(ks, func(i, j int) bool {
			var in, jn int
			for _, b := range []byte(ks[i]) {
				in += int(b)
			}
			for _, b := range []byte(ks[j]) {
				jn += int(b)
			}
			return in < jn
		})

		for _, k := range ks {
			w := h[k]
			var st, fi string
			if w.StartedAt != nil {
				st = w.StartedAt.Format(timeFormat)
			}
			if w.FinishedAt != nil {
				fi = w.FinishedAt.Format(timeFormat)
			}
			fmt.Fprintf(f, "\n%s,%s,%s,%d", k, st, fi, w.RestMin)
		}
		f.Close()
	}
	return nil
}

func (s *historyStore) SaveRestMin(projectName string, date time.Time, min int64) error {
	c, ok := s.histories[projectName]
	if !ok {
		s.histories[projectName] = make(History, 0)
		c = s.histories[projectName]
	}
	h := c[GetHistoryKey(date)]
	h.RestMin = min
	c[GetHistoryKey(date)] = h
	return s.writeCSV()
}

var timeFormat = "15:04:05"

func (s *historyStore) init() error {
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		if err := os.Mkdir(s.dir, 0755); err != nil {
			return err
		}
	}

	fs, err := ioutil.ReadDir(s.dir)
	if err != nil {
		return err
	}
	for _, info := range fs {
		f, err := os.OpenFile(filepath.Join(s.dir, info.Name()), os.O_RDONLY, 0755)
		if err != nil {
			return err
		}
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		project := strings.Split(strings.Split(info.Name(), "-")[1], ".")[0]
		hs := strings.Split(string(bs), "\n")

		for i, h := range hs {
			if i == 0 {
				s.histories[project] = make(History, 0)
				continue
			}
			if h == "" {
				continue
			}
			sp := strings.Split(h, ",")
			if len(sp) != 4 {
				return fmt.Errorf("invalid format %s: %s", project, h)
			}

			var st, fi time.Time
			if sp[1] != "" {
				st, err = time.Parse(timeFormat, sp[1])
				if err != nil {
					return err
				}
			}

			if sp[2] != "" {
				fi, err = time.Parse(timeFormat, sp[2])
				if err != nil {
					return err
				}
			}

			var rest int64
			if sp[3] != "" {
				rest, err = strconv.ParseInt(sp[3], 10, 64)
				if err != nil {
					return err
				}
			}
			s.histories[project][HistoryKey(sp[0])] = WorkingTime{
				StartedAt:  &st,
				FinishedAt: &fi,
				RestMin:    rest,
			}
		}
	}
	return nil
}

func (s *historyStore) FindHistory(projectName string) (History, error) {
	return s.histories[projectName], nil
}
