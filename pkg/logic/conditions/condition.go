package conditions

import "strings"

type StringCondition struct {
	Equal      *string `json:"equal"      name:"Equals"`
	StartsWith *string `json:"startsWith" name:"Starts with"`
}

func (c StringCondition) Match(value string) bool {
	if c.Equal != nil {
		return value == *c.Equal
	}
	if c.StartsWith != nil {
		return strings.HasPrefix(value, *c.StartsWith)
	}
	return false
}
