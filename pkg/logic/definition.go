package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.temporal.io/sdk/workflow"
)

type Empty struct{}

type ContextValue string

const (
	BotContextValue       ContextValue = "bot"
	BotUpdateContextValue ContextValue = "bot-update"
)

var Nodes = map[NodeID]json.Marshaler{
	BotUpdateEvent.ID:       BotUpdateEvent,
	BotSendMessageAction.ID: BotSendMessageAction,
}

type UpdateEventOut struct {
	UpdateID int                    `json:"updateID"`
	Message  *UpdateEventOutMessage `json:"message,omitempty"`
}

type UpdateEventOutMessage struct {
	ID   int                `json:"id"`
	Chat UpdateEventOutChat `json:"chat"`
	Text string             `json:"text"`
}

type UpdateEventOutChat struct {
	ID int64 `json:"id"`
}

var BotUpdateEvent = &Node[Empty, UpdateEventOut]{
	ID:   uuid.MustParse("ab2283b9-2553-474f-9d8c-14267eb026af"),
	Type: NodeTypeEvent,
	Func: func(ctx context.Context, _ *Empty) (*UpdateEventOut, error) {
		update, ok := ctx.Value(BotUpdateContextValue).(telego.Update)
		if !ok {
			return nil, fmt.Errorf("no context value %q", BotUpdateContextValue)
		}

		upd := &UpdateEventOut{
			UpdateID: update.UpdateID,
		}

		if update.Message != nil {
			message := update.Message

			upd.Message = &UpdateEventOutMessage{
				ID: message.MessageID,
				Chat: UpdateEventOutChat{
					ID: message.Chat.ID,
				},
				Text: message.Text,
			}
		}

		return upd, nil
	},
}

type BotSendMessageActionIn struct {
	ChatID int64  `json:"chatID"`
	Text   string `json:"text"`
}

var BotSendMessageAction = &Node[BotSendMessageActionIn, Empty]{
	ID:   uuid.MustParse("7b42a5e7-0c59-4034-94f9-8310c531e521"),
	Type: NodeTypeAction,
	Func: func(ctx context.Context, input *BotSendMessageActionIn) (*Empty, error) {
		bot, ok := ctx.Value(BotContextValue).(*telego.Bot)
		if !ok {
			return nil, fmt.Errorf("no context value %q", BotContextValue)
		}

		_, err := bot.SendMessage(tu.Message(tu.ID(input.ChatID), input.Text))
		if err != nil {
			return nil, err
		}

		return &Empty{}, nil
	},
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
