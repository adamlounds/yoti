package store

type StoreFac struct {
	DocumentStore DocumentStore
}

func NewStoreFactory() (*StoreFac, error) {
	documentStore, err := NewDocumentS3Store()
	if err != nil {
		return nil, err
	}
	return &StoreFac{documentStore}, nil
}
