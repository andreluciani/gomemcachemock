/*
The package memcachemock is used to mock a memcache client.
*/
package memcachemock

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

func New(server ...string) *memcachemock {
	mock := &memcachemock{}
	return mock
}

func NewFromSelector(ss *memcache.ServerSelector) *memcachemock {
	mock := &memcachemock{}
	return mock
}

type gomemcacheMockIface interface {
	// ExpectationsWereMet checks whether all queued expectations
	// were met in order (unless MatchExpectationsInOrder set to false).
	// If any of them was not met - an error is returned.
	ExpectationsWereMet() error

	// ExpectAdd expects Add() to be called with memcache.Item.
	// The *ExpectedAdd allows to mock the response.
	ExpectAdd() *ExpectedAdd

	// ExpectAppend expects Append() to be called with memcache.Item.
	// The *ExpectedAppend allows to mock the response.
	ExpectAppend() *ExpectedAppend

	// ExpectClose expects Close() to be called.
	// The *ExpectedClose allows to mock the response.
	ExpectClose() *ExpectedClose

	// ExpectCompareAndSwap expects CompareAndSwap() to be called with memcache.Item.
	// The *ExpectedCompareAndSwap allows to mock the response.
	ExpectCompareAndSwap() *ExpectedCompareAndSwap

	// ExpectDecrement expects Decrement() to be called with a key and a value.
	// The *ExpectedDecrement allows to mock the response.
	ExpectDecrement() *ExpectedDecrement

	// ExpectDelete expects Delete() to be called with a key.
	// The *ExpectedDelete allows to mock the response.
	ExpectDelete() *ExpectedDelete

	// ExpectDeleteAll expects DeleteAll() to be called.
	// The *ExpectedDeleteAll allows to mock the response.
	ExpectDeleteAll() *ExpectedDeleteAll

	// ExpectFlushAll expects FlushAll() to be called.
	// The *ExpectedFlushAll allows to mock the response.
	ExpectFlushAll() *ExpectedFlushAll

	// ExpectGet expects Get() to be called with a key.
	// The *ExpectedGet allows to mock the response.
	ExpectGet() *ExpectedGet

	// ExpectGetMulti expects GetMulti() to be called with a slice of keys.
	// The *ExpectedGetMulti allows to mock the response.
	ExpectGetMulti() *ExpectedGetMulti

	// ExpectIncrement expects Increment() to be called with a key and a value.
	// The *ExpectedIncrement allows to mock the response.
	ExpectIncrement() *ExpectedIncrement

	// ExpectPing expects Ping() to be called.
	// The *ExpectedPing allows to mock the response.
	ExpectPing() *ExpectedPing

	// ExpectPrepend expects Prepend() to be called with memcache.Item.
	// The *ExpectedPrepend allows to mock the response.
	ExpectPrepend() *ExpectedPrepend

	// ExpectReplace expects Replace() to be called with memcache.Item.
	// The *ExpectedReplace allows to mock the response.
	ExpectReplace() *ExpectedReplace

	// ExpectSet expects Set() to be called with memcache.Item.
	// The *ExpectedSet allows to mock the response.
	ExpectSet() *ExpectedSet

	// ExpectTouch expects Touch() to be called with a key and a number of seconds.
	// The *ExpectedTouch allows to mock the response.
	ExpectTouch() *ExpectedTouch
}

type gomemcacheIface interface {
	gomemcacheMockIface
	Add(item *memcache.Item) error
	Append(item *memcache.Item) error
	Close() error
	CompareAndSwap(item *memcache.Item) error
	Decrement(key string, delta uint64) (newValue uint64, err error)
	Delete(key string) error
	DeleteAll() error
	FlushAll() error
	Get(key string) (item *memcache.Item, err error)
	GetMulti(keys []string) (map[string]*memcache.Item, error)
	Increment(key string, delta uint64) (newValue uint64, err error)
	Ping() error
	Prepend(item *memcache.Item) error
	Replace(item *memcache.Item) error
	Set(item *memcache.Item) error
	Touch(key string, seconds int32) (err error)
}

var _ gomemcacheIface

type memcachemock struct {
	expectations []Expectation
}

func (c *memcachemock) ExpectationsWereMet() error {
	for _, e := range c.expectations {
		fulfilled := e.fulfilled() || !e.required()
		if !fulfilled {
			return fmt.Errorf("there is a remaining expectation which was not matched: %s", e)
		}
	}
	return nil
}

// Expectations Definition Methods
func (c *memcachemock) ExpectAdd() *ExpectedAdd {
	e := &ExpectedAdd{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectAppend() *ExpectedAppend {
	e := &ExpectedAppend{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectClose() *ExpectedClose {
	e := &ExpectedClose{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectCompareAndSwap() *ExpectedCompareAndSwap {
	e := &ExpectedCompareAndSwap{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectDecrement() *ExpectedDecrement {
	e := &ExpectedDecrement{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectDelete() *ExpectedDelete {
	e := &ExpectedDelete{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectDeleteAll() *ExpectedDeleteAll {
	e := &ExpectedDeleteAll{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectFlushAll() *ExpectedFlushAll {
	e := &ExpectedFlushAll{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectGet() *ExpectedGet {
	e := &ExpectedGet{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectGetMulti() *ExpectedGetMulti {
	e := &ExpectedGetMulti{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectIncrement() *ExpectedIncrement {
	e := &ExpectedIncrement{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectPing() *ExpectedPing {
	e := &ExpectedPing{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectPrepend() *ExpectedPrepend {
	e := &ExpectedPrepend{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectReplace() *ExpectedReplace {
	e := &ExpectedReplace{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectSet() *ExpectedSet {
	e := &ExpectedSet{}
	c.expectations = append(c.expectations, e)
	return e
}

func (c *memcachemock) ExpectTouch() *ExpectedTouch {
	e := &ExpectedTouch{}
	c.expectations = append(c.expectations, e)
	return e
}

// Memcache Methods Mocks
func (c *memcachemock) Add(item *memcache.Item) (err error) {
	ex, err := findExpectationFunc[*ExpectedAdd](c, "Add()", func(addExp *ExpectedAdd) error {
		if err := addExp.itemMatches(item); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Append(item *memcache.Item) (err error) {
	ex, err := findExpectationFunc[*ExpectedAppend](c, "Append()", func(appendExp *ExpectedAppend) error {
		if err := appendExp.itemMatches(item); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Close() (err error) {
	ex, err := findExpectation[*ExpectedClose](c, "Close()")
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) CompareAndSwap(item *memcache.Item) (err error) {
	ex, err := findExpectationFunc[*ExpectedCompareAndSwap](c, "CompareAndSwap()", func(compareAndSwapExp *ExpectedCompareAndSwap) error {
		if err := compareAndSwapExp.itemMatches(item); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Decrement(key string, delta uint64) (newValue uint64, err error) {
	ex, err := findExpectationFunc[*ExpectedDecrement](c, "Decrement()", func(decrementExp *ExpectedDecrement) error {
		if err := decrementExp.keyMatches(key); err != nil {
			return err
		}
		if err := decrementExp.deltaMatches(delta); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return ex.value, ex.error()
}

func (c *memcachemock) Delete(key string) (err error) {
	ex, err := findExpectationFunc[*ExpectedDelete](c, "Delete()", func(deleteExp *ExpectedDelete) error {
		if err := deleteExp.keyMatches(key); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) DeleteAll() (err error) {
	ex, err := findExpectation[*ExpectedDeleteAll](c, "DeleteAll()")
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) FlushAll() (err error) {
	ex, err := findExpectation[*ExpectedFlushAll](c, "FlushAll()")
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Get(key string) (item *memcache.Item, err error) {
	ex, err := findExpectationFunc[*ExpectedGet](c, "Get()", func(getExp *ExpectedGet) error {
		if err := getExp.keyMatches(key); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ex.item, ex.error()
}

func (c *memcachemock) GetMulti(keys []string) (items map[string]*memcache.Item, err error) {
	ex, err := findExpectationFunc[*ExpectedGetMulti](c, "GetMulti()", func(getMultiExp *ExpectedGetMulti) error {
		if err := getMultiExp.keysMatch(keys); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ex.items, ex.error()
}

func (c *memcachemock) Increment(key string, delta uint64) (newValue uint64, err error) {
	ex, err := findExpectationFunc[*ExpectedIncrement](c, "Increment()", func(incrementExp *ExpectedIncrement) error {
		if err := incrementExp.keyMatches(key); err != nil {
			return err
		}
		if err := incrementExp.deltaMatches(delta); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return ex.value, ex.error()
}

func (c *memcachemock) Ping() (err error) {
	ex, err := findExpectation[*ExpectedPing](c, "Ping()")
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Prepend(item *memcache.Item) (err error) {
	ex, err := findExpectationFunc[*ExpectedPrepend](c, "Prepend()", func(prependExp *ExpectedPrepend) error {
		if err := prependExp.itemMatches(item); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Replace(item *memcache.Item) (err error) {
	ex, err := findExpectationFunc[*ExpectedReplace](c, "Replace()", func(replaceExp *ExpectedReplace) error {
		if err := replaceExp.itemMatches(item); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Set(item *memcache.Item) (err error) {
	ex, err := findExpectationFunc[*ExpectedSet](c, "Set()", func(setExp *ExpectedSet) error {
		if err := setExp.itemMatches(item); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

func (c *memcachemock) Touch(key string, seconds int32) (err error) {
	ex, err := findExpectationFunc[*ExpectedTouch](c, "Touch()", func(touchExp *ExpectedTouch) error {
		if err := touchExp.keyMatches(key); err != nil {
			return err
		}
		if err := touchExp.secondsMatch(seconds); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ex.error()
}

type ExpectationType[t any] interface {
	*t
	Expectation
}

func findExpectationFunc[ET ExpectationType[t], t any](c *memcachemock, method string, cmp func(ET) error) (ET, error) {
	var expected ET
	var fulfilled int
	var ok bool
	var err error
	for _, next := range c.expectations {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(ET); ok {
			err = cmp(expected)
			if err == nil {
				break
			}
		}
		if (!ok || err != nil) && !next.required() {
			next.Unlock()
			continue
		}
		next.Unlock()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("call to method %s, was not expected, next expectation is: %s", method, next)
	}

	if expected == nil {
		msg := fmt.Sprintf("call to method %s was not expected", method)
		if fulfilled == len(c.expectations) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg)
	}
	defer expected.Unlock()

	expected.fulfill()
	return expected, nil
}

func findExpectation[ET ExpectationType[t], t any](c *memcachemock, method string) (ET, error) {
	return findExpectationFunc[ET, t](c, method, func(_ ET) error { return nil })
}
