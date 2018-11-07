// Package todotxt provides functionality related to todo.txt lists
package todotxt

import (
	"io/ioutil"
	"strings"

	"github.com/phux/pomogoro/config"
)

func parseList(lines string) []string {
	return strings.Split(lines, "\n")
}

// ReadTodoTxt returns all todos as a string array
func ReadTodoTxt(conf *config.Config) ([]string, error) {
	content, err := ioutil.ReadFile(conf.TodoFile)
	if err != nil {
		return nil, err
	}
	list := parseList(string(content))
	return list, nil
}
