package store

// Database is the root storage.
type Database struct {
	data map[string]string
	counts map[string]int
}

var _ Store = &Database{}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]string),
		counts: make(map[string]int),
	}
}

func (db *Database) Set(key, value string) {

	v, ok := db.data[key]
	db.data[key] = value

	// update the counts
	if ok {
		if v != value {
			db.decrement(v)
			db.increment(value)
		}
	} else {
		db.increment(value)
	}
}

func (db *Database) decrement(value string) {
	count := db.counts[value] - 1
	if count > 0 {
		db.counts[value] = count
	} else {
		delete(db.counts, value)
	}
}

func (db *Database) increment(value string) {
	db.counts[value]++
}

func (db *Database) Get(key string) (string, bool) {
	v, ok := db.data[key]
	if ok {
		return v, true
	}
	return "", false
}

func (db *Database) Delete(key string) {
	v, ok := db.data[key]
	if ok {
		delete(db.data, key)
		db.decrement(v)
	}
}

func (db *Database) Count(value string) int {
	return db.counts[value]
}

func (db *Database) Begin() Store {
	return NewTransaction(db)
}

func (db *Database) Rollback() (Store, error) {
	return nil, ErrTxNotFound
}

func (db *Database) Commit() Store {
	return db
}
