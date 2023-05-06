package actions

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type BotAction struct {
	bot *telego.Bot
}

func NewBotAction(bot *telego.Bot) *BotAction {
	return &BotAction{
		bot: bot,
	}
}

type SendMessageContext struct {
	ChatID int64
}

type SendMessageData struct {
	Text string `json:"text" name:"Message text"`
}

type SendMessageReturns struct {
	MessageID int64 `json:"messageID" name:"Message ID"`
}

func (b *BotAction) SendMessage(args ActionArgs) (any, error) {
	data, context, err := ActionArgsAs[SendMessageData, SendMessageContext](args)
	if err != nil {
		return nil, err
	}

	msg, err := b.bot.SendMessage(&telego.SendMessageParams{
		ChatID: tu.ID(context.ChatID),
		Text:   data.Text,
	})

	return &SendMessageReturns{
		MessageID: int64(msg.MessageID),
	}, err
}
