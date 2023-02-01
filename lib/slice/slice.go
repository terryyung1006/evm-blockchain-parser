package slice

import "strings"

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func GetSliceOfKeys[T comparable, V interface{}](targetMap map[T]V) []T {

	keys := make([]T, 0, len(targetMap))
	for k := range targetMap {
		keys = append(keys, k)
	}
	return keys
}

func MinInt(array []int) int {
	var min int = array[0]
	for _, value := range array {
		if min > value {
			min = value
		}
	}
	return min
}

func StringSliceEqualFoldContainCheck(strList []string, str string) bool {
	for _, val := range strList {
		if strings.EqualFold(val, str) {
			return true
		}
	}
	return false
}
