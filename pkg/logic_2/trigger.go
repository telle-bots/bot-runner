package logic_2

import "encoding/json"

type TriggerType uint

const (
	_ TriggerType = iota
	TriggerTypeBot
)

type Trigger struct {
	Type   TriggerType `json:"type"`
	Events []Event     `json:"events"`
}

type EventType uint

const (
	_ EventType = iota
	EventTypeUpdate
)

type Event struct {
	Type       EventType        `json:"type"`
	Conditions []EventCondition `json:"conditions,omitempty"`
}

type EventConditionType uint

const (
	_ EventConditionType = iota
	EventConditionTypeMessage
)

type EventCondition struct {
	Type EventConditionType `json:"type"`
	Data json.RawMessage    `json:"data"`
}
