package logic_2

type Workflows []Workflow

type Workflow struct {
	Name string `json:"name"`

	Triggers   []Trigger         `json:"triggers"`
	Conditions []ConditionSource `json:"conditions,omitempty"`
	Actions    []ActionSource    `json:"actions"`
}
