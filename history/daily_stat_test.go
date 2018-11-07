package history

import (
	"fmt"
	"testing"

	"github.com/stvp/assert"
)

func Test_add_addsFirstStat(t *testing.T) {
	date := "2018-01-01"
	dailyStat := NewDailyStat(date)
	todo := "foo"
	stat := &PomodoroStat{todo: todo, date: date}

	dailyStat.add(stat)

	assert.Equal(
		t,
		stat,
		dailyStat.todos[todo][0],
		fmt.Sprintf("Did not find expected %v in %v", stat, dailyStat),
	)
}
func Test_add_appendsToExistingStats(t *testing.T) {
	date := "2018-01-01"
	dailyStat := NewDailyStat(date)
	todo := "foo"
	stat1 := &PomodoroStat{todo: todo, date: date}
	stat2 := &PomodoroStat{todo: todo, date: date}
	dailyStat.add(stat1)
	dailyStat.add(stat2)

	if len(dailyStat.todos[todo]) != 2 {
		t.Fatalf("dailyStat does not contain 2 stats for '%s', contains %v", todo, dailyStat.todos[todo])
	}
	assert.Equal(
		t,
		stat1,
		dailyStat.todos[todo][0],
	)
	assert.Equal(
		t,
		stat2,
		dailyStat.todos[todo][1],
	)
}

func Test_aggregateDurationPerDay(t *testing.T) {
	todo := "my todo"
	date := "2018-10-01"

	var data = []struct {
		name             string
		stats            []*PomodoroStat
		expectedDuration int
	}{
		{
			"single stat",
			[]*PomodoroStat{
				{todo: todo, date: date, duration: 25},
			},
			25,
		},
		{
			"multiple stats",
			[]*PomodoroStat{
				{todo: todo, date: date, duration: 25},
				{todo: todo, date: date, duration: 15},
			},
			40,
		},
		{
			"stats from different todos",
			[]*PomodoroStat{
				{todo: todo, date: date, duration: 25},
				{todo: "another todo", date: date, duration: 25},
			},
			25,
		},
	}

	for _, tt := range data {
		dailyStat := NewDailyStat(date)
		for _, stat := range tt.stats {
			dailyStat.add(stat)
		}

		aggregated := dailyStat.aggregateDuration(todo)

		assert.Equal(t, tt.expectedDuration, aggregated)
	}
}
