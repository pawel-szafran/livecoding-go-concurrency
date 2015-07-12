package search

type Query string
type Result string

type Search func(Query) Result
type SearchType string
