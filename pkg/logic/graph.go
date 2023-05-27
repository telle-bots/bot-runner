package logic

type Graph struct {
	Nodes       GraphNodes   `json:"nodes"`
	Connections []Connection `json:"connections"`
}

type GraphNodes map[NodeID]NodeID

type UserValues map[NodeID]map[string]any
