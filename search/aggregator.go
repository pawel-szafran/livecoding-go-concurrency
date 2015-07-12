package search

type Searches map[SearchType]Search
type Results map[SearchType]Result

type Aggregator struct {
	Searches Searches
}

func (a *Aggregator) Search(query Query) Results {
	resultChan := make(chan typedResult)
	a.searchAll(query, resultChan)
	return a.collectResults(resultChan)
}

type typedResult struct {
	searchType SearchType
	result     Result
}

func (a *Aggregator) searchAll(query Query, resultChan chan typedResult) {
	for searchType, search := range a.Searches {
		searchType, search := searchType, search
		go func() {
			resultChan <- typedResult{searchType, search(query)}
		}()
	}
}

func (a *Aggregator) collectResults(resultChan chan typedResult) Results {
	results := Results{}
	for i := 0; i < len(a.Searches); i++ {
		result := <-resultChan
		results[result.searchType] = result.result
	}
	return results
}
