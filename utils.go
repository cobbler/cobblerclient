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
	if err != nil {
		return nil, err
	}

	resConverted, ok := res.([]interface{})

	if !ok {
		return nil, errors.New("result is not a slice")
	}

	return convertToStringSlice(resConverted)
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

func convertToStringSlice(data []interface{}) ([]string, error) {
	result := make([]string, 0)
	for _, name := range data {
		convertedData, ok := name.(string)
		if !ok {
			return nil, errors.New("convertToStringSlice: data is not a string")
		}
		result = append(result, convertedData)
	}
	return result, nil
}

func convertXmlRpcBool(data interface{}) (bool, error) {
	convertedData, ok := data.(int)
	if !ok {
		return false, errors.New("convertXmlRpcBool: data is not a number")
	}
	if convertedData == 1 {
		return true, nil
	} else if convertedData == 0 {
		return false, nil
	} else {
		return false, errors.New("boolean needs to be 0 or 1 according to XML-RPC spec")
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

func convertToFloat(float interface{}) (float64, error) {
	switch v := float.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return -1, errors.New("float could not be converted")
	}
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

func stringInSlice(a string, list []string) bool {
	// https://stackoverflow.com/a/15323988
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
