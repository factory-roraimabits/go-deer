package utils

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

func ConvertStringToList(value string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(value))
	v, err := r.Read()

	if err != nil {
		return nil, err
	}

	return v, nil
}

func ConvertStringArrayToIntArray(values []string) ([]int, error) {
	v := []int{}

	for _, i := range values {
		j, err := strconv.Atoi(i)
		v = append(v, j)

		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func ConvertStringArrayToFloatArray(values []string) ([]float64, error) {
	v := []float64{}

	for _, i := range values {
		j, err := strconv.ParseFloat(i, 64)
		v = append(v, j)

		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func ConvertStringSliceToString(values []string) string {
	return strings.Join(values, ",")
}

func ConvertIntSliceToString(values []int) string {
	if len(values) == 0 {
		return ""
	}

	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = fmt.Sprint(v)
	}

	return strings.Join(strValues, ",")
}

func ConvertFloatSliceToString(values []float64) string {
	if len(values) == 0 {
		return ""
	}

	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = fmt.Sprint(v)
	}

	return strings.Join(strValues, ",")
}
