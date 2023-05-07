package actions

import (
	"encoding/json"
	"fmt"
)

type ActionDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Arguments   json.RawMessage `json:"arguments,omitempty"`
	Returns     json.RawMessage `json:"returns,omitempty"`
}

type Action func(args ActionArgs) (any, error)

type ActionArgs struct {
	Data    any
	Context any
}

func ActionArgsAs[Data, Context any](args ActionArgs) (data Data, context Context, err error) {
	if data, err = ActionDataAs[Data](args); err != nil {
		return
	}
	if context, err = ActionContextAs[Context](args); err != nil {
		return
	}
	return
}

func ActionDataAs[Data any](args ActionArgs) (Data, error) {
	data, ok := args.Data.(Data)
	if !ok {
		return data, fmt.Errorf("invalida data type %T", data)
	}
	return data, nil
}

func ActionContextAs[Context any](args ActionArgs) (Context, error) {
	context, ok := args.Context.(Context)
	if !ok {
		return context, fmt.Errorf("invalida cotext type %T", context)
	}
	return context, nil
}
