package value

type Type uint

const (
	_ Type = iota
	Struct
	Array
	Map
	String
	Integer
	Float
)

type Value struct {
	Name     string `json:"name"`
	Type     Type   `json:"type"`
	Optional bool   `json:"optional,omitempty"`

	Struct  *map[string]Value `json:"struct,omitempty"`
	Array   *[]Value          `json:"array,omitempty"`
	Map     *map[Value]Value  `json:"map,omitempty"`
	String  *string           `json:"string,omitempty"`
	Integer *int64            `json:"integer,omitempty"`
	Float   *float64          `json:"float,omitempty"`
}

func (v Value) GetStruct() map[string]Value {
	if v.Struct == nil {
		return nil
	}

	return *v.Struct
}

func (v Value) GetString() string {
	if v.String == nil {
		return ""
	}

	return *v.String
}

func (v Value) GetInteger() int64 {
	if v.Integer == nil {
		return 0
	}

	return *v.Integer
}
