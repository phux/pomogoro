package todo

import (
	"bufio"
	"io"
	"log"
)

type Provider interface {
	Load() []*Todo
}

type TodoTxtProvider struct {
	reader io.Reader
}

func NewTodoTxtProvider(reader io.Reader) TodoTxtProvider {
	return TodoTxtProvider{reader: reader}
}

func (tp TodoTxtProvider) Load() ([]*Todo, error) {
	todos := make([]*Todo, 0)

	scanner := bufio.NewScanner(tp.reader)
	for scanner.Scan() {
		todo, _ := tp.parse(scanner.Text())
		todos = append(todos, todo)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return todos, nil
}

func (tp TodoTxtProvider) parse(line string) (*Todo, error) {
	return NewTodo(line), nil
}
