package todo

import (
	"bufio"
	"io"
	"log"
	"time"
)

type Provider interface {
	Load(time.Duration) []*Todo
}

type TodoTxtProvider struct {
	reader io.Reader
}

func NewTodoTxtProvider(reader io.Reader) TodoTxtProvider {
	return TodoTxtProvider{reader: reader}
}

func (tp TodoTxtProvider) Load(start time.Time) []*Todo {
	todos := make([]*Todo, 0)

	scanner := bufio.NewScanner(tp.reader)
	for scanner.Scan() {
		t, err := tp.parse(scanner.Text())
		if err != nil {
			log.Fatal("could not parse todo from file", err)
		} else {
			todos = append(todos, t)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return todos
}

func (tp TodoTxtProvider) parse(line string) (*Todo, error) {
	return NewTodo(line), nil
}
