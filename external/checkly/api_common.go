package external

import "fmt"

func checkValueString(x string, y string) (value string) {
	if x == "" {
		value = y
	} else {
		value = x
	}
	return
}

func checkValueInt(x int, y int) (value int) {
	if x == 0 {
		value = y
	} else {
		value = x
	}
	return
}

func checkValueArray(x []string, y []string) (value []string) {
	if len(x) == 0 {
		value = y
	} else {
		value = x
	}
	return
}

func getTags(labels map[string]string) (tags []string) {

	for k, v := range labels {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}

	return
}
