package elasticsearch

type statistics struct {
	Indices map[string]interface{}
}

type Client interface {
	FetchIndices() ([]string, error)
	DeleteIndex(key string) error
}
