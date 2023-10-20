package memcachemock

import (
	"fmt"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
)

// Call Modifiers
func TestTimes(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectPing().Times(3)
	err := mock.Ping()
	a.NoError(err)
	err = mock.Ping()
	a.NoError(err)
	a.Error(mock.ExpectationsWereMet())
	err = mock.Ping()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestMaybe(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectPing().Maybe()
	mock.ExpectClose()
	err := mock.Close()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestWillReturnError(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectPing().
		WillReturnError(memcache.ErrServerError)
	err := mock.Ping()
	a.ErrorIs(memcache.ErrServerError, err)
	a.NoError(mock.ExpectationsWereMet())
}

// Common methods
func TestKeyMatches(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDelete().
		WithKey("some-key")
	err := mock.Delete("some-key")
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestKeyMatches_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDelete().
		WithKey("some-key")
	err := mock.Delete("some-other-key")
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestKeysMatch(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	expectedItems := map[string]*memcache.Item{
		"some-key": {
			Key: "some-key",
		},
		"some-other-key": {
			Key: "some-other-key",
		},
	}
	mock.ExpectGetMulti().
		WithKeys([]string{"some-key", "some-other-key"}).
		WillReturnItems(expectedItems)
	items, err := mock.GetMulti([]string{"some-key", "some-other-key"})
	a.NoError(err)
	a.Equal(expectedItems, items)
	a.NoError(mock.ExpectationsWereMet())
}

func TestKeysMatch_KeysSliceWithDifferentLengths(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	expectedKeys := []string{"some-key", "some-other-key"}
	keys := []string{"some-key", "some-other-key", "additional-key"}
	mock.ExpectGetMulti().
		WithKeys(expectedKeys)
	_, err := mock.GetMulti(keys)
	a.Error(err)
	a.ErrorContains(err, fmt.Sprintf("expected keys %s, but got keys %s", expectedKeys, keys))
	a.Error(mock.ExpectationsWereMet())
}

func TestKeysMatch_KeysSliceWithDifferentValues(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	expectedKeys := []string{"some-key", "some-other-key"}
	keys := []string{"some-key", "different-key"}
	mock.ExpectGetMulti().
		WithKeys(expectedKeys)
	_, err := mock.GetMulti(keys)
	a.Error(err)
	a.ErrorContains(err, fmt.Sprintf("expected keys %s, but got keys %s", expectedKeys, keys))
	a.Error(mock.ExpectationsWereMet())
}

func TestItemMatches(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectSet().
		WithItem(item)
	err := mock.Set(item)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestItemMatches_Nil(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectSet().
		WithItem(nil)
	err := mock.Set(nil)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestItemMatches_NilError(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectSet().
		WithItem(nil)
	err := mock.Set(item)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestItemMatches_DifferentKey(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	item := &memcache.Item{
		Key: "some-key",
	}
	anotherItem := &memcache.Item{
		Key: "some-other-key",
	}
	mock.ExpectSet().
		WithItem(item)
	err := mock.Set(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestItemMatches_DifferentValue(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	item := &memcache.Item{
		Value: []byte("some value"),
	}
	anotherItem := &memcache.Item{
		Value: []byte("another value"),
	}
	mock.ExpectSet().
		WithItem(item)
	err := mock.Set(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestItemMatches_DifferentFlags(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	item := &memcache.Item{
		Flags: 10,
	}
	anotherItem := &memcache.Item{
		Flags: 15,
	}
	mock.ExpectSet().
		WithItem(item)
	err := mock.Set(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestItemMatches_DifferentExpiration(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	item := &memcache.Item{
		Expiration: 100,
	}
	anotherItem := &memcache.Item{
		Expiration: 150,
	}
	mock.ExpectSet().
		WithItem(item)
	err := mock.Set(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestItemMatches_DifferentCasID(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	item := &memcache.Item{
		CasID: 10,
	}
	anotherItem := &memcache.Item{
		CasID: 15,
	}
	mock.ExpectSet().
		WithItem(item)
	err := mock.Set(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestDeltaMatches(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDecrement().
		WithKeyAndDelta("some-key", 10).
		WillReturnValue(20)
	newValue, err := mock.Decrement("some-key", 10)
	a.NoError(err)
	a.Equal(uint64(20), newValue)
	a.NoError(mock.ExpectationsWereMet())
}

func TestDeltaMatches_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDecrement().
		WithKeyAndDelta("some-key", 10)
	newValue, err := mock.Decrement("some-key", 15)
	a.Error(err)
	a.Equal(uint64(0), newValue)
	a.Error(mock.ExpectationsWereMet())
}

func TestSecondsMatch(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectTouch().
		WithKeyAndSeconds("some-key", 10)
	err := mock.Touch("some-key", 10)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestSecondsMatches_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectTouch().
		WithKeyAndSeconds("some-key", 10)
	err := mock.Touch("some-key", 15)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestExpectationStrings(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectPing().Maybe().
		WillReturnError(memcache.ErrServerError)
	mock.ExpectDeleteAll()
	mock.ExpectFlushAll()
	mock.ExpectGet().WillReturnItem(&memcache.Item{})
	mock.ExpectGetMulti().WillReturnItems(map[string]*memcache.Item{
		"some-key":    {},
		"another-key": {},
	})
	for _, ex := range mock.expectations {
		fmt.Sprintln(ex)
	}
	a.Error(mock.ExpectationsWereMet())
}
