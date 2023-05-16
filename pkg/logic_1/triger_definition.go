package logic_1

import "github.com/telle-bots/bot-runner/pkg/logic_1/triggers"

type TriggerName string

const (
	TriggerMessageText TriggerName = "messageText"
)

var TriggerDefinitions = map[TriggerName]triggers.TriggerDefinition{
	TriggerMessageText: {
		Name:      "Message text",
		Condition: MustStructureOf[triggers.MessageTextCondition](),
	},
}

func Triggers(triggerUpdate *triggers.TriggerUpdate) map[TriggerName]triggers.Trigger {
	return map[TriggerName]triggers.Trigger{
		TriggerMessageText: triggerUpdate.MessageText,
	}
}
