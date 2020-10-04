package client

type ClientInstance struct {
	DataStore map[string][]byte
}

func (this ClientInstance) Store(id, payload []byte) (aesKey []byte, err error) {
	return []byte("crypt"), nil
}
