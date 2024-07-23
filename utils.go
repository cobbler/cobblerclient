package cobblerclient

import (
	"errors"
)

func returnString(res interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	} else {
		return res.(string), err
	}
}

func returnStringSlice(res interface{}, err error) ([]string, error) {
	var result []string

	if err != nil {
		return nil, err
	}

	for _, name := range res.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

func returnIntSlice(res interface{}, err error) ([]int, error) {
	var result []int

	if err != nil {
		return nil, err
	}

	for _, name := range res.([]interface{}) {
		var parsedInt, err = convertToInt(name)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedInt)
	}

	return result, nil
}

func returnBool(res interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	} else {
		return res.(bool), err
	}
}

func convertIntBool(integer int) (bool, error) {
	if integer == 0 {
		return false, nil
	}
	if integer == 1 {
		return true, nil
	}
	return false, errors.New("integer was neither 0 nor 1")
}

func convertToInt(integer interface{}) (int, error) {
	switch integer.(type) {
	case int8:
		return int(integer.(int8)), nil
	case int16:
		return int(integer.(int16)), nil
	case int32:
		return int(integer.(int32)), nil
	case int64:
		return int(integer.(int64)), nil
	default:
		return -1, errors.New("integer could not be converted")
	}
}
