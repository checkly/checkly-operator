/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
