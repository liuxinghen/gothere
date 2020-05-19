package gothere

import "reflect"

// load data from a slice of target type to a data map slice
func LoadFromSlice(sourceSlice interface{}, itemType reflect.Type, fields ...string) []map[string]interface{} {
	// TODO : implementation
	return nil
}

func UnloadToSlice(sourceData []map[string]interface{}, t reflect.Type) interface{} {
	// TODO : implementation
	return nil
}
