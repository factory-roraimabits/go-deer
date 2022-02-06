package configtest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/factory-roraimabits/go-deer/pkg/config/utils"
	"github.com/magiconair/properties"
)

type Config struct {
	prop *properties.Properties
}

// Load load the configurations.
func Load(m map[string]string) *Config {
	return &Config{
		prop: properties.LoadMap(m),
	}
}

// GetBool retrieve the property as bool value
func (p *Config) GetBool(key string, value bool) bool {
	return p.prop.GetBool(key, value)
}

// GetString retrieve the property as string value
func (p *Config) GetString(key string, value string) string {
	return p.prop.GetString(key, value)
}

// GetInt retrieve the property as int value
func (p *Config) GetInt(key string, value int) int {
	return p.prop.GetInt(key, value)
}

// GetFloat64 retrieve the property as float value
func (p *Config) GetFloat64(key string, value float64) float64 {
	return p.prop.GetFloat64(key, value)
}

// GetUint retrieve the property as uint value
func (p *Config) GetUint(key string, value uint) uint {
	return p.prop.GetUint(key, value)
}

// GetDuration retrieve the property as duration value
func (p *Config) GetDuration(key string, value time.Duration) time.Duration {
	return p.prop.GetDuration(key, value)
}

// GetAll retrieve all properties
func (p *Config) GetAll() map[string]string {
	return p.prop.Map()
}

// GetStringSlice retrieve the property as string list values
func (p *Config) GetStringSlice(key string, defaultValues []string) []string {
	if v, exist := p.getList(key); exist {
		return v
	}

	return defaultValues
}

// GetIntSlice retrieve the property as int list values
func (p *Config) GetIntSlice(key string, defaultValues []int) []int {
	if v, exist := p.getList(key); exist {
		if result, err := utils.ConvertStringArrayToIntArray(v); err == nil {
			return result
		}
	}

	return defaultValues
}

// GetFloatSlice retrieve the property as float list values
func (p *Config) GetFloatSlice(key string, defaultValues []float64) []float64 {
	if v, exist := p.getList(key); exist {
		if result, err := utils.ConvertStringArrayToFloatArray(v); err == nil {
			return result
		}
	}

	return defaultValues
}

// GetParsedDuration retrieve the property as duration parsed with time.ParseDuration()
func (p *Config) GetParsedDuration(key string, value time.Duration) time.Duration {
	return p.prop.GetParsedDuration(key, value)
}

func (p *Config) getList(key string) (values []string, exist bool) {
	in, exist := p.prop.Get(key)
	v, err := utils.ConvertStringToList(in)

	if err != nil {
		return nil, false
	}

	return v, exist
}

// GetJSONPropertyAndUnmarshal Retrieve json property and unmarshal
func (p *Config) GetJSONPropertyAndUnmarshal(key string, structType interface{}) error {
	in, exist := p.prop.Get(key)

	if !exist {
		return fmt.Errorf("key %s nonexistent ", key)
	}

	return json.Unmarshal([]byte(in), structType)
}
