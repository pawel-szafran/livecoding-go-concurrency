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
			"Photos": func(query Query) Result { return Result(fmt.Sprint("Photos result for ", query)) },
			"Videos": func(query Query) Result { return Result(fmt.Sprint("Videos result for ", query)) },
		}}
		results := aggregator.Search("golang")
		Expect(results).To(Equal(Results{
			"Photos": "Photos result for golang",
			"Videos": "Videos result for golang",
		}))
	})

})
