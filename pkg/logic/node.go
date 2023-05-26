package logic

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type NodeType uint

const (
	_ NodeType = iota
	NodeTypeEvent
	NodeTypeAction
	NodeTypeCondition
	NodeTypeProcessor
	NodeTypeProvider
)

type NodeID = uuid.UUID

var NilNodeID = uuid.Nil

type NodeFunc[In, Out any] func(ctx context.Context, input *In) (*Out, error)

type Node[In, Out any] struct {
	ID   NodeID            `json:"id"`
	Type NodeType          `json:"type"`
	Func NodeFunc[In, Out] `json:"-"`
}

func (n *Node[In, Out]) MarshalJSON() ([]byte, error) {
	in, err := StructureOf[In]()
	if err != nil {
		return nil, err
	}

	out, err := StructureOf[Out]()
	if err != nil {
		return nil, err
	}

	nodeData := struct {
		ID     NodeID          `json:"id"`
		Type   NodeType        `json:"type"`
		Input  json.RawMessage `json:"input,omitempty"`
		Output json.RawMessage `json:"output,omitempty"`
	}{
		ID:     n.ID,
		Type:   n.Type,
		Input:  in,
		Output: out,
	}

	return json.Marshal(nodeData)
}
