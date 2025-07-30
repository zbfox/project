package util

import "strings"

type StringArray []string

func (a *StringArray) NormalizeSlashes() {
	for i, v := range *a {
		(*a)[i] = strings.ReplaceAll(v, "\\", "/")
	}
}
