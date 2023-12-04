# About

This is a [gomemcache](https://github.com/bradfitz/gomemcache) client mock library.

The implementation is based on the [pgxmock](https://github.com/pashagolub/pgxmock) library, and it has a similar purpose: to simulate a [**gomemcache**](https://github.com/bradfitz/gomemcache) client without an actual memcached server connection.

# Install

```shell
go get github.com/andreluciani/gomemcachemock/memcachemock
```

# Usage

The mock API allows little to no change in the code, see the example below:

## Code that uses gomemcache

```go
package example

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheInterface interface {
	Set(item *memcache.Item) error
	Get(key string) (item *memcache.Item, err error)
}

func SetAndGet(mc MemcacheInterface, item *memcache.Item) (*memcache.Item, error) {
	if err := mc.Set(item); err != nil {
		return nil, err
	}
	return mc.Get(item.Key)
}
```

## Testing the code above

```go
package example

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
```

# Tests

```shell
go test -cover -v ./...
```

# Docs

See https://pkg.go.dev/github.com/andreluciani/gomemcachemock/memcachemock

Or, with [pkgsite](https://github.com/golang/pkgsite) installed, run:

```shell
pkgsite
```

# License

This library is distributed under the [3-Clause BSD License](https://opensource.org/license/bsd-3-clause/)
