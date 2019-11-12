package todo_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phux/pomogoro/internal/pkg/todo"
)

func TestTodo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Suite")
}

var _ = Describe("Todo", func() {
	var t *todo.Todo
	var title string

	Context("initially", func() {
		title = "some title"
		t = todo.NewTodo(title)

		It("has no duration tracked", func() {
			Expect(t.TrackedTime()).Should(BeZero())
		})
		It("has a title", func() {
			Expect(t.Title()).Should(Equal(title))
		})
	})

	Context("tracking 25 minutes spent", func() {
		t := todo.NewTodo(title)

		tracked := 25 * time.Minute
		t.Track(tracked)

		It("should increase the total duration to 25 minutes", func() {
			Expect(t.TrackedTime()).Should(Equal(tracked))
		})
	})

	Context("tracking 10 and 20 minutes spent", func() {
		t := todo.NewTodo(title)

		t.Track(10 * time.Minute)
		t.Track(20 * time.Minute)

		It("should increase the total duration to 30 minutes", func() {
			Expect(t.TrackedTime()).Should(Equal(30 * time.Minute))
		})
	})

	Context("adding a tracking with date", func() {
		t := todo.NewTodo(title)

		t.Track(10 * time.Minute)
		t.TrackWithDate(20*time.Minute, time.Now())

		It("should increase the total duration", func() {
			Expect(t.TrackedTime()).Should(Equal(30 * time.Minute))
		})
	})
})
