package nodes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/telle-bots/bot-runner/pkg/logic"
)

const BotContextValue ContextValue = "bot"

type BotSendMessageArgs struct {
	ChatID int64  `json:"chatID"`
	Text   string `json:"text"`
}

var BotSendMessageAction = &logic.Node[BotSendMessageArgs, Empty]{
	ID:   uuid.MustParse("7b42a5e7-0c59-4034-94f9-8310c531e521"),
	Type: logic.NodeTypeAction,
	Func: func(ctx context.Context, input *BotSendMessageArgs) (*Empty, error) {
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
