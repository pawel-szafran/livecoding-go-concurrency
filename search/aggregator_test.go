package search_test

import (
	"fmt"
	"runtime"
	"time"

	. "github.com/pawel-szafran/livecoding-go-concurrency/search"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregator", func() {

	It("aggregates results from all searches", func() {
		aggregator := Aggregator{
			Searches: Searches{
				"Photos": Replicas{fakeSearch("Photos")},
				"Videos": Replicas{fakeSearch("Videos")},
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
				"Photos": Replicas{fakeLongSearch("Photos", time.Millisecond)},
				"Videos": Replicas{fakeLongSearch("Videos", time.Millisecond)},
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
				"Photos": Replicas{fakeLongSearch("Photos", 3*time.Millisecond)},
				"Videos": Replicas{fakeLongSearch("Videos", 1*time.Millisecond)},
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

	It("leaks goroutines... see output :/", func() {
		aggregator := Aggregator{
			Searches: Searches{
				"Photos": Replicas{
					fakeLongSearch("Photos1", 2*time.Millisecond),
					fakeLongSearch("Photos2", 1*time.Millisecond),
				},
				"Videos": Replicas{
					fakeLongSearch("Videos1", 4*time.Millisecond),
				},
			},
			Timeout: 3 * time.Millisecond,
		}
		fmt.Println("Goroutines:", runtime.NumGoroutine())
		for i := 1; i <= 1000; i++ {
			aggregator.Search("golang")
			if i%10 == 0 {
				fmt.Println("Goroutines:", runtime.NumGoroutine())
			}
		}
	})

	It("doesn't leak goroutines", func() {
		aggregator := Aggregator{
			Searches: Searches{
				"Photos": Replicas{
					fakeLongSearch("Photos1", 2*time.Millisecond),
					fakeLongSearch("Photos2", 1*time.Millisecond),
				},
				"Videos": Replicas{
					fakeLongSearch("Videos1", 4*time.Millisecond),
				},
			},
			Timeout: 3 * time.Millisecond,
		}
		before := runtime.NumGoroutine()
		for i := 0; i < 20; i++ {
			aggregator.Search("golang")
		}
		time.Sleep(3 * time.Millisecond)
		left := runtime.NumGoroutine() - before
		Expect(left).To(BeNumerically("<", 5))
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
