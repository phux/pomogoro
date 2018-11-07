package history_test

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phux/pomogoro/config"
	"github.com/phux/pomogoro/history"
)

var _ = Describe("History", func() {
	var conf *config.Config

	BeforeEach(func() {
		conf = &config.Config{
			LogFile: setupLogFile(),
		}
	})

	Describe("Writing log files", func() {

		Context("On first write without existing file", func() {
			It("should create the file and write the first line", func() {
				err := os.Remove(conf.LogFile)
				if err != nil {
					log.Fatal(err)
				}

				history.Append("Some Todo", 10, conf)

				content, _ := ioutil.ReadFile(conf.LogFile)
				lines := strings.Split(string(content), "\n")
				Expect(len(lines)).To(Equal(2))
			})
		})

		Context("On writing to an existing file", func() {
			It("should just append the line", func() {
				history.Append("First Todo", 10, conf)
				history.Append("Second Todo", 10, conf)

				content, _ := ioutil.ReadFile(conf.LogFile)
				lines := strings.Split(string(content), "\n")
				Expect(len(lines)).To(Equal(3))
			})
		})

	})
})

func setupLogFile() string {
	tmp, err := ioutil.TempFile("", "pomogoro_test_")
	if err != nil {
		panic(err)
	}
	defer tmp.Close()

	return tmp.Name()
}
