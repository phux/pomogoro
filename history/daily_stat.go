package history

import (
	"fmt"
)

// NewDailyStat creates new instance of *DailyStat
func NewDailyStat(date string) *DailyStat {
	return &DailyStat{
		todos: make(map[string][]*PomodoroStat),
		date:  date,
	}
}

// DailyStat contains aggregated PomodoroStats per day
type DailyStat struct {
	todos map[string][]*PomodoroStat
	date  string
}

// GetDate returns date of DailyStat
func (d *DailyStat) GetDate() string {
	return d.date
}

func (d *DailyStat) add(stat *PomodoroStat) {
	if _, exists := d.todos[stat.todo]; !exists {
		d.todos[stat.todo] = []*PomodoroStat{stat}
	} else {
		d.todos[stat.todo] = append(d.todos[stat.todo], stat)
	}
}

func (d *DailyStat) aggregateDuration(todo string) int {
	duration := 0
	for _, v := range d.todos[todo] {
		duration += v.duration
	}
	return duration
}

// ToList returns a formatted list of finished pomodoros in d
func (d *DailyStat) ToList() []string {
	list := []string{}
	durations := make(map[string]int)
	for _, stats := range d.todos {
		for _, stat := range stats {
			list = append(list, fmt.Sprintf("%s %s | %s | %d minutes", stat.date, stat.startTime, stat.todo, stat.duration))
			durations[stat.todo] += stat.duration
		}

	}
	return list
}

// Summary aggregates all duraions per todo
func (d *DailyStat) Summary() []string {
	list := []string{}
	durations := make(map[string]int)
	for _, stats := range d.todos {
		for _, stat := range stats {
			durations[stat.todo] += stat.duration
		}

	}
	sum := 0
	for todo, duration := range durations {
		list = append(list, fmt.Sprintf("%s\n\tTotal: %d:%02d h\n", todo, duration/60, duration%60))
		sum += duration
	}

	list = append(list, fmt.Sprintf("Total: %d:%02d h\n\n", sum/60, sum%60))

	return list
}
