package logic_2

import "encoding/json"

type ActionSourceType uint

const (
	_ ActionSourceType = iota
	ActionSourceTypeBot
)

type ActionSource struct {
	Type    ActionSourceType `json:"type"`
	Actions []Action         `json:"actions"`
}

type ActionType uint

const (
	_ ActionType = iota
	ActionTypeSendMessage
)

type Action struct {
	Type ActionType      `json:"type"`
	Data json.RawMessage `json:"data"`
}
