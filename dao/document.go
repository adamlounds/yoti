package dao

type DocumentDao struct {
	store DocumentStore
}
type DocumentStore interface {
	Retrieve(id []byte) (payload []byte, err error)
	Store(id, payload []byte) (err error)
}

func (d *DocumentDao) Store(id, payload []byte) (err error) {
	err = d.store.Store(id, payload)
	return err
}

func (d *DocumentDao) Retrieve(id []byte) (payload []byte, err error) {
	payload, err = d.store.Retrieve(id)
	return payload, err
}
