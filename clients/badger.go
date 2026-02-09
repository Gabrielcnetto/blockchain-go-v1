package clients

import (
	"github.com/dgraph-io/badger/v4"
)

func StartBadger() (*badger.DB, error) {
	path := "/home/netto/Desktop/blockchain/my-blockchain/tmp"
	instance_badger := badger.DefaultOptions(path)
	instance_badger.ValueDir = path
	instance_badger.Dir = path
	db, err := badger.Open(instance_badger)
	if err != nil {
		return nil, err
	}
	return db, nil
}
