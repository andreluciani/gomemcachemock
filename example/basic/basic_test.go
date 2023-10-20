package basic

import (
	"testing"

	"github.com/andreluciani/gomemcachemock/memcachemock"
	"github.com/bradfitz/gomemcache/memcache"

	"github.com/stretchr/testify/require"
)

func TestGetSet_ShouldReturnItem(t *testing.T) {
	mock := memcachemock.New("10.0.0.1:11211")
	item := &memcache.Item{
		Key:   "foo",
		Value: []byte("my value"),
	}
	mock.ExpectSet().
		WithItem(item)
	mock.ExpectGet().
		WithKey("foo").
		WillReturnItem(item)
	it, err := SetAndGet(mock, item)
	require.NoError(t, err)
	require.Equal(t, item, it)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSet_GetShouldReturnError(t *testing.T) {
	mock := memcachemock.New("10.0.0.1:11211")
	item := &memcache.Item{
		Key:   "foo",
		Value: []byte("my value"),
	}
	mock.ExpectSet().
		WithItem(item)
	mock.ExpectGet().
		WithKey("foo").
		WillReturnError(memcache.ErrServerError)
	it, err := SetAndGet(mock, item)
	require.Error(t, err)
	require.ErrorIs(t, memcache.ErrServerError, err)
	require.Nil(t, it)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSet_SetShouldReturnError(t *testing.T) {
	mock := memcachemock.New("10.0.0.1:11211")
	item := &memcache.Item{
		Key:   "invalid key",
		Value: []byte("my value"),
	}
	mock.ExpectSet().
		WithItem(item).
		WillReturnError(memcache.ErrMalformedKey)
	it, err := SetAndGet(mock, item)
	require.Error(t, err)
	require.ErrorIs(t, memcache.ErrMalformedKey, err)
	require.Nil(t, it)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
