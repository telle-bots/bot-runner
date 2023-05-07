package triggers

import (
	"github.com/mymmrac/telego"
	"github.com/telle-bots/bot-runner/pkg/logic/conditions"
)

type UpdateTrigger struct {
	update telego.Update
}

func NewUpdateTrigger(update telego.Update) *UpdateTrigger {
	return &UpdateTrigger{
		update: update,
	}
}

type MessageTextCondition struct {
	Text conditions.StringCondition `json:"text" name:"Text"`
}

func (t *UpdateTrigger) MessageText(args TriggerArgs) (bool, error) {
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
