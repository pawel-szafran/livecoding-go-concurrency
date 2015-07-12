package search_test

import (
	"fmt"
	"time"

	. "github.com/pawel-szafran/livecoding-go-concurrency/search"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregator", func() {

	It("aggregates results from all searches", func() {
		aggregator := Aggregator{Searches{
			"Photos": fakeSearch("Photos"),
			"Videos": fakeSearch("Videos"),
		}}
		results := aggregator.Search("golang")
		Expect(results).To(Equal(Results{
			"Photos": "Photos result for golang",
			"Videos": "Videos result for golang",
		}))
	})

	It("aggregates results concurrently", func() {
		aggregator := Aggregator{Searches{
			"Photos": fakeLongSearch("Photos", time.Millisecond),
			"Videos": fakeLongSearch("Videos", time.Millisecond),
		}}
		start := time.Now()
		aggregator.Search("golang")
		Expect(time.Since(start)).To(BeNumerically("<", 2*time.Millisecond))
	})

})

func fakeSearch(name string) Search {
	return func(query Query) Result {
		return Result(fmt.Sprint(name, " result for ", query))
	}
}

func fakeLongSearch(name string, duration time.Duration) Search {
	return func(query Query) Result {
		time.Sleep(duration)
		return fakeSearch(name)(query)
	}
}
