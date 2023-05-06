package logic

import (
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/logic/actions"
	"github.com/telle-bots/bot-runner/pkg/logic/value"
)

type ActionDo func(args *value.Value) (*value.Value, error)

type Action struct {
	Do        ActionDo     `json:"-"`
	Arguments *value.Value `json:"arguments,omitempty"`
	Returns   *value.Value `json:"returns,omitempty"`
}

type ActionManager struct{}

func NewActionManager(_ *do.Injector) (*ActionManager, error) {
	return &ActionManager{}, nil
}

func (a *ActionManager) Actions(botAction *actions.BotAction) map[string]Action {
	return map[string]Action{
		"send_message": {
			Do: botAction.SendMessage,
			Arguments: &value.Value{
				Name: "Parameters",
				Type: value.Struct,
				Struct: &map[string]value.Value{
					"chat_id": {
						Name: "Chat ID",
						Type: value.Integer,
					},
					"text": {
						Name: "Text",
						Type: value.String,
					},
				},
			},
			Returns: nil,
		},
	}
}
