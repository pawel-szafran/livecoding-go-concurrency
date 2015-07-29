package search

import "time"

type Engines map[Type]Replicas
type Replicas []Engine
type Results map[Type]Result

type Aggregator struct {
	Engines Engines
	Timeout time.Duration
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
	Type   Type
	Result Result
}

func (a *Aggregator) newSearchRequest(query Query) *searchRequest {
	return &searchRequest{
		query:       query,
		resultChan:  make(chan typedResult, len(a.Engines)),
		timeoutChan: newTimeoutChan(a.Timeout),
	}
}

func newTimeoutChan(timeout time.Duration) <-chan time.Time {
	if timeout <= 0 {
		return nil
	}
	return time.After(timeout)
}

func (a *Aggregator) searchAll(req *searchRequest) {
	for t, replicas := range a.Engines {
		t, replicas := t, replicas
		go func() {
			req.resultChan <- typedResult{t, replicas.firstResult(req.query)}
		}()
	}
}

func (replicas Replicas) firstResult(query Query) Result {
	resultChan := make(chan Result, len(replicas))
	for _, engine := range replicas {
		engine := engine
		go func() { resultChan <- engine(query) }()
	}
	return <-resultChan
}

func (a *Aggregator) collectResults(req *searchRequest) (results Results) {
	results = Results{}
	for i := 0; i < len(a.Engines); i++ {
		select {
		case r := <-req.resultChan:
			results[r.Type] = r.Result
		case <-req.timeoutChan:
			return
		}
	}
	return
}
