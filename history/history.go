package history

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/phux/pomogoro/config"
)

// History TODO: NEEDS COMMENT INFO
type History struct {
	stats map[string]*DailyStat
}

// PomodoroStat contains information of a finished pomodoro
type PomodoroStat struct {
	todo      string
	duration  int
	date      string
	startTime string
}

// Append logs to configured history file
func Append(todo string, duration int, conf *config.Config) {
	f := openLogFile(conf)
	defer f.Close()
	writer := bufio.NewWriter(f)
	fmt.Fprintln(writer, formatLine(duration, todo))
	err := writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// Aggregate the logs from file specified in conf
func Aggregate(conf *config.Config) (*History, error) {
	content, err := ioutil.ReadFile(conf.LogFile)
	if err != nil {
		return nil, err
	}

	return parseList(strings.Split(string(content), "\n")), nil
}

// GetStatsForLastDays aggregates the stats for last days
func (h *History) GetStatsForLastDays(days int) []*DailyStat {
	stats := []*DailyStat{}
	for i := 0; i <= days; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		if dailyStat, ok := h.stats[date]; ok {
			stats = append(stats, dailyStat)
		}
	}
	return stats
}

// NewHistory creates an instance of *History
func NewHistory() *History {
	return &History{stats: make(map[string]*DailyStat)}
}

func parseList(lines []string) *History {
	h := NewHistory()
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		stat, err := lineToStat(line)
		if err != nil {
			log.Fatal(err)
		} else {
			// TODO: really continue?
			h.addStat(stat)
		}
	}
	return h
}

func (h *History) addStat(stat *PomodoroStat) {
	if _, exists := h.stats[stat.date]; !exists {
		h.stats[stat.date] = NewDailyStat(stat.date)
	}
	h.stats[stat.date].add(stat)
}

func lineToStat(line string) (*PomodoroStat, error) {
	stat := &PomodoroStat{}
	parts := strings.Split(line, " | ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Line does not have 2 | separators: %s", line)
	}

	dateParts := strings.Split(parts[0], " ")
	stat.date = dateParts[0]
	if len(dateParts) == 2 {
		stat.startTime = dateParts[1]
	}

	stat.todo = parts[1]

	d, err := strconv.ParseInt(strings.Split(parts[2], " ")[0], 10, 64)
	if err != nil {
		return nil, err
	}
	stat.duration = int(d)

	return stat, nil
}

func formatLine(duration int, todo string) string {
	text := fmt.Sprintf(
		"%s | %s | %d minutes",
		time.Now().Format("2006-01-02 15:04:05"),
		todo,
		int(math.Ceil(float64(duration)/60.0)),
	)
	return text
}

func openLogFile(conf *config.Config) *os.File {
	if _, err := os.Stat(conf.LogFile); os.IsNotExist(err) {
		_, createErr := os.Create(conf.LogFile)
		if createErr != nil {
			fmt.Println("could not create log file at " + conf.LogFile)
			panic(createErr)
		}
	}

	f, err := os.OpenFile(conf.LogFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("could not open log file")
		panic(err)
	}
	return f
}
