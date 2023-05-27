package logic

import "github.com/google/uuid"

type FlowID = uuid.UUID

type Flow struct {
	ID         FlowID     `json:"id"`
	Graph      Graph      `json:"graph"`
	UserValues UserValues `json:"userValues"`
}
