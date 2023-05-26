package logic

type Connection struct {
	Source      ConnectionIO `json:"source"`
	Destination ConnectionIO `json:"destination"`
}

type IOType uint

const (
	IOTypeInOut IOType = iota
	IOTypeTrigger
	IOTypeUserValue
)

type ConnectionIO struct {
	NodeID   NodeID `json:"nodeID"`
	Type     IOType `json:"type"`
	DataPath string `json:"dataPath,omitempty"`
}
