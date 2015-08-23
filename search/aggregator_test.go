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
		enginesSyncer := NewSyncer()
		aggregator := Aggregator{
			Engines: Engines{
				"Photos": Replicas{fakeSyncedEngine("Photos", enginesSyncer)},
				"Videos": Replicas{fakeSyncedEngine("Videos", enginesSyncer)},
			}}
		go func() { enginesSyncer.WaitForAllReady().LetAllRun() }()
		results := make(chan Results)
		go func() { results <- aggregator.Search("golang") }()
		Eventually(results).Should(Receive())
	})

	It("has timeout and aggregates only ready results", func() {
		slowEnginesSyncer := NewSyncer()
		aggregator := Aggregator{
			Engines: Engines{
				"Photos": Replicas{fakeSyncedEngine("Photos", slowEnginesSyncer)},
				"Videos": Replicas{fakeEngine("Videos")},
			},
			Timeout: 10 * time.Millisecond,
		}
		results := make(chan Results)
		go func() { results <- aggregator.Search("golang") }()
		Eventually(results).Should(Receive(Equal(Results{
			"Videos": "Videos result for golang",
		})))
		slowEnginesSyncer.LetAllRun()
	})

	It("uses replicas to minimize tail latency impact", func() {
		slowEnginesSyncer := NewSyncer()
		aggregator := Aggregator{
			Engines: Engines{
				"Photos": Replicas{
					fakeSyncedEngine("Photos1", slowEnginesSyncer),
					fakeSyncedEngine("Photos2", slowEnginesSyncer),
					fakeEngine("Photos3"),
				},
				"Videos": Replicas{
					fakeEngine("Videos1"),
					fakeSyncedEngine("Videos2", slowEnginesSyncer),
				},
			}}
		results := make(chan Results)
		go func() { results <- aggregator.Search("golang") }()
		Eventually(results).Should(Receive(Equal(Results{
			"Photos": "Photos3 result for golang",
			"Videos": "Videos1 result for golang",
		})))
		slowEnginesSyncer.LetAllRun()
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

func fakeSyncedEngine(name string, syncer *Syncer) Engine {
	syncer.Register()
	return func(query Query) Result {
		syncer.Sync()
		return fakeEngine(name)(query)
	}
}

func fakeSlowEngine(name string, duration time.Duration) Engine {
	return func(query Query) Result {
		time.Sleep(duration)
		return fakeEngine(name)(query)
	}
}
