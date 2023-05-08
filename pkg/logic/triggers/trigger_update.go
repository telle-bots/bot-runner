package triggers

import (
	"github.com/mymmrac/telego"
	"github.com/telle-bots/bot-runner/pkg/logic/conditions"
)

type TriggerUpdate struct {
	update telego.Update
}

func NewTriggerUpdate(update telego.Update) *TriggerUpdate {
	return &TriggerUpdate{
		update: update,
	}
}

type MessageTextCondition struct {
	Text conditions.ConditionString `json:"text" name:"Text"`
}

func (t *TriggerUpdate) MessageText(args TriggerArgs) (bool, error) {
	if t.update.Message == nil {
		return false, nil
	}

	condition, err := TriggerConditionAs[MessageTextCondition](args)
	if err != nil {
		return false, err
	}

	if condition.Text.Match(t.update.Message.Text) {
		return true, nil
	}

	return false, nil
}
