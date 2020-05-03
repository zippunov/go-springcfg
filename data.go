package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type springReplacer struct {
	Opener  string
	Closer  string
	Default string
}

// Replace replaces all variables with the defined opening and
// closing strings with and default separator the value of the
// key or when available with the default value.
func (rpl *springReplacer) Replace(str string, d *Data) string {
	var f, s int
	f = strings.Index(str, rpl.Opener) + len(rpl.Opener)
	for f-len(rpl.Opener) > -1 {
		s = f + strings.Index(str[f:], rpl.Closer)
		key := str[f:s]
		var replacement interface{}
		var found, useDefVal bool
		defPosition := strings.Index(key, rpl.Default)
		if defPosition >= 0 {
			useDefVal = true
			replacement, found = d.getEnvOrRaw(key[:defPosition])
		} else {
			replacement, found = d.getEnvOrRaw(key)
		}
		if !found {
			if useDefVal {
				replacement = key[defPosition+1:]
			} else {
				replacement = ""
			}
		}
		str = str[:f-len(rpl.Opener)] + fmt.Sprintf("%v", replacement) + str[s+len(rpl.Closer):]
		f = strings.Index(str, rpl.Opener) + len(rpl.Opener)
	}
	return str
}

var defaultReplacer = springReplacer{
	Opener:  "${",
	Closer:  "}",
	Default: ":",
}

// Data holds configuration data and provides queries with data path
type Data struct {
	Data     map[string]interface{} `json:"data"`
	keyDelim string
	replacer *springReplacer
}

// NewData creates new Data struct instance
func NewData(maps ...map[string]interface{}) Data {
	ms := append(maps, nil)
	copy(ms[1:], ms)
	ms[0] = map[string]interface{}{}
	return Data{
		keyDelim: ".",
		replacer: &defaultReplacer,
		Data:     DeepMergeAll(ms...),
	}
}

// Merge returns new instance of Data with deep merged fields
func (c *Data) Merge(d Data) Data {
	if d.Data == nil {
		return Data{
			keyDelim: ".",
			Data:     c.Data,
			replacer: &defaultReplacer,
		}
	}
	if c.Data == nil {
		return Data{
			keyDelim: ".",
			Data:     d.Data,
			replacer: &defaultReplacer,
		}
	}
	return Data{
		keyDelim: ".",
		Data:     DeepMergeAll(c.Data, d.Data),
		replacer: &defaultReplacer,
	}
}

func (c *Data) getEnvOrRaw(key string) (interface{}, bool) {
	val, ok := os.LookupEnv(key)
	if ok {
		return val, true
	}
	if c.Data == nil {
		return nil, false
	}
	path := buildPath(key, c.keyDelim)
	return searchMap(path, c.Data)
}

// Get returns the value associated with the key as a any type.
func (c *Data) Get(key string) interface{} {
	if c.Data == nil {
		return nil
	}
	path := buildPath(key, c.keyDelim)
	val, _ := searchMap(path, c.Data)
	switch rslt := val.(type) {
	case string:
		processed := c.replacer.Replace(rslt, c)
		return processed
	default:
		return rslt
	}
}

func (c *Data) Has(key string) bool {
	if c.Data == nil {
		return false
	}
	path := buildPath(key, c.keyDelim)
	_, ok := searchMap(path, c.Data)
	return ok
}

// Sub returns data subtree as string map
func (c *Data) Sub(key string) *Data {
	val := c.Get(key)
	switch data := val.(type) {
	case map[string]interface{}:
		return &Data{
			keyDelim: ".",
			Data:     data,
			replacer: &defaultReplacer,
		}
	default:
		return nil
	}
}

// GetString returns the value associated with the key as a string.
func (c *Data) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (c *Data) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (c *Data) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

// GetInt32 returns the value associated with the key as an integer.
func (c *Data) GetInt32(key string) int32 {
	return cast.ToInt32(c.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Data) GetInt64(key string) int64 {
	return cast.ToInt64(c.Get(key))
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Data) GetUint(key string) uint {
	return cast.ToUint(c.Get(key))
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func (c *Data) GetUint32(key string) uint32 {
	return cast.ToUint32(c.Get(key))
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Data) GetUint64(key string) uint64 {
	return cast.ToUint64(c.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Data) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

// GetTime returns the value associated with the key as time.
func (c *Data) GetTime(key string) time.Time {
	return cast.ToTime(c.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (c *Data) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (c *Data) GetIntSlice(key string) []int {
	return cast.ToIntSlice(c.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Data) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Data) GetStringMap(key string) map[string]interface{} {
	return c.Sub(key).Data
}

func searchMap(path []string, source map[string]interface{}) (interface{}, bool) {
	if len(path) == 0 {
		return source, true
	}

	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next, true
		}
		if m, ok := next.(map[string]interface{}); ok {
			return searchMap(path[1:], m)
		}
	}
	return nil, false
}
