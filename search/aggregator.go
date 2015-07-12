package search

import "time"

type Searches map[SearchType]Replicas
type Replicas []Search
type Results map[SearchType]Result

type Aggregator struct {
	Searches Searches
	Timeout  time.Duration
}

func (a *Aggregator) Search(query Query) Results {
	req := a.newSearchRequest(query)
	a.searchAll(req)
	return a.collectResults(req)
}

type searchRequest struct {
	query       Query
	resultChan  chan typedResult
	timeoutChan <-chan time.Time
}

type typedResult struct {
	searchType SearchType
	result     Result
}

func (a *Aggregator) newSearchRequest(query Query) *searchRequest {
	return &searchRequest{
		query:       query,
		resultChan:  make(chan typedResult),
		timeoutChan: time.After(a.Timeout),
	}
}

func (a *Aggregator) searchAll(req *searchRequest) {
	for searchType, replicas := range a.Searches {
		searchType, replicas := searchType, replicas
		go func() {
			req.resultChan <- typedResult{searchType, replicas.firstResult(req.query)}
		}()
	}
}

func (replicas Replicas) firstResult(query Query) Result {
	resultChan := make(chan Result)
	for _, search := range replicas {
		search := search
		go func() {
			resultChan <- search(query)
		}()
	}
	return <-resultChan
}

func (a *Aggregator) collectResults(req *searchRequest) (results Results) {
	results = Results{}
	for i := 0; i < len(a.Searches); i++ {
		select {
		case result := <-req.resultChan:
			results[result.searchType] = result.result
		case <-req.timeoutChan:
			return
		}
	}
	return
}
