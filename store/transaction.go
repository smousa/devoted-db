package store

// Transaction is the transactional storage
type Transaction struct {
	data    map[string]string
	deletes map[string]struct{}
	counts  map[string]int
	store   Store
}

var _ Store = &Transaction{}

func NewTransaction(store Store) *Transaction {
	return &Transaction{
		data:    make(map[string]string),
		deletes: make(map[string]struct{}),
		counts:  make(map[string]int),
		store:   store,
	}
}

func (tx *Transaction) Set(key, value string) {
	if _, ok := tx.deletes[key]; ok {
		// undelete locally deleted key
		delete(tx.deletes, key)
		tx.increment(value)
	} else if v, ok := tx.data[key]; ok {
		// update locally set key
		if v != value {
			tx.decrement(v)
			tx.increment(value)
		} else {
			return
		}
	} else if v, ok := tx.store.Get(key); ok {
		// update remote set key
		if v != value {
			tx.decrement(v)
			tx.increment(value)
		} else {
			return
		}
	} else {
		tx.increment(value)
	}

	tx.data[key] = value
}

func (tx *Transaction) decrement(value string) {
	count := tx.counts[value] - 1
	if count != 0 {
		tx.counts[value] = count
	} else {
		delete(tx.counts, value)
	}
}

func (tx *Transaction) increment(value string) {
	tx.counts[value]++
}

func (tx *Transaction) Get(key string) (string, bool) {
	if _, ok := tx.deletes[key]; ok {
		return "", false
	} else if v, ok := tx.data[key]; ok {
		return v, true
	} else if v, ok := tx.store.Get(key); ok {
		return v, true
	}
	return "", false
}

func (tx *Transaction) Delete(key string) {
	if _, ok := tx.deletes[key]; ok {
		// don't care, already deleted
		return
	} else if v, ok := tx.data[key]; ok {
		// remove locally set key
		delete(tx.data, key)
		tx.decrement(v)
	}

	if v, ok := tx.store.Get(key); ok {
		// mark remote key as deleted
		tx.deletes[key] = struct{}{}
		tx.decrement(v)
	}
}

func (tx *Transaction) Count(value string) int {
	return tx.counts[value] + tx.store.Count(value)
}

func (tx *Transaction) Begin() Store {
	return NewTransaction(tx)
}

func (tx *Transaction) Rollback() (Store, error) {
	return tx.store, nil
}

func (tx *Transaction) Commit() Store {
	for key := range tx.deletes {
		tx.store.Delete(key)
	}
	for key, value := range tx.data {
		tx.store.Set(key, value)
	}
	return tx.store.Commit()
}
