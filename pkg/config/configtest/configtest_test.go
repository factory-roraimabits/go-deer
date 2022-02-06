package configtest

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/factory-roraimabits/go-deer/pkg/config/utils"
	"github.com/stretchr/testify/require"
)

const _invalidKey = "invalid"

const (
	_stringKey      = "string"
	_boolKey        = "bool"
	_intKey         = "int"
	_floatKey       = "float"
	_durationKey    = "duration"
	_stringSliceKey = "stringSliceKey"
	_intSliceKey    = "intSliceKey"
	_floatSliceKey  = "floatSliceKey"
	_jsonKey        = "jsonKey"
)

const (
	_stringValue   = "value"
	_boolValue     = true
	_intValue      = 10
	_floatValue    = 12.13
	_durationValue = 1000
)

const (
	_defBool   = false
	_defString = "_defString"
	_defInt    = 18
	_defFloat  = 111.12
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

var (
	_stringSliceValue = []string{"value1", "value2", "value3"}
	_intSliceValue    = []int{1, 2, 3}
	_floatSliceValue  = []float64{1.0, 2.0, 3.0}
	_jsonValue        = Person{"John", "Doe", 35}
	_encodedJsonValue = "{\"FirstName\": \"John\", \"LastName\": \"Doe\", \"Age\": 35}"
)

var (
	_defStringSlice = []string{"default1", "default2", "default3"}
	_defIntSlice    = []int{10, 20, 30}
	_defFloatSlice  = []float64{10.0, 20.0, 30.0}
)

func TestLoadProperties(t *testing.T) {
	properties := newPropsMap()
	manager := Load(properties)
	assertDefaultValues(t, manager)
	assertExistingValues(t, manager)
}

func newPropsMap() map[string]string {
	return map[string]string{
		_stringKey:      _stringValue,
		_boolKey:        strconv.FormatBool(_boolValue),
		_intKey:         strconv.Itoa(_intValue),
		_floatKey:       fmt.Sprint(_floatValue),
		_durationKey:    strconv.Itoa(_durationValue),
		_stringSliceKey: utils.ConvertStringSliceToString(_stringSliceValue),
		_intSliceKey:    utils.ConvertIntSliceToString(_intSliceValue),
		_floatSliceKey:  utils.ConvertFloatSliceToString(_floatSliceValue),
		_jsonKey:        _encodedJsonValue,
	}
}

func assertDefaultValues(t *testing.T, c *Config) {
	require.Equal(t, _defString, c.GetString(_invalidKey, _defString))
	require.Equal(t, false, c.GetBool(_invalidKey, _defBool))
	require.Equal(t, _defInt, c.GetInt(_invalidKey, _defInt))
	require.Equal(t, uint(_defInt), c.GetUint(_invalidKey, _defInt))
	require.Equal(t, _defFloat, c.GetFloat64(_invalidKey, _defFloat))
	require.Equal(t, time.Duration(_defInt), c.GetDuration(_invalidKey, time.Duration(_defInt)))
	require.Equal(t, _defStringSlice, c.GetStringSlice(_invalidKey, _defStringSlice))
	require.Equal(t, _defIntSlice, c.GetIntSlice(_invalidKey, _defIntSlice))
	require.Equal(t, _defFloatSlice, c.GetFloatSlice(_invalidKey, _defFloatSlice))
}

func assertExistingValues(t *testing.T, c *Config) {
	require.Equal(t, _stringValue, c.GetString(_stringKey, ""))
	require.Equal(t, _boolValue, c.GetBool(_boolKey, false))
	require.Equal(t, _intValue, c.GetInt(_intKey, 0))
	require.Equal(t, uint(_intValue), c.GetUint(_intKey, 0))
	require.Equal(t, _floatValue, c.GetFloat64(_floatKey, 0))
	require.Equal(t, time.Duration(_durationValue), c.GetDuration(_durationKey, time.Duration(0)))
	require.Equal(t, _stringSliceValue, c.GetStringSlice(_stringSliceKey, _defStringSlice))
	require.Equal(t, _intSliceValue, c.GetIntSlice(_intSliceKey, _defIntSlice))
	require.Equal(t, _floatSliceValue, c.GetFloatSlice(_floatSliceKey, _defFloatSlice))

	var person Person
	err := c.GetJSONPropertyAndUnmarshal(_jsonKey, &person)
	require.NoError(t, err)
	require.Equal(t, _jsonValue, person)

	expectedProperties := map[string]string{
		_stringKey:      _stringValue,
		_boolKey:        strconv.FormatBool(_boolValue),
		_intKey:         strconv.Itoa(_intValue),
		_floatKey:       fmt.Sprint(_floatValue),
		_durationKey:    strconv.Itoa(_durationValue),
		_stringSliceKey: utils.ConvertStringSliceToString(_stringSliceValue),
		_intSliceKey:    utils.ConvertIntSliceToString(_intSliceValue),
		_floatSliceKey:  utils.ConvertFloatSliceToString(_floatSliceValue),
		_jsonKey:        _encodedJsonValue,
	}
	require.Equal(t, expectedProperties, c.GetAll())
}
