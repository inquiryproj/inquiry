// Package replacer provides a Replacer interface that replaces a string with a new string.
package replacer

import (
	"fmt"
	"strings"
	"time"
)

// Replacer is an interface that replaces a string with a new string.
type Replacer interface {
	Replace(string) string
}

// NewFuncReplacer returns a Replacer that replaces the following functions:
// - unix() -> unix timestamp
// - unixNano() -> unix timestamp in nanoseconds
// every function definition is replaced with a dynamic value at run time.
func NewFuncReplacer() Replacer {
	return &funcReplacer{
		replacementMap: map[string]func() string{
			"unix()": func() string {
				return fmt.Sprintf("%d", time.Now().Unix())
			},
			"unixNano()": func() string {
				return fmt.Sprintf("%d", time.Now().UnixNano())
			},
		},
	}
}

type funcReplacer struct {
	replacementMap map[string]func() string
}

func (f *funcReplacer) Replace(s string) string {
	for k, v := range f.replacementMap {
		s = strings.ReplaceAll(s, formatKey(k), v())
	}
	return s
}

// NewMapReplacer returns a Replacer that replaces the following variables.
// - ${key} -> value
// every variable definition is replaced with a static value at run time.
func NewMapReplacer(replacementMap map[string]string) Replacer {
	return &variableReplacer{
		replacementMap: replacementMap,
	}
}

type variableReplacer struct {
	replacementMap map[string]string
}

func (v *variableReplacer) Replace(s string) string {
	for k, v := range v.replacementMap {
		s = strings.ReplaceAll(s, formatKey(k), v)
	}
	return s
}

func formatKey(key string) string {
	return fmt.Sprintf("${%s}", key)
}
