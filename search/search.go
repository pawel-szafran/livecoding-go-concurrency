package search

type Query string
type Result string

type Engine func(Query) Result
type Type string
