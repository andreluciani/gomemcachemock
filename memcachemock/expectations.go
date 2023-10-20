package memcachemock

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
)

// an Expectation interface
type Expectation interface {
	error() error
	required() bool
	fulfilled() bool
	fulfill()
	sync.Locker
	fmt.Stringer
}

type CallModifier interface {
	Maybe() CallModifier
	Times(n uint) CallModifier
	WillReturnError(err error)
}

// commonExpectation struct
// satisfies the Expectation interface
type commonExpectation struct {
	triggered    uint  // how many times method was called
	err          error // should method return error
	optional     bool  // can method be skipped
	plannedCalls uint  // how many sequentional calls should be made
	sync.Mutex
}

func (e *commonExpectation) error() error {
	return e.err
}

func (e *commonExpectation) required() bool {
	return !e.optional
}

func (e *commonExpectation) fulfilled() bool {
	maxPlannedCalls := uint(1)
	if e.plannedCalls > maxPlannedCalls {
		maxPlannedCalls = e.plannedCalls
	}
	return e.triggered >= maxPlannedCalls
}

func (e *commonExpectation) fulfill() {
	e.triggered++
}

// Maybe allows the expected method call to be optional.
// Not calling an optional method will not cause an error while asserting expectations
func (e *commonExpectation) Maybe() CallModifier {
	e.optional = true
	return e
}

// Times indicates that the expected method should only fire the indicated number of times.
// Zero value is ignored and means the same as one.
func (e *commonExpectation) Times(n uint) CallModifier {
	e.plannedCalls = n
	return e
}

// WillReturnError allows to set an error for the expected method.
func (e *commonExpectation) WillReturnError(err error) {
	e.err = err
}

// String returns string representation
func (e *commonExpectation) String() string {
	w := new(strings.Builder)
	if e.err != nil {
		fmt.Fprintf(w, "\t- returns error: %v\n", e.err)
	}
	if e.optional {
		fmt.Fprint(w, "\t- execution is optional\n")
	}
	if e.plannedCalls > 0 {
		fmt.Fprintf(w, "\t- execution calls awaited: %d\n", e.plannedCalls)
	}
	return w.String()
}

// keyBasedExpectation is a base class that adds a key matching logic
type keyBasedExpectation struct {
	expectedKey string
}

func (e *keyBasedExpectation) keyMatches(key string) error {
	if key != e.expectedKey {
		return fmt.Errorf("expected key %s, but got key %s", e.expectedKey, key)
	}
	return nil
}

// keysBasedExpectation is a base class that adds keys matching logic
type keysBasedExpectation struct {
	expectedKeys []string
}

func (e *keysBasedExpectation) keysMatch(keys []string) error {
	if len(keys) != len(e.expectedKeys) {
		return fmt.Errorf("expected keys %v, but got keys %v", e.expectedKeys, keys)
	}
	if !reflect.DeepEqual(keys, e.expectedKeys) {
		return fmt.Errorf("expected keys %v, but got keys %v", e.expectedKeys, keys)
	}
	return nil
}

// itemBasedExpectation is a base class that adds an memcache.Item matching logic
type itemBasedExpectation struct {
	expectedItem *memcache.Item
}

func (e *itemBasedExpectation) itemMatches(item *memcache.Item) error {
	if e.expectedItem == nil && item == nil {
		return nil
	}
	if e.expectedItem == nil && item != nil {
		return fmt.Errorf("did not expect item, but got item with key %s", item.Key)
	}
	if item.Key != e.expectedItem.Key {
		return fmt.Errorf("expected item with key %s, but got item with key %s", e.expectedItem.Key, item.Key)
	}
	if string(item.Value) != string(e.expectedItem.Value) {
		return fmt.Errorf("expected item with value %s, but got item with value %s", string(e.expectedItem.Value), string(item.Value))
	}
	if item.Flags != e.expectedItem.Flags {
		return fmt.Errorf("expected item with flags %d, but got item with flags %d", e.expectedItem.Flags, item.Flags)
	}
	if item.Expiration != e.expectedItem.Expiration {
		return fmt.Errorf("expected item with expiration %d, but got item with expiration %d", e.expectedItem.Expiration, item.Expiration)
	}
	if item.CasID != e.expectedItem.CasID {
		return fmt.Errorf("expected item with casID %d, but got item with casID %d", e.expectedItem.CasID, item.CasID)
	}
	return nil
}

// deltaBasedExpectation is a base class that adds a delta matching logic
type deltaBasedExpectation struct {
	expectedDelta uint64
}

func (e *deltaBasedExpectation) deltaMatches(delta uint64) error {
	if delta != e.expectedDelta {
		return fmt.Errorf("expected call with delta %d, but got delta %d", e.expectedDelta, delta)
	}
	return nil
}

// secondsBasedExpectation is a base class that adds a seconds matching logic
type secondsBasedExpectation struct {
	expectedSeconds int32
}

func (e *secondsBasedExpectation) secondsMatch(seconds int32) error {
	if seconds != e.expectedSeconds {
		return fmt.Errorf("expected call with seconds: %d, but got seconds: %d", e.expectedSeconds, seconds)
	}
	return nil
}

// Methods Expectations

// ExpectedAdd is used to manage *memcache.Client.Add expectations
type ExpectedAdd struct {
	commonExpectation
	itemBasedExpectation
}

// WithItem will match given expected memcache.Item to actual memcache.Item used when calling memcache.Client.Add().
// If at least one field does not match, it will return an error.
func (e *ExpectedAdd) WithItem(item *memcache.Item) *ExpectedAdd {
	e.expectedItem = item
	return e
}

// String returns string representation
func (e *ExpectedAdd) String() string {
	msg := "ExpectedAdd => expecting call to Add():\n"
	if e.expectedItem != nil {
		msg += fmt.Sprintf("\t- is with item with key: %s\n", e.expectedItem.Key)
		msg += fmt.Sprintf("\t- and with value: %s\n", string(e.expectedItem.Value))
		msg += fmt.Sprintf("\t- and with flags: %d\n", e.expectedItem.Flags)
		msg += fmt.Sprintf("\t- and expiration date: %d\n", e.expectedItem.Expiration)
		msg += fmt.Sprintf("\t- and with casID: %d\n", e.expectedItem.CasID)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedAppend is used to manage *memcache.Client.Append expectations
type ExpectedAppend struct {
	commonExpectation
	itemBasedExpectation
}

// WithItem will match given expected memcache.Item to actual memcache.Item used when calling memcache.Client.Append().
// If at least one field does not match, it will return an error.
func (e *ExpectedAppend) WithItem(item *memcache.Item) *ExpectedAppend {
	e.expectedItem = item
	return e
}

// String returns string representation
func (e *ExpectedAppend) String() string {
	msg := "ExpectedAppend => expecting call to Append():\n"
	if e.expectedItem != nil {
		msg += fmt.Sprintf("\t- is with item with key: %s\n", e.expectedItem.Key)
		msg += fmt.Sprintf("\t- and with value: %s\n", string(e.expectedItem.Value))
		msg += fmt.Sprintf("\t- and with flags: %d\n", e.expectedItem.Flags)
		msg += fmt.Sprintf("\t- and expiration date: %d\n", e.expectedItem.Expiration)
		msg += fmt.Sprintf("\t- and with casID: %d\n", e.expectedItem.CasID)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedClose is used to manage *memcache.Client.Close expectations
type ExpectedClose struct {
	commonExpectation
}

// String returns string representation
func (e *ExpectedClose) String() string {
	return "ExpectedClose => expecting call to Close()\n" + e.commonExpectation.String()
}

// ExpectedCompareAndSwap is used to manage *memcache.Client.CompareAndSwap expectations
type ExpectedCompareAndSwap struct {
	commonExpectation
	itemBasedExpectation
}

// WithItem will match given expected memcache.Item to actual memcache.Item used when calling memcache.Client.CompareAndSwap().
// If at least one field does not match, it will return an error.
func (e *ExpectedCompareAndSwap) WithItem(item *memcache.Item) *ExpectedCompareAndSwap {
	e.expectedItem = item
	return e
}

// String returns string representation
func (e *ExpectedCompareAndSwap) String() string {
	msg := "ExpectedCompareAndSwap => expecting call to CompareAndSwap():\n"
	if e.expectedItem != nil {
		msg += fmt.Sprintf("\t- is with item with key: %s\n", e.expectedItem.Key)
		msg += fmt.Sprintf("\t- and with value: %s\n", string(e.expectedItem.Value))
		msg += fmt.Sprintf("\t- and with flags: %d\n", e.expectedItem.Flags)
		msg += fmt.Sprintf("\t- and expiration date: %d\n", e.expectedItem.Expiration)
		msg += fmt.Sprintf("\t- and with casID: %d\n", e.expectedItem.CasID)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedDecrement is used to manage *memcache.Client.Decrement expectations
type ExpectedDecrement struct {
	commonExpectation
	keyBasedExpectation
	deltaBasedExpectation
	value uint64
}

// WithKeyAndDelta will match given expected key and delta value to actual key and delta used when calling memcache.Client.Decrement().
// If at least one parameter does not match, it will return an error.
func (e *ExpectedDecrement) WithKeyAndDelta(key string, delta uint64) *ExpectedDecrement {
	e.expectedKey = key
	e.expectedDelta = delta
	return e
}

// WillReturnValue specifies the value that will be returned when calling memcache.Client.Decrement().
func (e *ExpectedDecrement) WillReturnValue(value uint64) *ExpectedDecrement {
	e.value = value
	return e
}

// String returns string representation
func (e *ExpectedDecrement) String() string {
	msg := "ExpectedDecrement => expecting call to Decrement():\n"
	msg += fmt.Sprintf("\t- is with key: %s\n", e.expectedKey)
	msg += fmt.Sprintf("\t- and with delta: %d\n", e.expectedDelta)
	return msg + e.commonExpectation.String()
}

// ExpectedDelete is used to manage *memcache.Client.Delete expectations
type ExpectedDelete struct {
	commonExpectation
	keyBasedExpectation
}

// WithKey will match given expected key to actual key used when calling memcache.Client.Delete().
// if the keys do not match, it will return an error.
func (e *ExpectedDelete) WithKey(key string) *ExpectedDelete {
	e.expectedKey = key
	return e
}

// String returns string representation
func (e *ExpectedDelete) String() string {
	msg := "ExpectedDelete => expecting call to Delete():\n"
	msg += fmt.Sprintf("\t- is with key: %s\n", e.expectedKey)
	return msg + e.commonExpectation.String()
}

// ExpectedDeleteAll is used to manage *memcache.Client.DeleteAll expectations
type ExpectedDeleteAll struct {
	commonExpectation
}

// String returns string representation
func (e *ExpectedDeleteAll) String() string {
	return "ExpectedDeleteAll => expecting call to DeleteAll()\n" + e.commonExpectation.String()
}

// ExpectedFlushAll is used to manage *memcache.Client.FlushAll expectations
type ExpectedFlushAll struct {
	commonExpectation
}

// String returns string representation
func (e *ExpectedFlushAll) String() string {
	return "ExpectedFlushAll => expecting call to FlushAll()\n" + e.commonExpectation.String()
}

// ExpectedGet is used to manage *memcache.Client.Get expectations
type ExpectedGet struct {
	commonExpectation
	keyBasedExpectation
	item *memcache.Item
}

// WithKey will match given expected key to actual key used when calling memcache.Client.Get().
// if the keys do not match, it will return an error.
func (e *ExpectedGet) WithKey(key string) *ExpectedGet {
	e.expectedKey = key
	return e
}

// WillReturnItem specifies the memcache.Item that will be returned when calling memcache.Client.Get().
func (e *ExpectedGet) WillReturnItem(item *memcache.Item) *ExpectedGet {
	e.item = item
	return e
}

// String returns string representation
func (e *ExpectedGet) String() string {
	msg := "ExpectedGet => expecting call to Get():\n"
	msg += fmt.Sprintf("\t- is with key: %s\n", e.expectedKey)
	if e.item != nil {
		msg += fmt.Sprintf("\t- returns item with key: %s\n", e.item.Key)
		msg += fmt.Sprintf("\t- and expiration date: %d\n", e.item.Expiration)
		msg += fmt.Sprintf("\t- and with flags: %d\n", e.item.Flags)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedGetMulti is used to manage *memcache.Client.GetMulti expectations
type ExpectedGetMulti struct {
	commonExpectation
	keysBasedExpectation
	items map[string]*memcache.Item
}

// WithKeys will match given expected keys to actual keys used when calling memcache.Client.GetMulti().
// If at least one of the keys does not match, it will return an error.
func (e *ExpectedGetMulti) WithKeys(keys []string) *ExpectedGetMulti {
	e.expectedKeys = keys
	return e
}

// WillReturnItems specifies the map of memcache.Item that will be returned when calling memcache.Client.GetMulti().
func (e *ExpectedGetMulti) WillReturnItems(items map[string]*memcache.Item) *ExpectedGetMulti {
	e.items = items
	return e
}

// String returns string representation
func (e *ExpectedGetMulti) String() string {
	msg := "ExpectedGetMulti => expecting call to GetMulti():\n"
	msg += fmt.Sprintf("\t- is with keys: %s\n", e.expectedKeys)
	if e.items != nil {
		msg += fmt.Sprintf("\t- returns items: %v\n", e.items)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedIncrement is used to manage *memcache.Client.Increment expectations
type ExpectedIncrement struct {
	commonExpectation
	keyBasedExpectation
	deltaBasedExpectation
	value uint64
}

// WithKeyAndDelta will match given expected key and delta value to actual key and delta used when calling memcache.Client.Increment().
// If at least one parameter does not match, it will return an error.
func (e *ExpectedIncrement) WithKeyAndDelta(key string, delta uint64) *ExpectedIncrement {
	e.expectedKey = key
	e.expectedDelta = delta
	return e
}

// WillReturnValue specifies the value that will be returned when calling memcache.Client.Increment().
func (e *ExpectedIncrement) WillReturnValue(value uint64) *ExpectedIncrement {
	e.value = value
	return e
}

// String returns string representation
func (e *ExpectedIncrement) String() string {
	msg := "ExpectedIncrement => expecting call to Increment():\n"
	msg += fmt.Sprintf("\t- is with key: %s\n", e.expectedKey)
	msg += fmt.Sprintf("\t- and with delta: %d\n", e.expectedDelta)
	return msg + e.commonExpectation.String()
}

// ExpectedPing is used to manage *memcache.Client.Ping expectations
type ExpectedPing struct {
	commonExpectation
}

// String returns string representation
func (e *ExpectedPing) String() string {
	msg := "ExpectedPing => expecting call to Ping()\n"
	return msg + e.commonExpectation.String()
}

// ExpectedPrepend is used to manage *memcache.Client.Prepend expectations
type ExpectedPrepend struct {
	commonExpectation
	itemBasedExpectation
}

// WithItem will match given expected memcache.Item to actual memcache.Item used when calling memcache.Client.Prepend().
// If at least one field does not match, it will return an error.
func (e *ExpectedPrepend) WithItem(item *memcache.Item) *ExpectedPrepend {
	e.expectedItem = item
	return e
}

// String returns string representation
func (e *ExpectedPrepend) String() string {
	msg := "ExpectedPrepend => expecting call to Prepend():\n"
	if e.expectedItem != nil {
		msg += fmt.Sprintf("\t- is with item with key: %s\n", e.expectedItem.Key)
		msg += fmt.Sprintf("\t- and with value: %s\n", string(e.expectedItem.Value))
		msg += fmt.Sprintf("\t- and with flags: %d\n", e.expectedItem.Flags)
		msg += fmt.Sprintf("\t- and expiration date: %d\n", e.expectedItem.Expiration)
		msg += fmt.Sprintf("\t- and with casID: %d\n", e.expectedItem.CasID)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedReplace is used to manage *memcache.Client.Replace expectations
type ExpectedReplace struct {
	commonExpectation
	itemBasedExpectation
}

// WithItem will match given expected memcache.Item to actual memcache.Item used when calling memcache.Client.Replace().
// If at least one field does not match, it will return an error.
func (e *ExpectedReplace) WithItem(item *memcache.Item) *ExpectedReplace {
	e.expectedItem = item
	return e
}

// String returns string representation
func (e *ExpectedReplace) String() string {
	msg := "ExpectedReplace => expecting call to Replace():\n"
	if e.expectedItem != nil {
		msg += fmt.Sprintf("\t- is with item with key: %s\n", e.expectedItem.Key)
		msg += fmt.Sprintf("\t- and with value: %s\n", string(e.expectedItem.Value))
		msg += fmt.Sprintf("\t- and with flags: %d\n", e.expectedItem.Flags)
		msg += fmt.Sprintf("\t- and expiration date: %d\n", e.expectedItem.Expiration)
		msg += fmt.Sprintf("\t- and with casID: %d\n", e.expectedItem.CasID)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedSet is used to manage *memcache.Client.Set expectations
type ExpectedSet struct {
	commonExpectation
	itemBasedExpectation
}

// WithItem will match given expected memcache.Item to actual memcache.Item used when calling memcache.Client.Set().
// If at least one field does not match, it will return an error.
func (e *ExpectedSet) WithItem(item *memcache.Item) *ExpectedSet {
	e.expectedItem = item
	return e
}

// String returns string representation
func (e *ExpectedSet) String() string {
	msg := "ExpectedSet => expecting call to Set():\n"
	if e.expectedItem != nil {
		msg += fmt.Sprintf("\t- is with item with key: %s\n", e.expectedItem.Key)
		msg += fmt.Sprintf("\t- and with value: %s\n", string(e.expectedItem.Value))
		msg += fmt.Sprintf("\t- and with flags: %d\n", e.expectedItem.Flags)
		msg += fmt.Sprintf("\t- and expiration date: %d\n", e.expectedItem.Expiration)
		msg += fmt.Sprintf("\t- and with casID: %d\n", e.expectedItem.CasID)
	}
	return msg + e.commonExpectation.String()
}

// ExpectedTouch is used to manage *memcache.Client.Touch expectations
type ExpectedTouch struct {
	commonExpectation
	keyBasedExpectation
	secondsBasedExpectation
}

// WithKeyAndSeconds will match given expected key and seconds value to actual key and seconds used when calling memcache.Client.Touch().
// If at least one parameter does not match, it will return an error.
func (e *ExpectedTouch) WithKeyAndSeconds(key string, seconds int32) *ExpectedTouch {
	e.expectedKey = key
	e.expectedSeconds = seconds
	return e
}

// String returns string representation
func (e *ExpectedTouch) String() string {
	msg := "ExpectedTouch => expecting call to Touch():\n"
	msg += fmt.Sprintf("\t- is with key: %s\n", e.expectedKey)
	msg += fmt.Sprintf("\t- and with seconds: %d\n", e.expectedSeconds)
	return msg + e.commonExpectation.String()
}
