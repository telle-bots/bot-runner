package actions

import "fmt"

type Action struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Do          ActionDo `json:"-"`
	Arguments   []byte   `json:"arguments,omitempty"`
	Returns     []byte   `json:"returns,omitempty"`
}

type ActionDo func(args ActionArgs) (any, error)

type ActionArgs struct {
	Data    any
	Context any
}

func ActionArgsAs[D, C any](args ActionArgs) (data D, context C, err error) {
	if data, err = ActionDataAs[D](args); err != nil {
		return
	}
	if context, err = ActionContextAs[C](args); err != nil {
		return
	}
	return
}

func ActionDataAs[T any](args ActionArgs) (T, error) {
	data, ok := args.Data.(T)
	if !ok {
		return data, fmt.Errorf("invalida data type %T", data)
	}
	return data, nil
}

func ActionContextAs[T any](args ActionArgs) (T, error) {
	context, ok := args.Context.(T)
	if !ok {
		return context, fmt.Errorf("invalida cotext type %T", context)
	}
	return context, nil
}
