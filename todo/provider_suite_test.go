package todo_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

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
		tmpFile, err := ioutil.TempFile(os.TempDir(), "todotxt-provider-")
		if err != nil {
			log.Fatal("cannot create temp file", err)
		}
		defer os.Remove(tmpFile.Name())

		provider := todo.NewTodoTxtProvider(tmpFile)

		It("should return an empty todo list", func() {
			todos := provider.Load(time.Now())
			Expect(todos).To(BeEmpty())
		})
	})

	Context("Given an file with a single todo", func() {
		tmpFile, err := ioutil.TempFile(os.TempDir(), "todotxt-provider-")
		if err != nil {
			log.Fatal("cannot create temp file", err)
		}
		defer os.Remove(tmpFile.Name())

		text := []byte("plain todo without formatting")
		if _, err = tmpFile.Write(text); err != nil {
			log.Fatal("Failed to write to temporary file", err)
		}

		provider := todo.NewTodoTxtProvider(tmpFile)
		todos := provider.Load(time.Now())

		It("should return 1 todo", func() {
			Expect(todos).To(HaveLen(1))
		})
	})
})
