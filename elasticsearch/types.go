package elasticsearch

type statistics struct {
	Indices map[string]interface{} `json:"indices"`
}

type Client interface {
	FetchIndices(pattern string) ([]string, error)
	DeleteIndex(key string) error
}
