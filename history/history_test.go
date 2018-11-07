package history

import (
	"regexp"
	"testing"

	"github.com/stvp/assert"
)

func Test_formatLine(t *testing.T) {
	lineRegex := `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \| `

	type args struct {
		duration int
		todo     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"plain todo",
			args{
				duration: 10,
				todo:     "plain todo",
			},
			lineRegex + `plain todo \| 1 minutes`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatLine(tt.args.duration, tt.args.todo)
			if !regexp.MustCompile(tt.want).MatchString(got) {
				t.Errorf("formatLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lineToStat(t *testing.T) {
	line := "2018-10-14 | My Todo | 20 minutes"

	stat, err := lineToStat(line)

	assert.Nil(t, err, "A valid todo should not return an error")
	assert.Equal(t, "2018-10-14", stat.date)
	assert.Equal(t, "My Todo", stat.todo)
	assert.Equal(t, 20, stat.duration)
}

func Test_lineToStat_InvalidDuration(t *testing.T) {
	line := "2018-10-14 | My Todo | ten minutes"

	stat, err := lineToStat(line)

	assert.Nil(t, stat, "Invalid duration should not return a *Stat")
	assert.NotNil(t, err, "Invalid duration should return an error")
}

func Test_parseList(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want *History
	}{
		{
			name: "empty history",
			args: args{
				lines: []string{},
			},
			want: &History{
				stats: map[string]*DailyStat{},
			},
		},
		{
			name: "single history entry",
			args: args{
				lines: []string{
					"2018-10-14 | My Todo | 20 minutes",
				},
			},
			want: &History{
				stats: map[string]*DailyStat{
					"2018-10-14": {
						todos: map[string][]*PomodoroStat{
							"My Todo": {
								{
									todo:     "My Todo",
									date:     "2018-10-14",
									duration: 20,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple history entries on same day",
			args: args{
				lines: []string{
					"2018-10-14 | My Todo | 20 minutes",
					"2018-10-14 | My Todo | 10 minutes",
				},
			},
			want: &History{
				stats: map[string]*DailyStat{
					"2018-10-14": {
						todos: map[string][]*PomodoroStat{
							"My Todo": {
								{
									todo:     "My Todo",
									date:     "2018-10-14",
									duration: 20,
								},
								{
									todo:     "My Todo",
									date:     "2018-10-14",
									duration: 10,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple history entries on different days",
			args: args{
				lines: []string{
					"2018-10-14 | My Todo | 20 minutes",
					"2018-10-15 | My Todo | 10 minutes",
				},
			},
			want: &History{
				stats: map[string]*DailyStat{
					"2018-10-14": {
						todos: map[string][]*PomodoroStat{
							"My Todo": {
								{
									todo:     "My Todo",
									date:     "2018-10-14",
									duration: 20,
								},
							},
						},
					},
					"2018-10-15": {
						todos: map[string][]*PomodoroStat{
							"My Todo": {
								{
									todo:     "My Todo",
									date:     "2018-10-15",
									duration: 10,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseList(tt.args.lines)
			if len(got.stats) != len(tt.want.stats) {
				t.Errorf("parseList() = %v, want %v", got, tt.want)
			}
			for k := range tt.want.stats {
				if _, ok := got.stats[k]; !ok {
					t.Errorf("parseList() missed key, want %s, got %v", k, tt.want.stats)
				}
			}
		})
	}
}
