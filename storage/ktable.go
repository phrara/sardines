package storage

import (
	"bytes"
	"sardines/config"
	"fmt"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type KeyTable struct {
	db *leveldb.DB
	mu sync.Mutex
}

func InitKeyTab() (*KeyTable, error) {

	db, err := leveldb.OpenFile(config.Ktab, nil)
	if err != nil {
		return nil, err
	}
	
	return &KeyTable{
		db: db,
		mu: sync.Mutex{},
	}, nil
}

func (k *KeyTable) Put(key string, cids []string) bool {
	val := k.toBytes(cids)
	err := k.db.Put([]byte(key), val, nil)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (k *KeyTable) Get(key string) []string {
	value, err := k.db.Get([]byte(key), nil)
	if err != nil {
		return nil
	} else {
		return k.split(value)
	}
}

func (k *KeyTable) Append(key string, cids []string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	s := k.Get(key)
	s = append(s, cids...)
	return k.Put(key, s)
}

func (k *KeyTable) GetAll() string {
	k.mu.Lock()
	s, _ := k.db.GetSnapshot()
	defer s.Release()
	res := s.String()
	k.mu.Unlock()
	return res
}

func (k *KeyTable) split(val []byte) []string {
	s := strings.Split(string(val), "|")
	return s
}

func (k *KeyTable) toBytes(cids []string) []byte {
	var builder bytes.Buffer
	for _, v := range cids {
		builder.WriteString(v)
		builder.WriteByte('|')
	}
	return builder.Bytes()[:builder.Len()-1]
}

func (k *KeyTable) Close() {
	err := k.db.Close()
	if err != nil {
		fmt.Println(err)
	}
}




