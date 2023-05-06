package logic

import (
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/logic/actions"
)

type ActionManager struct{}

func NewActionManager(_ *do.Injector) (*ActionManager, error) {
	return &ActionManager{}, nil
}

func (a *ActionManager) Actions(botAction *actions.BotAction) map[string]actions.Action {
	return map[string]actions.Action{
		"send_message": {
			Name:      "Send message",
			Do:        botAction.SendMessage,
			Arguments: MustStructureOf[actions.SendMessageData](),
			Returns:   MustStructureOf[actions.SendMessageReturns](),
		},
	}
}
