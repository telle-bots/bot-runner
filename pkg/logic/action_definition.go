package logic

import "github.com/telle-bots/bot-runner/pkg/logic/actions"

type ActionName string

const (
	ActionSendMessage ActionName = "sendMessage"
)

var ActionDefinitions = map[ActionName]actions.ActionDefinition{
	ActionSendMessage: {
		Name:      "Send message",
		Arguments: MustStructureOf[actions.SendMessageData](),
		Returns:   MustStructureOf[actions.SendMessageReturns](),
	},
}

func Actions(actionBot *actions.ActionBot) map[ActionName]actions.Action {
	return map[ActionName]actions.Action{
		ActionSendMessage: actionBot.SendMessage,
	}
}
