package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/byuoitav/common/log"
	"go.etcd.io/bbolt"
)

const (
	dbPath        = "/tmp/cache.db"
	dbPerms       = 0600
	personsBucket = "persons"
	loginsBucket  = "logins"
)

// Person is the in-cache representation of a person containing all the information about
// an individual that lab-attendance cares about
type Person struct {
	BYUID  string
	Name   string
	CardID string
	NetID  string
}

// Cache contains all of the information and things that an instantiation of
// cache needs to run
type Cache struct {
	db *bbolt.DB
}

// New returns a new instantiation of the cache
func New() (*Cache, error) {

	db, err := bbolt.Open(dbPath, dbPerms, &bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("Error while trying to create bbolt db: %s", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(personsBucket))
		if err != nil {
			return fmt.Errorf("Error while creating persons bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(loginsBucket))
		if err != nil {
			return fmt.Errorf("Error while creating logins bucket: %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error while initializing bbolt database: %s", err)
	}

	c := &Cache{
		db: db,
	}

	return c, nil
}

// GetPersonByBYUID gets the person identified by the given BYU ID from cache and returns the record
func (c *Cache) GetPersonByBYUID(byuID string) (Person, error) {

	key := fmt.Sprintf("byuid:%s", byuID)
	record, err := c.getKeyFromCache(key, personsBucket)
	if err != nil {
		return Person{}, fmt.Errorf("Error while trying to retrieve person from cache: %s", err)
	}

	if record == nil {
		return Person{}, errors.New("Person does not exist in cache")
	}

	var p Person
	err = json.Unmarshal(record, &p)
	if err != nil {
		return Person{}, fmt.Errorf("Error while trying to unmarshal person from cache: %s", err)
	}

	return p, nil
}

// GetPersonByCardID gets the person identified by the given Card ID from cache and returns the record
func (c *Cache) GetPersonByCardID(cardID string) (Person, error) {
	key := fmt.Sprintf("cardid:%s", cardID)
	record, err := c.getKeyFromCache(key, personsBucket)
	if err != nil {
		return Person{}, fmt.Errorf("Error while trying to retrieve person from cache: %s", err)
	}

	if record == nil {
		return Person{}, errors.New("Person does not exist in cache")
	}

	var p Person
	err = json.Unmarshal(record, &p)
	if err != nil {
		return Person{}, fmt.Errorf("Error while trying to unmarshal person from cache: %s", err)
	}

	return p, nil
}

// getKeyFromCache gets the given key from the given bucket in cache and returns the value
func (c *Cache) getKeyFromCache(key, bucket string) ([]byte, error) {
	var record []byte
	err := c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		record = b.Get([]byte(key))
		return nil
	})
	return record, err
}

// SavePersonToCache pushes the given person into the cache
func (c *Cache) SavePersonToCache(p Person) error {

	bKey := fmt.Sprintf("byuid:%s", p.BYUID)
	cKey := fmt.Sprintf("cardid:%s", p.CardID)
	val, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Error while trying to marshal person: %s", err)
	}

	err = c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(personsBucket))
		b.Put([]byte(bKey), val)
		if p.CardID != "" {
			b.Put([]byte(cKey), val)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error while trying to save to cache: %s", err)
	}

	log.L.Debugf("Saved person record to cache: %v+", p)

	return nil
}
