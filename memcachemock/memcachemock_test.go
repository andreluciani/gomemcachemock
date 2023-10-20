package memcachemock

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
)

func TestNewFromSelector(t *testing.T) {
	var ss memcache.ServerSelector
	mock := NewFromSelector(&ss)
	a := assert.New(t)
	mock.ExpectPing()
	err := mock.Ping()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestUnexpectedClose(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectPing()
	err := mock.Close()
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestUnexpectedDeleteAll(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectClose()
	err := mock.DeleteAll()
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestUnexpectedFlushAll(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectClose()
	err := mock.FlushAll()
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestUnexpectedPing(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectClose()
	err := mock.Ping()
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestUnexpectedPingAfterExpectationsFullfilled(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)

	mock.ExpectClose()
	err := mock.Close()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
	err = mock.Ping()
	a.Error(err)
}

func TestAdd(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectAdd().
		WithItem(item)
	err := mock.Add(item)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestAdd_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	anotherItem := &memcache.Item{
		Key: "another-key",
	}
	mock.ExpectAdd().
		WithItem(item)
	err := mock.Add(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestAppend(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectAppend().
		WithItem(item)
	err := mock.Append(item)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestAppend_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	anotherItem := &memcache.Item{
		Key: "another-key",
	}
	mock.ExpectAppend().
		WithItem(item)
	err := mock.Append(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestClose(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectClose()
	err := mock.Close()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestClose_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectClose().
		WillReturnError(memcache.ErrServerError)
	err := mock.Close()
	a.ErrorIs(err, memcache.ErrServerError)
	a.NoError(mock.ExpectationsWereMet())
}

func TestCompareAndSwap(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectCompareAndSwap().
		WithItem(item)
	err := mock.CompareAndSwap(item)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestCompareAndSwap_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	anotherItem := &memcache.Item{
		Key: "another-key",
	}
	mock.ExpectCompareAndSwap().
		WithItem(item)
	err := mock.CompareAndSwap(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestDecrement(t *testing.T) {
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

func TestDecrement_ErrorKey(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDecrement().
		WithKeyAndDelta("some-key", 10)
	newValue, err := mock.Decrement("another-key", 10)
	a.Error(err)
	a.Equal(uint64(0), newValue)
	a.Error(mock.ExpectationsWereMet())
}

func TestDecrement_ErrorDelta(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDecrement().
		WithKeyAndDelta("some-key", 10)
	newValue, err := mock.Decrement("some-key", 15)
	a.Error(err)
	a.Equal(uint64(0), newValue)
	a.Error(mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDelete().
		WithKey("some-key")
	err := mock.Delete("some-key")
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestDelete_ErrorKey(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDelete().
		WithKey("some-key")
	err := mock.Delete("another-key")
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestDeleteAll(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDeleteAll()
	err := mock.DeleteAll()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestDeleteAll_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectDeleteAll().
		WillReturnError(memcache.ErrServerError)
	err := mock.DeleteAll()
	a.ErrorIs(err, memcache.ErrServerError)
	a.NoError(mock.ExpectationsWereMet())
}

func TestFlushAll(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectFlushAll()
	err := mock.FlushAll()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestFlushAll_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectFlushAll().
		WillReturnError(memcache.ErrServerError)
	err := mock.FlushAll()
	a.ErrorIs(err, memcache.ErrServerError)
	a.NoError(mock.ExpectationsWereMet())
}

func TestGet(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectGet().
		WithKey("some-key").
		WillReturnItem(item)
	result, err := mock.Get("some-key")
	a.NoError(err)
	a.Equal(item, result)
	a.NoError(mock.ExpectationsWereMet())
}

func TestGet_ErrorKey(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectGet().
		WithKey("some-key")
	result, err := mock.Get("another-key")
	a.Error(err)
	a.Nil(result)
	a.Error(mock.ExpectationsWereMet())
}

func TestGetMulti(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	items := map[string]*memcache.Item{
		"some-key": {
			Key: "some-key",
		},
		"another-key": {
			Key: "another-key",
		},
	}

	mock.ExpectGetMulti().
		WithKeys([]string{"some-key", "another-key"}).
		WillReturnItems(items)
	result, err := mock.GetMulti([]string{"some-key", "another-key"})
	a.NoError(err)
	a.Equal(items, result)
	a.NoError(mock.ExpectationsWereMet())
}

func TestGetMulti_ErrorKeys(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectGetMulti().
		WithKeys([]string{"some-key", "another-key"})
	result, err := mock.GetMulti([]string{"some-key", "unknown-key"})
	a.Error(err)
	a.Nil(result)
	a.Error(mock.ExpectationsWereMet())
}

func TestIncrement(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectIncrement().
		WithKeyAndDelta("some-key", 10).
		WillReturnValue(20)
	newValue, err := mock.Increment("some-key", 10)
	a.NoError(err)
	a.Equal(uint64(20), newValue)
	a.NoError(mock.ExpectationsWereMet())
}

func TestIncrement_ErrorKey(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectIncrement().
		WithKeyAndDelta("some-key", 10)
	newValue, err := mock.Increment("another-key", 10)
	a.Error(err)
	a.Equal(uint64(0), newValue)
	a.Error(mock.ExpectationsWereMet())
}

func TestIncrement_ErrorDelta(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectIncrement().
		WithKeyAndDelta("some-key", 10)
	newValue, err := mock.Increment("some-key", 15)
	a.Error(err)
	a.Equal(uint64(0), newValue)
	a.Error(mock.ExpectationsWereMet())
}

func TestPing(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectPing()
	err := mock.Ping()
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestPing_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectPing().
		WillReturnError(memcache.ErrServerError)
	err := mock.Ping()
	a.ErrorIs(err, memcache.ErrServerError)
	a.NoError(mock.ExpectationsWereMet())
}

func TestPrepend(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectPrepend().
		WithItem(item)
	err := mock.Prepend(item)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestPrepend_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	anotherItem := &memcache.Item{
		Key: "another-key",
	}
	mock.ExpectPrepend().
		WithItem(item)
	err := mock.Prepend(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestReplace(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	mock.ExpectReplace().
		WithItem(item)
	err := mock.Replace(item)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestReplace_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	anotherItem := &memcache.Item{
		Key: "another-key",
	}
	mock.ExpectReplace().
		WithItem(item)
	err := mock.Replace(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestSet(t *testing.T) {
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

func TestSet_Error(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	item := &memcache.Item{
		Key: "some-key",
	}
	anotherItem := &memcache.Item{
		Key: "another-key",
	}
	mock.ExpectSet().
		WithItem(item)
	err := mock.Set(anotherItem)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestTouch(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectTouch().
		WithKeyAndSeconds("some-key", 10)
	err := mock.Touch("some-key", 10)
	a.NoError(err)
	a.NoError(mock.ExpectationsWereMet())
}

func TestTouch_ErrorKey(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectTouch().
		WithKeyAndSeconds("some-key", 10)
	err := mock.Touch("another-key", 10)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}

func TestTouch_ErrorSeconds(t *testing.T) {
	mock := New("localhost:11211")
	a := assert.New(t)
	mock.ExpectTouch().
		WithKeyAndSeconds("some-key", 10)
	err := mock.Touch("some-key", 15)
	a.Error(err)
	a.Error(mock.ExpectationsWereMet())
}
