package util

import (
	"APIStatsInterceptor/types"
	"errors"
	"strconv"
	"strings"
)

// GetValNestedMap => gets the right value
func GetValNestedMap(m map[string]interface{}, keys []string) interface{} {
	var ok bool

	for k, v := range keys {
		if k+1 < len(keys) {
			m, ok = m[v].(map[string]interface{})

			if !ok {
				return nil
			}
		} else {
			return m[v]
		}
	}

	return nil
}

// ParseSets => parses the sets
func ParseSets(str string) (resp []types.Set, err error) {
	allSets := strings.Split(str, "|")
	var keyValue []string
	var path []string

	for _, setStr := range allSets {
		keyValue = strings.Split(setStr, "^")
		if len(keyValue) != 3 {
			err = errors.New("invalid set info, expected 3 values but found only " + strconv.Itoa(len(keyValue)) + " " + setStr)
			return
		}

		path = strings.Split(keyValue[1], ",")

		switch keyValue[2] {
		case "COMMA":
			resp = append(resp, types.Set{Name: keyValue[0], Path: path, Option: types.Comma})
		case "PERCENT":
			resp = append(resp, types.Set{Name: keyValue[0], Path: path, Option: types.Percent})
		case "DATA":
			resp = append(resp, types.Set{Name: keyValue[0], Path: path, Option: types.Data})
		default:
			if strings.HasPrefix(keyValue[2], "P") {
				resp = append(resp, types.Set{Name: keyValue[0], Path: path, Option: types.Prefix, OptionalVal: strings.SplitN(keyValue[2], "P", 2)[1]})
			} else if strings.HasPrefix(keyValue[2], "S") {
				resp = append(resp, types.Set{Name: keyValue[0], Path: path, Option: types.Suffix, OptionalVal: strings.SplitN(keyValue[2], "S", 2)[1]})
			} else {
				resp = append(resp, types.Set{Name: keyValue[0], Path: path, Option: types.None})
			}
		}
	}

	return
}

// ParseHeaders => parses the headers
func ParseHeaders(str string) (headers map[string]string, err error) {
	headers = make(map[string]string)
	sets := strings.Split(str, ",")

	for _, v := range sets {
		keyValues := strings.Split(v, ":")

		if len(keyValues) != 2 {
			err = errors.New("invalid set info")
			return
		}

		headers[keyValues[0]] = keyValues[1]
	}

	return
}
