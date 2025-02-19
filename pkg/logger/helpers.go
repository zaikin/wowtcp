package logger

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	delimiter = "==========================="
)

func Prettify(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")

	return string(s)
}

func PP(i ...interface{}) {
	//nolint:forbidigo
	fmt.Println(delimiter)

	for _, j := range i {
		//nolint:forbidigo
		fmt.Println(Prettify(j))
	}

	//nolint:forbidigo
	fmt.Println(delimiter)
}

// GetPP returns Prettified Print.
func GetPP(i ...interface{}) string {
	r := []string{}
	for _, j := range i {
		r = append(r, Prettify(j))
	}

	return strings.Join(r, "")
}

// GetJSON returns JSON.
func GetJSON(i interface{}) string {
	bytes, _ := json.Marshal(i)

	return string(bytes)
}
