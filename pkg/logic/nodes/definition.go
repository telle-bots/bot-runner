package nodes

import (
	"encoding/json"

	"github.com/mymmrac/telego"
	"github.com/telle-bots/bot-runner/pkg/logic"
	"go.temporal.io/sdk/workflow"
)

type Empty struct{}

type ContextValue string

var Nodes = map[logic.NodeID]json.Marshaler{
	BotUpdateEvent.ID:       BotUpdateEvent,
	BotSendMessageAction.ID: BotSendMessageAction,
}

// =======

func BotWorkflow(ctx workflow.Context, botID int64) error {
	stop := false
	selector := workflow.NewSelector(ctx)

	updateChan := workflow.GetSignalChannel(ctx, "bot-update-event")
	selector.AddReceive(updateChan, func(c workflow.ReceiveChannel, _ bool) {
		var update telego.Update
		c.Receive(ctx, &update)

		workflow.ExecuteChildWorkflow(ctx, BotUpdateWorkflow, botID, update)
	})

	disableChan := workflow.GetSignalChannel(ctx, "bot-disable")
	selector.AddReceive(disableChan, func(c workflow.ReceiveChannel, _ bool) {
		var empty struct{}
		c.Receive(ctx, &empty)

		stop = true
	})

	for !stop {
		selector.Select(ctx)
	}

	return nil
}

func BotUpdateWorkflow(ctx workflow.Context, botID int64, update telego.Update) error {
	return nil
}
