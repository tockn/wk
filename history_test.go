package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var st = timeToPtr(time.Date(2020, 9, 1, 22, 0, 0, 0, time.UTC))
var fi = timeToPtr(time.Date(2020, 9, 1, 23, 0, 0, 0, time.UTC))
var p = "test"

func Test_historyStore_writeCSV(t *testing.T) {
	tests := []struct {
		name      string
		histories map[string]History
		dir       string
		wantErr   bool
		expect    []byte
	}{
		{
			name: "projectが1件でデータが存在しない時、ヘッダのみ書き込まれる",
			histories: map[string]History{
				p: {},
			},
			expect: expectCSV1(),
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
			expect: expectCSV2(),
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
			expect: expectCSV3(),
		},
		{
			name: "projectが1件でデータが存在せず、1件WorkingTimeがある場合",
			histories: map[string]History{
				p: {
					"2020-09-01": WorkingTime{
						StartedAt:  st,
						FinishedAt: fi,
					},
				},
			},
			expect: expectCSV4(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &historyStore{
				histories: tt.histories,
				dir:       tt.dir,
			}
			if err := s.writeCSV(); err != nil {
				t.Fatal(err)
			}
			got, err := getCSV(p)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, got)
		})
	}
}

func getCSV(project string) ([]byte, error) {
	f, err := os.OpenFile(fmt.Sprintf("wk-%s.csv", project), os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func expectCSV1() []byte {
	return []byte(`date,started_at,finished_at,rest_min`)
}

func expectCSV2() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,22:00:00,`)
}

func expectCSV3() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,,23:00:00,`)
}

func expectCSV4() []byte {
	return []byte(`date,started_at,finished_at,rest_min
2020-09-01,22:00:00,23:00:00,`)
}
