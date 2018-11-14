package config

import (
	"fmt"
	"io/ioutil"
	"os"
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

type ValidationError struct {
	s string
}

func (e *ValidationError) Error() string {
	return e.s
}

// Load builds Config struct from given confPath
func Load(confPath string) (*Config, error) {
	var conf Config

	tomlData, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println("could not read config path at " + confPath)
		panic(err)
	}
	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		panic(err)
	}

	err = validate(conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// GetPomodoroDuration returns configured pomodoro duration
func (conf *Config) GetPomodoroDuration() time.Duration {
	return time.Duration(conf.PomodoroDuration) * time.Minute
}

// GetBreakDuration returns configured pomodoro duration
func (conf *Config) GetBreakDuration() time.Duration {
	return time.Duration(conf.BreakDuration) * time.Minute
}

func validate(conf Config) error {
	if _, err := os.Stat(conf.TodoFile); os.IsNotExist(err) {
		return &ValidationError{s: fmt.Sprintf("specified TodoFile does not exist: '%s'", conf.TodoFile)}
	}

	if conf.LogFile == "" {
		return &ValidationError{
			s: "LogFile not specified (provide an absolute path, file will be created if not exists)",
		}
	}

	if conf.PomodoroEnabled {
		if conf.PomodoroDuration < 1 {
			return &ValidationError{
				s: "PomodoroEnabled is true, so PomodoroDuration must be specified as well (greater than 0)",
			}
		}
	}

	return nil
}
