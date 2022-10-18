package main

import "encoding/json"

func Parse[T any](source string) (T, error) {
	var target T
	err := json.Unmarshal([]byte(source), &target)
	return target, err
}

// func Parse[T any](source interface{}) (T, error) {
// 	var target T
// 	err := InterToStruct(source, &target)
// 	return target, err
// }

func InterToStruct(source interface{}, Target interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, Target)
}
