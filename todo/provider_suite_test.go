package todo_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phux/pomogoro/todo"
)

func TestTodoProvider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TodoProvider Suite")
}

var _ = Describe("TodoTxt Provider", func() {
	Context("Given an empty file", func() {

		It("should return an empty todo list", func() {
			tmpFile := newTmpFile()
			defer os.Remove(tmpFile.Name())
			provider := todo.NewTodoTxtProvider(tmpFile)

			todos, _ := provider.Load()
			Expect(todos).To(BeEmpty())
		})
	})

	Context("Given an file with a single todo", func() {
		It("should return 1 todo", func() {
			tmpFile := newTmpFile()
			defer os.Remove(tmpFile.Name())

			text := []byte("plain todo without formatting")
			err := ioutil.WriteFile(tmpFile.Name(), text, 0)
			if err != nil {
				log.Fatal("Failed to write to temporary file", err)
			}

			provider := todo.NewTodoTxtProvider(tmpFile)

			todos, _ := provider.Load()
			Expect(todos).To(HaveLen(1))
		})
	})

	Context("Given an file with two todos", func() {
		It("should return 2 todos", func() {
			tmpFile := newTmpFile()
			defer tmpFile.Close()

			text := []byte("plain todo without formatting\nsecond todo without formatting")

			err := ioutil.WriteFile(tmpFile.Name(), text, 0)
			if err != nil {
				log.Fatal("Failed to write to temporary file", err)
			}

			provider := todo.NewTodoTxtProvider(tmpFile)

			todos, err := provider.Load()
			Expect(err).To(BeNil())
			Expect(todos).To(HaveLen(2))
		})
	})
})

func newTmpFile() *os.File {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "todotxt-provider-*.txt")
	if err != nil {
		log.Fatal("cannot create temp file", err)
	}
	return tmpFile
}
