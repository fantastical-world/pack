package pack

import (
	"sort"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

//KVPack is a simple key/value store to pack stuff in.
type KVPack struct {
	sync.Mutex
	dbLocation string
	timeout    time.Duration
}

//New creates a new KVPack for storing things.
func New(dbLocation string) Pack {
	return &KVPack{dbLocation: dbLocation, timeout: time.Second * 10}
}

//Save will pack the provided thing in the location specified using the name as it's key.
func (k *KVPack) Save(location string, thing Packable) error {
	k.Lock()
	defer k.Unlock()
	db, err := bbolt.Open(k.dbLocation, 0600, &bbolt.Options{Timeout: k.timeout})
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(location))
		name, bytes := thing.Pack()
		err = b.Put([]byte(name), bytes)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (k *KVPack) Get(location, name string) ([]byte, error) {
	var thing []byte
	k.Lock()
	defer k.Unlock()
	db, err := bbolt.Open(k.dbLocation, 0600, &bbolt.Options{Timeout: k.timeout})
	if err != nil {
		return nil, err
	}
	defer db.Close()
	err = db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(location))
		bytes := b.Get([]byte(name))
		if bytes == nil {
			return ErrThingDoesNotExist
		}

		thing = make([]byte, len(bytes))
		copy(thing, bytes)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return thing, nil
}

func (k *KVPack) Delete(location, name string) error {
	k.Lock()
	defer k.Unlock()
	db, err := bbolt.Open(k.dbLocation, 0600, &bbolt.Options{Timeout: k.timeout})
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(location))
		bytes := b.Get([]byte(name))
		if bytes == nil {
			return ErrThingDoesNotExist
		}

		_ = b.Delete([]byte(name))

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (k *KVPack) List(location string) ([]string, error) {
	k.Lock()
	defer k.Unlock()
	var things []string
	db, err := bbolt.Open(k.dbLocation, 0600, &bbolt.Options{Timeout: k.timeout})
	if err != nil {
		return nil, err
	}
	defer db.Close()
	db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(location))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			things = append(things, string(k))
		}

		return nil
	})

	sort.SliceStable(things, func(p, q int) bool {
		return things[p] < things[q]
	})

	return things, nil
}
