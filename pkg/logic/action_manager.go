package logic

import (
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/logic/actions"
)

type ActionName string

const (
	ActionSendMessage ActionName = "sendMessage"
)

var Actions = map[ActionName]actions.Action{
	ActionSendMessage: {
		Name:      "Send message",
		Arguments: MustStructureOf[actions.SendMessageData](),
		Returns:   MustStructureOf[actions.SendMessageReturns](),
	},
}

type ActionManager struct{}

func NewActionManager(_ *do.Injector) (*ActionManager, error) {
	return &ActionManager{}, nil
}

func (a *ActionManager) Actions(botAction *actions.BotAction) map[ActionName]actions.ActionDo {
	return map[ActionName]actions.ActionDo{
		ActionSendMessage: botAction.SendMessage,
	}
}
