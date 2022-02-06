package config

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/factory-roraimabits/go-deer/pkg/config/utils"
	"github.com/magiconair/properties"
)

const (
	_defaultConfigPath      = "/configs/latest/application.properties"
	_propertyConfigFileName = "configFileName"
	_checksumEnabled        = "checksumEnabled"
)

// Config provides all configurations loaded from the fury's configuration.
type Config struct {
	prop     *properties.Properties
	filename string
}

// Load loads the configurations.
func Load() (*Config, error) {
	if c := os.Getenv(_propertyConfigFileName); c != "" {
		return load(c)
	}

	return load(_defaultConfigPath)
}

func load(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading configuration: %v", err)
	}

	if err = verify(b, filename); err != nil {
		return nil, fmt.Errorf("verifying configuration: %v", err)
	}

	prop, err := properties.Load(b, properties.UTF8)
	if err != nil {
		return nil, fmt.Errorf("loading configuration: %v", err)
	}

	return &Config{
		prop:     prop,
		filename: filename,
	}, nil
}

func verify(b []byte, filename string) error {
	if c := os.Getenv(_checksumEnabled); c == "false" {
		return nil
	}

	if len(b) == 0 {
		return fmt.Errorf("the file %s is empty", filename)
	}

	filenameMD5 := filename + ".md5"

	expectedMD5, err := ioutil.ReadFile(filenameMD5)
	if err != nil {
		return err
	}

	if len(expectedMD5) == 0 {
		return fmt.Errorf("the file %s is empty", filenameMD5)
	}

	currentMD5, err := md5FromBytes(b)
	if err != nil {
		return err
	}

	if !bytes.Equal(currentMD5, expectedMD5) {
		return fmt.Errorf("different md5 contents")
	}

	return nil
}

func md5FromBytes(b []byte) ([]byte, error) {
	hash := md5.New() //nolint:gosec
	if _, err := hash.Write(b); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := hex.NewEncoder(&buf).Write(hash.Sum(nil)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

// getList retrieve the property as list values
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
