package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all settings
type Config struct {
	TodoFile         string
	LogFile          string
	PomodoroDuration int
	BreakDuration    int
	PomodoroEnabled  bool
	LogIdleTime      bool
	LogBreakTime     bool
}

// Load builds Config struct from given confPath
func Load(confPath string) *Config {
	var conf Config

	tomlData, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println("could not read config path at " + confPath)
		panic(err)
	}
	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		panic(err)
	}

	return &conf
}

// GetPomodoroDuration returns configured pomodoro duration
func (conf *Config) GetPomodoroDuration() time.Duration {
	return time.Duration(conf.PomodoroDuration) * time.Minute
}

// GetBreakDuration returns configured pomodoro duration
func (conf *Config) GetBreakDuration() time.Duration {
	return time.Duration(conf.BreakDuration) * time.Minute
}
