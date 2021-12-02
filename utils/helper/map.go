package helper

import "reflect"

func MapKeysToArrayOfInterface(mapObject interface{}) []interface{} {
	objType := reflect.TypeOf(mapObject)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	if objType.Kind() != reflect.Map {
		return nil
	}

	objVal := reflect.ValueOf(mapObject)
	keys := objVal.MapKeys()
	keySlice := make([]interface{}, len(keys))
	i := 0
	for _, key := range keys {
		keySlice[i] = key.Interface()
		i++
	}

	return keySlice
}
