package flow

import (
	"context"

	"github.com/mymmrac/telego"
	"go.temporal.io/sdk/workflow"

	"github.com/telle-bots/bot-runner/pkg/logic/nodes"

	"github.com/telle-bots/bot-runner/pkg/bot"
)

type Manager struct {
	botRepository *bot.Repository // TODO: Pass
}

func (m *Manager) RetrieveBot(ctx context.Context, botID int64) (*telego.Bot, error) {
	botInfo, err := m.botRepository.BotInfo(ctx, botID)
	if err != nil {
		return nil, err
	}

	return telego.NewBot(botInfo.Token)
}

func (m *Manager) BotWorkflow(ctx workflow.Context, botID int64) error {
	b := &telego.Bot{}
	if err := workflow.ExecuteActivity(ctx, m.RetrieveBot, botID).Get(ctx, b); err != nil {
		return err
	}
	ctx = workflow.WithValue(ctx, nodes.BotContextValue, b)

	selector := workflow.NewSelector(ctx)

	// TODO: Receive flow
	// var flow *logic.Flow

	// ====

	// flow.Graph.Nodes

	updateChan := workflow.GetSignalChannel(ctx, "bot-update")
	selector.AddReceive(updateChan, func(c workflow.ReceiveChannel, _ bool) {
		var update telego.Update
		c.Receive(ctx, &update)

		workflow.ExecuteChildWorkflow(ctx, BotUpdateWorkflow, update)
	})

	// ====
	stop := false

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

func BotUpdateWorkflow(ctx workflow.Context, update telego.Update) error {
	return nil
}
