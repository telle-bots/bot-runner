package actions

import (
	"errors"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/telle-bots/bot-runner/pkg/logic/value"
)

type BotAction struct {
	bot *telego.Bot
}

func NewBotAction(bot *telego.Bot) *BotAction {
	return &BotAction{
		bot: bot,
	}
}

func (b *BotAction) SendMessage(args *value.Value) (*value.Value, error) {
	if args == nil || args.Type != value.Struct || args.Struct == nil {
		return nil, errors.New("invalid argument")
	}

	_, err := b.bot.SendMessage(&telego.SendMessageParams{
		ChatID: tu.ID(args.GetStruct()["chat_id"].GetInteger()),
		Text:   args.GetStruct()["text"].GetString(),
	})

	return nil, err
}
