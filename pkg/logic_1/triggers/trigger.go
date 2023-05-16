package triggers

import (
	"encoding/json"
	"fmt"
)

type TriggerDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Condition   json.RawMessage `json:"condition,omitempty"`
}

type Trigger func(args TriggerArgs) (bool, error)

type TriggerArgs struct {
	Condition any
	Context   any
}

func TriggerArgsAs[Condition, Context any](args TriggerArgs) (condition Condition, context Context, err error) {
	if condition, err = TriggerConditionAs[Condition](args); err != nil {
		return
	}
	if context, err = TriggerContextAs[Context](args); err != nil {
		return
	}
	return
}

func TriggerConditionAs[Condition any](args TriggerArgs) (Condition, error) {
	condition, ok := args.Condition.(Condition)
	if !ok {
		return condition, fmt.Errorf("invalida condition type %T", condition)
	}
	return condition, nil
}

func TriggerContextAs[Context any](args TriggerArgs) (Context, error) {
	context, ok := args.Context.(Context)
	if !ok {
		return context, fmt.Errorf("invalida cotext type %T", context)
	}
	return context, nil
}
