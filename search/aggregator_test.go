package search_test

import (
	"fmt"

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

})

func fakeSearch(name string) Search {
	return func(query Query) Result {
		return Result(fmt.Sprint(name, " result for ", query))
	}
}
