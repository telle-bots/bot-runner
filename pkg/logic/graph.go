package logic

type Graph struct {
	Nodes       []NodeID     `json:"nodes"`
	Connections []Connection `json:"connections"`
}

type UserValues map[NodeID]map[string]any
