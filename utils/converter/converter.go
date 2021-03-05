package converter

import (
	"errors"
	"reflect"
	"strconv"
)

// ToInt function
func ToInt(v interface{}) (int, error) {
	switch v.(type) {
	case int:
		return v.(int), nil
	case int32:
		return int(v.(int32)), nil
	case uint32:
		return int(v.(uint32)), nil
	case float32:
		return int(v.(float32)), nil
	case float64:
		return int(v.(float64)), nil
	case string:
		t, err := strconv.Atoi(v.(string))
		if err != nil {
			return 0.0, errors.New("MongoDB.ToInt: value of string " + v.(string) + " cannot be converted to int")
		}
		return t, nil
	default:
		return 0.0, errors.New("MongoDB.ToInt:" + reflect.ValueOf(v).Kind().String() + " cannot be converted to int")
	}
}

// ToString function
func ToString(v interface{}) (string, error) {
	switch v.(type) {
	case int:
		return strconv.Itoa(v.(int)), nil
	case int32:
		return strconv.Itoa(int(v.(int32))), nil
	case uint32:
		return strconv.Itoa(int(v.(uint32))), nil
	case float32:
		return strconv.Itoa(int(v.(float32))), nil
	case float64:
		return strconv.Itoa(int(v.(float64))), nil
	case string:
		return v.(string), nil
	default:
		return "", errors.New("MongoDB.ToInt:" + reflect.ValueOf(v).Kind().String() + " cannot be converted to int")
	}
}

// ToFloat64 function
func ToFloat64(v interface{}) (float64, error) {
	switch v.(type) {
	case int:
		return float64(v.(int)), nil
	case int32:
		return float64(v.(int32)), nil
	case uint32:
		return float64(v.(uint32)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	case string:
		f, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return 0.0, errors.New("MongoDB: value of string " + v.(string) + " cannot be converted to float64")
		}
		return f, nil
	default:
		return 0.0, errors.New("MongoDB.ToFloat64:" + reflect.ValueOf(v).Kind().String() + " cannot be converted to float64")
	}
}

// ToArrString function
func ToArrString(arr []interface{}) ([]string, error) {
	r := make([]string, len(arr))

	for i, v := range arr {
		switch v.(type) {
		case int:
			r[i] = strconv.Itoa(v.(int))
		case int32:
			r[i] = strconv.Itoa(int(v.(int32)))
		case uint32:
			r[i] = strconv.Itoa(int(v.(uint32)))
		case float32:
			r[i] = strconv.FormatFloat(float64(v.(float32)), 'f', 3, 32)
		case float64:
			r[i] = strconv.FormatFloat(v.(float64), 'f', 3, 32)
		case string:
			r[i] = v.(string)
		default:
			return nil, errors.New("MongoDB.ToArrString:" + reflect.ValueOf(v).Kind().String() + " cannot be converted to string")
		}
	}

	return r, nil
}

// ToArrInt function
func ToArrInt(arr []interface{}) ([]int, error) {
	r := make([]int, len(arr))

	for i, v := range arr {
		switch v.(type) {
		case int:
			r[i] = v.(int)
		case int32:
			r[i] = int(v.(int32))
		case uint32:
			r[i] = int(v.(uint32))
		case float32:
			r[i] = int(v.(float32))
		case float64:
			r[i] = int(v.(float64))
		case string:
			f, err := strconv.Atoi(v.(string))
			if err != nil {
				return nil, errors.New("MongoDB: value of string " + v.(string) + " cannot be converted to int")
			}
			r[i] = f
		default:
			return nil, errors.New("MongoDB.ToArrInt:" + reflect.ValueOf(v).Kind().String() + " cannot be converted to int")
		}
	}

	return r, nil
}
