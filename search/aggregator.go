package search

type Searches map[SearchType]Search
type Results map[SearchType]Result

type Aggregator struct {
	Searches Searches
}

func (a *Aggregator) Search(query Query) Results {
	results := Results{}
	for searchType, search := range a.Searches {
		results[searchType] = search(query)
	}
	return results
}
