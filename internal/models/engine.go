package models

type TermFreq map[string]uint

type DocInfo struct {
	Terms      TermFreq
	TotalTerms uint
}

type DocIndex map[string]DocInfo

type SearchQueryResult struct {
	path string
	rank float32
}

func NewSearchQueryResult(path string, rank float32) SearchQueryResult {
	return SearchQueryResult{
		path: path,
		rank: rank,
	}
}

func (s *SearchQueryResult) Path() string {
	return s.path
}

func (s *SearchQueryResult) Rank() float32 {
	return s.rank
}
