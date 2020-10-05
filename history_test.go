package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var st = timeToPtr(time.Date(2020, 9, 1, 22, 0, 0, 0, time.UTC))
var fi = timeToPtr(time.Date(2020, 9, 1, 23, 0, 0, 0, time.UTC))
var testDir = ".wk_test"
var p = "test"

type csv struct {
	fileName string
	body     []byte
}

func init() {
	os.Remove(testDir)
	os.Mkdir(testDir, 0755)
}

func Test_historyStore_writeCSV(t *testing.T) {
	tests := []struct {
		name      string
		histories map[string]History
		wantErr   bool
		expect    []*csv
	}{
		{
			name: "projectが1件でデータが存在しない時データが消え、何も書き込まれない",
			histories: map[string]History{
				p: {},
			},
			expect: []*csv{},
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがあり、startedAtのみの場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						StartedAt:  st,
						FinishedAt: nil,
					},
				},
			},
			expect: []*csv{
				{
					fileName: "wk-2020-09.csv",
					body:     expectCSV1(),
				},
			},
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがあり、finishedAtのみの場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						StartedAt:  nil,
						FinishedAt: fi,
					},
				},
			},
			expect: []*csv{
				{
					fileName: "wk-2020-09.csv",
					body:     expectCSV2(),
				},
			},
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがあり、startedAt, finishedAtのみの場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						StartedAt:  st,
						FinishedAt: fi,
					},
				},
			},
			expect: []*csv{
				{
					fileName: "wk-2020-09.csv",
					body:     expectCSV3(),
				},
			},
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがあり、restMinのみの場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						RestMin: 100,
					},
				},
			},
			expect: []*csv{
				{
					fileName: "wk-2020-09.csv",
					body:     expectCSV4(),
				},
			},
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがあり、startedAt, restMinのみの場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						StartedAt: st,
						RestMin:   100,
					},
				},
			},
			expect: []*csv{
				{
					fileName: "wk-2020-09.csv",
					body:     expectCSV5(),
				},
			},
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがあり、finishedAt, restMinのみの場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						FinishedAt: fi,
						RestMin:    100,
					},
				},
			},
			expect: []*csv{
				{
					fileName: "wk-2020-09.csv",
					body:     expectCSV6(),
				},
			},
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがあり、startedAt, finishedAt, restMinがある場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						StartedAt:  st,
						FinishedAt: fi,
						RestMin:    100,
					},
				},
			},
			expect: []*csv{
				{
					fileName: "wk-2020-09.csv",
					body:     expectCSV7(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &historyStore{
				histories: tt.histories,
				dir:       testDir,
			}
			if err := s.writeCSV(); err != nil {
				t.Fatal(err)
			}
			got, err := getCSVs(p)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, got)
		})
	}
}

func getCSVs(project string) ([]*csv, error) {
	cs := make([]*csv, 0)
	fs, err := ioutil.ReadDir(filepath.Join(testDir, project))
	if err != nil {
		return nil, err
	}
	for _, info := range fs {
		f, err := os.OpenFile(filepath.Join(testDir, project, info.Name()), os.O_RDONLY, 0755)
		if err != nil {
			return nil, err
		}
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		cs = append(cs, &csv{fileName: info.Name(), body: bs})
	}
	return cs, nil
}

func expectCSV1() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,22:00:00,,0`)
}

func expectCSV2() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,,23:00:00,0`)
}

func expectCSV3() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,22:00:00,23:00:00,0`)
}

func expectCSV4() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,,,100`)
}

func expectCSV5() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,22:00:00,,100`)
}

func expectCSV6() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,,23:00:00,100`)
}

func expectCSV7() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,22:00:00,23:00:00,100`)
}
