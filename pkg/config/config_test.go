package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Given
	configFileName := "testdata/valid.properties"
	os.Setenv("configFileName", configFileName)
	os.Setenv("checksumEnabled", "false")

	// When
	cfg, err := Load()

	// Then
	require.NoError(t, err)

	require.Equal(t, "value", cfg.GetString("string", ""))
	require.Equal(t, "default-value", cfg.GetString("non-existent-value", "default-value"))
	require.Equal(t, "other-default-value", cfg.GetString("non-existent-value", "other-default-value"))

	require.Equal(t, true, cfg.GetBool("bool", false))
	require.Equal(t, true, cfg.GetBool("non-existent-value", true))
	require.Equal(t, false, cfg.GetBool("non-existent-value", false))

	require.Equal(t, 9, cfg.GetInt("int", 0))
	require.Equal(t, 12, cfg.GetInt("non-existent-value", 12))
	require.Equal(t, 18, cfg.GetInt("non-existent-value", 18))

	require.Equal(t, uint(9), cfg.GetUint("int", 0))
	require.Equal(t, uint(12), cfg.GetUint("non-existent-value", 12))
	require.Equal(t, uint(18), cfg.GetUint("non-existent-value", 18))

	require.Equal(t, 9.12, cfg.GetFloat64("float", 0))
	require.Equal(t, 12.18, cfg.GetFloat64("non-existent-value", 12.18))
	require.Equal(t, 912.18, cfg.GetFloat64("non-existent-value", 912.18))

	require.Equal(t, time.Duration(91218), cfg.GetDuration("duration", time.Duration(0)))
	require.Equal(t, time.Duration(14515), cfg.GetDuration("non-existent-value", time.Duration(14515)))
	require.Equal(t, time.Duration(14318), cfg.GetDuration("non-existent-value", time.Duration(14318)))

	require.Equal(t, time.Duration(91218), cfg.GetParsedDuration("format.duration", time.Duration(0)))
	require.Equal(t, time.Duration(14515), cfg.GetParsedDuration("non-existent-value", time.Duration(14515)))
	require.Equal(t, time.Duration(14318), cfg.GetParsedDuration("non-existent-value", time.Duration(14318)))

	expectedProperties := map[string]string{
		"bool":              "true",
		"duration":          "91218",
		"float":             "9.12",
		"int":               "9",
		"string":            "value",
		"int.list":          "10,15,90",
		"float.list":        "1.0,1.5,9.0",
		"int.invalid.list":  "10,15,a",
		"string.list":       "\"a1,a2\",b1,c1",
		"format.duration":   "91.218µs",
		"json.car.property": `{"id":10,"model":"Gol 1.6","year":2016,"price":"50.00","maker":{"id":13,"name":"Volkswagen","offices":[{"id":89,"name":"The Volkswagen 387","address":{"id":56,"street":"Av. Paulista","number":1007,"neighborhood":"Bela vista"}},{"id":75,"name":"The Volkswagen 1002","address":{"id":56,"street":"Av. Brigadeiro Faria Lima","number":6583,"neighborhood":"Faria Lima"}}]}}`,
	}

	require.Equal(t, expectedProperties, cfg.GetAll())

	expectedStringListValues := []string{"a1,a2", "b1", "c1"}
	require.Equal(t, expectedStringListValues, cfg.GetStringSlice("string.list", []string{}))

	expectedDefaultStringListValues := []string{"a", "b", "c"}
	require.Equal(t, expectedDefaultStringListValues, cfg.GetStringSlice("non-existent-value", expectedDefaultStringListValues))

	expectedIntegerListValues := []int{10, 15, 90}
	require.Equal(t, expectedIntegerListValues, cfg.GetIntSlice("int.list", []int{}))

	expectedFloatListValues := []float64{1.0, 1.5, 9.0}
	require.Equal(t, expectedFloatListValues, cfg.GetFloatSlice("float.list", []float64{}))
	expectedFloatListDfaultValues := []float64{1.0, 1.5, 9.0}
	require.Equal(t, expectedFloatListDfaultValues, cfg.GetFloatSlice("non-existent-value", expectedFloatListDfaultValues))
	expectedDefaultIntegerListValues := []int{231, 582}
	require.Equal(t, expectedDefaultIntegerListValues, cfg.GetIntSlice("int.invalid.list", expectedDefaultIntegerListValues))
	require.Equal(t, expectedDefaultIntegerListValues, cfg.GetIntSlice("non-existent-value", expectedDefaultIntegerListValues))

	var car Car

	noError := cfg.GetJSONPropertyAndUnmarshal("json.car.property", &car)
	require.NoError(t, noError)
	require.Equal(t, "Gol 1.6", car.Model)
	require.Equal(t, "Volkswagen", car.Maker.Name)
	require.Equal(t, "Av. Paulista", car.Maker.Offices[0].Address.Street)
	require.Equal(t, "Av. Brigadeiro Faria Lima", car.Maker.Offices[1].Address.Street)

}

func TestLoad_err(t *testing.T) {
	tt := []struct {
		name          string
		filename      string
		expectedError string
	}{
		{
			name:          "non-existent properties",
			filename:      "testdata/non-existent.properties",
			expectedError: "reading configuration: open testdata/non-existent.properties: no such file or directory",
		},
		{
			name:          "empty properties",
			filename:      "testdata/empty.properties",
			expectedError: "verifying configuration: the file testdata/empty.properties is empty",
		},
		{
			name:          "non-existent md5",
			filename:      "testdata/non-existent-md5.properties",
			expectedError: "verifying configuration: open testdata/non-existent-md5.properties.md5: no such file or directory",
		},
		{
			name:          "empty md5",
			filename:      "testdata/empty-md5.properties",
			expectedError: "verifying configuration: the file testdata/empty-md5.properties.md5 is empty",
		},
		{
			name:          "invalid md5",
			filename:      "testdata/invalid-md5.properties",
			expectedError: "verifying configuration: different md5 contents",
		},
		{
			name:          "default config",
			expectedError: "reading configuration: open /configs/latest/application.properties: no such file or directory",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// When
			_ = os.Setenv("configFileName", tc.filename)
			_ = os.Setenv("checksumEnabled", "true")
			_, err := Load()
			// Then
			require.EqualError(t, err, tc.expectedError)
		})
	}

	aa := []struct {
		name          string
		filename      string
		expectedError string
	}{
		{
			name:          "string",
			filename:      "testdata/valid.properties",
			expectedError: "invalid character 'v' looking for beginning of value",
		},
		{
			name:          "nonexistent.car.property",
			filename:      "testdata/valid.properties",
			expectedError: "key nonexistent.car.property nonexistent ",
		},
	}
	for _, tc := range aa {
		t.Run(tc.name, func(t *testing.T) {
			// When
			_ = os.Setenv("configFileName", tc.filename)
			_ = os.Setenv("checksumEnabled", "false")
			cfg, _ := Load()
			var car Car

			// Then
			err := cfg.GetJSONPropertyAndUnmarshal(tc.name, &car)
			require.EqualError(t, err, tc.expectedError)
		})
	}
}

type Car struct {
	ID    int
	Model string
	Year  int
	Price string
	Maker Maker
}

type Maker struct {
	ID      int
	Name    string
	Offices []Offices
}

type Offices struct {
	ID      int
	Name    string
	Address Address
}

type Address struct {
	ID           int
	Street       string
	Number       int
	Neighborhood string
}
