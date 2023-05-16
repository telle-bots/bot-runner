package logic_2

import "encoding/json"

type ConditionSourceType uint

const (
	_ ConditionSourceType = iota
	ConditionSourceTypeClock
)

type ConditionSource struct {
	Type       ConditionSourceType `json:"type"`
	Conditions []SourceCondition   `json:"conditions"`
}

type SourceConditionType uint

const (
	_ SourceConditionType = iota
	SourceConditionTypeTime
)

type SourceCondition struct {
	Type SourceConditionType `json:"type"`
	Data json.RawMessage     `json:"data"`
}
