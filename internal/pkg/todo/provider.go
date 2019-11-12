package todo

import (
	"bufio"
	"io"
	"log"
)

type TodoTxtProvider struct {
	reader io.Reader
}

func NewTodoTxtProvider(reader io.Reader) TodoTxtProvider {
	return TodoTxtProvider{reader: reader}
}

func (tp TodoTxtProvider) Load() ([]*Todo, error) {
	todos := []*Todo{}

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

func (TodoTxtProvider) parse(line string) (*Todo, error) {
	return NewTodo(line), nil
}
