package markov

import (
	"sync"
)

var (
	dict dictionary = NewDictionary()
)

type dictionary struct {
	pool map[token]*token // clumsy but wcyd
	mut  sync.Mutex
}

func NewDictionary() dictionary {
	return dictionary{
		pool: make(map[token]*token),
		mut:  sync.Mutex{},
	}
}

func (d *dictionary) Add(value token) {
	d.mut.Lock()
	d.pool[value] = &value
	d.mut.Unlock()
}

func (d *dictionary) Get(key token) (*token, bool) {
	d.mut.Lock()
	ptr, ok := d.pool[key]
	d.mut.Unlock()
	return ptr, ok
}

func (d *dictionary) Intern(key token) *token {
	if ptr, ok := d.Get(key); ok {
		return ptr
	}
	d.Add(key)
	return d.pool[key]
}
