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

	It("aggregates results from all engines", func() {
		aggregator := Aggregator{
			Engines: Engines{
				"Photos": Replicas{fakeEngine("Photos")},
				"Videos": Replicas{fakeEngine("Videos")},
			}}
		results := aggregator.Search("golang")
		Expect(results).To(Equal(Results{
			"Photos": "Photos result for golang",
			"Videos": "Videos result for golang",
		}))
	})

	It("aggregates results concurrently", func() {
		aggregator := Aggregator{
			Engines: Engines{
				"Photos": Replicas{fakeSlowEngine("Photos", time.Millisecond)},
				"Videos": Replicas{fakeSlowEngine("Videos", time.Millisecond)},
			}}
		start := time.Now()
		aggregator.Search("golang")
		Expect(time.Since(start)).To(BeNumerically("<", 2*time.Millisecond))
	})

	It("has timeout and aggregates only ready results", func() {
		aggregator := Aggregator{
			Engines: Engines{
				"Photos": Replicas{fakeSlowEngine("Photos", 3*time.Millisecond)},
				"Videos": Replicas{fakeSlowEngine("Videos", 1*time.Millisecond)},
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
			Engines: Engines{
				"Photos": Replicas{
					fakeSlowEngine("Photos1", 3*time.Millisecond),
					fakeSlowEngine("Photos2", 2*time.Millisecond),
					fakeSlowEngine("Photos3", 1*time.Millisecond),
				},
				"Videos": Replicas{
					fakeSlowEngine("Videos1", 1*time.Millisecond),
					fakeSlowEngine("Videos2", 3*time.Millisecond),
				},
			}}
		start := time.Now()
		results := aggregator.Search("golang")
		Expect(time.Since(start)).To(BeNumerically("<", 2*time.Millisecond))
		Expect(results).To(Equal(Results{
			"Photos": "Photos3 result for golang",
			"Videos": "Videos1 result for golang",
		}))
	})

	It("doesn't leak goroutines", func() {
		aggregator := Aggregator{
			Engines: Engines{
				"Photos": Replicas{
					fakeSlowEngine("Photos1", 2*time.Millisecond),
					fakeSlowEngine("Photos2", 1*time.Millisecond),
				},
				"Videos": Replicas{
					fakeSlowEngine("Videos1", 4*time.Millisecond),
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

func fakeEngine(name string) Engine {
	return func(query Query) Result {
		return Result(fmt.Sprint(name, " result for ", query))
	}
}

func fakeSlowEngine(name string, duration time.Duration) Engine {
	return func(query Query) Result {
		time.Sleep(duration)
		return fakeEngine(name)(query)
	}
}
