package logic

type Connection struct {
	Source      ConnectionIO `json:"source"`
	Destination ConnectionIO `json:"destination"`
}

type ConnectionIO struct {
	NodeID   NodeID `json:"nodeID"`
	DataPath string `json:"dataPath"`
}
