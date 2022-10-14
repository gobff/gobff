package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`^(\w+)(\[[\w*=^<>" ]+\])?(\{.+\})?$`)

type Field struct {
	Key         string
	Children    []Field
	ArrayFilter string
}

func GetFieldsFromPattern(pattern string) ([]Field, error) {
	return getAllFieldsInPattern(trim(pattern, "{", "}"))
}

func getAllFieldsInPattern(pattern string) ([]Field, error) {
	patterns := split(pattern)
	var all []Field
	for _, pattern := range patterns {
		exp, err := getFieldFromPattern(pattern)
		if err != nil {
			return nil, err
		}

		all = append(all, exp)
	}
	return all, nil
}

func getFieldFromPattern(pattern string) (Field, error) {
	var matches = regex.FindStringSubmatch(pattern)
	if len(matches) != 4 {
		return Field{}, fmt.Errorf("invalid expression: %s", pattern)
	}

	var (
		key         = matches[1]
		filter      = matches[2]
		children    = matches[3]
		childrenExp []Field
	)
	if filter != "" {
		filter = trim(filter, "[", "]")
	}
	if children != "" {
		children = trim(children, "{", "}")

		var err error
		childrenExp, err = getAllFieldsInPattern(children)
		if err != nil {
			return Field{}, err
		}
	}
	return Field{
		Key:         key,
		ArrayFilter: filter,
		Children:    childrenExp,
	}, nil
}

func split(str string) []string {
	var (
		char  rune
		slice []string

		ignore, index, lastIndex int
	)
	for index, char = range str {
		switch char {
		case '{':
			ignore++
		case '}':
			ignore--
		case ',':
			if ignore > 0 {
				continue
			}
			slice = append(slice, str[lastIndex:index])
			lastIndex = index + 1
		}
	}
	slice = append(slice, str[lastIndex:])
	return slice
}

func trim(str, prefix, suffix string) string {
	str = strings.TrimPrefix(str, prefix)
	str = strings.TrimSuffix(str, suffix)
	return str
}
