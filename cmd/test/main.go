package main

import (
	"encoding/json"
	"fmt"

	"github.com/telle-bots/bot-runner/pkg/logic_2"
)

func main() {
	workflows := logic_2.Workflows{
		{
			Name: "Ping - Time - Pong",
			Triggers: []logic_2.Trigger{
				{
					Type: logic_2.TriggerTypeBot,
					Events: []logic_2.Event{
						{
							Type: logic_2.EventTypeUpdate,
							Conditions: []logic_2.EventCondition{
								{
									Type: logic_2.EventConditionTypeMessage,
									Data: []byte(`{"field": "update.message.text", "equal": "Ping"}`),
								},
							},
						},
					},
				},
			},
			Conditions: []logic_2.ConditionSource{
				{
					Type: logic_2.ConditionSourceTypeClock,
					Conditions: []logic_2.SourceCondition{
						{
							Type: logic_2.SourceConditionTypeTime,
							Data: []byte(`{"field": "now.minutes", "even": true}`),
						},
					},
				},
			},
			Actions: []logic_2.ActionSource{
				{
					Type: logic_2.ActionSourceTypeBot,
					Actions: []logic_2.Action{
						{
							Type: logic_2.ActionTypeSendMessage,
							Data: []byte(`{"text": "Pong", "chat_id": "update.message.chat.id"}`),
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(workflows)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
