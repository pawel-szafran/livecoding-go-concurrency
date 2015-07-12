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
		aggregator := Aggregator{
			Searches: Searches{
				"Photos": fakeSearch("Photos"),
				"Videos": fakeSearch("Videos"),
			},
			Timeout: time.Second,
		}
		results := aggregator.Search("golang")
		Expect(results).To(Equal(Results{
			"Photos": "Photos result for golang",
			"Videos": "Videos result for golang",
		}))
	})

	It("aggregates results concurrently", func() {
		aggregator := Aggregator{
			Searches: Searches{
				"Photos": fakeLongSearch("Photos", time.Millisecond),
				"Videos": fakeLongSearch("Videos", time.Millisecond),
			},
			Timeout: time.Second,
		}
		start := time.Now()
		aggregator.Search("golang")
		Expect(time.Since(start)).To(BeNumerically("<", 2*time.Millisecond))
	})

	It("has timeout and aggregates only ready results", func() {
		aggregator := Aggregator{
			Searches: Searches{
				"Photos": fakeLongSearch("Photos", 3*time.Millisecond),
				"Videos": fakeLongSearch("Videos", 1*time.Millisecond),
			},
			Timeout: 2 * time.Millisecond,
		}
		start := time.Now()
		results := aggregator.Search("golang")
		Expect(time.Since(start)).To(BeNumerically("<", 3*time.Millisecond))
		Expect(results).To(Equal(Results{
			"Videos": "Videos result for golang",
		}))
	})

	It("uses replicas to minimize tail latency impact", func() {
		aggregator := Aggregator{
			Searches: Searches{
				"Photos": Replicas{
					fakeLongSearch("Photos1", 3*time.Millisecond),
					fakeLongSearch("Photos2", 2*time.Millisecond),
					fakeLongSearch("Photos3", 1*time.Millisecond),
				},
				"Videos": Replicas{
					fakeLongSearch("Videos1", 1*time.Millisecond),
					fakeLongSearch("Videos2", 3*time.Millisecond),
				},
			},
			Timeout: time.Second,
		}
		start := time.Now()
		results := aggregator.Search("golang")
		Expect(time.Since(start)).To(BeNumerically("<", 2*time.Millisecond))
		Expect(results).To(Equal(Results{
			"Photos": "Photos3 result for golang",
			"Videos": "Videos1 result for golang",
		}))
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
