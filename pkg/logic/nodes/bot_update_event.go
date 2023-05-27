package nodes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mymmrac/telego"

	"github.com/telle-bots/bot-runner/pkg/logic"
)

const BotUpdateContextValue ContextValue = "bot-update"

type BotUpdate struct {
	UpdateID int64       `json:"updateID"`
	Message  *BotMessage `json:"message,omitempty"`
}

type BotMessage struct {
	ID   int64   `json:"id"`
	Chat BotChat `json:"chat"`
	Text string  `json:"text"`
}

type BotChat struct {
	ID int64 `json:"id"`
}

var BotUpdateEvent = &logic.Node[Empty, BotUpdate]{
	ID:   uuid.MustParse("ab2283b9-2553-474f-9d8c-14267eb026af"),
	Type: logic.NodeTypeEvent,
	Func: func(ctx context.Context, _ *Empty) (*BotUpdate, error) {
		update, ok := ctx.Value(BotUpdateContextValue).(telego.Update)
		if !ok {
			return nil, fmt.Errorf("no context value %q", BotUpdateContextValue)
		}

		upd := &BotUpdate{
			UpdateID: int64(update.UpdateID),
		}

		if update.Message != nil {
			message := update.Message

			upd.Message = &BotMessage{
				ID: int64(message.MessageID),
				Chat: BotChat{
					ID: message.Chat.ID,
				},
				Text: message.Text,
			}
		}

		return upd, nil
	},
}
