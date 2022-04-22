package main

import (
	"github.com/dgraph-io/badger"
)

// https://colobu.com/2017/10/11/badger-a-performant-k-v-store/
// https://juejin.cn/post/6844903814571491335
type KvDbOp struct {
	DbConn *badger.DB
}

func NewKvDbOp() *KvDbOp {
	r := KvDbOp{}
	r.Init("dnsDbCache")
	return &r
}

func (r *KvDbOp) Init(szDb string) error {
	opts := badger.DefaultOptions(szDb)
	db, err := badger.Open(opts)
	if nil != err {
		return err
	}
	r.DbConn = db
	return nil
}

func (r *KvDbOp) Delete(key string) error {
	err := r.DbConn.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	return err
}

func (r *KvDbOp) Close() {
	r.DbConn.Close()
}

func (r *KvDbOp) Get(key string) (szRst []byte, err error) {
	err = r.DbConn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		// val, err := item.Value()
		err = item.Value(func(val []byte) error {
			szRst = val
			return nil
		})
		return err
	})
	return szRst, err
}

func (r *KvDbOp) Put(key string, data []byte) {
	err := r.DbConn.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), data)
		if err == badger.ErrTxnTooBig {
			_ = txn.Commit()
		}
		return err
	})
	if err != nil {

	}
}
