package main

import (
	"encoding/json"
	"fmt"

	"github.com/telle-bots/bot-runner/pkg/logic"
)

func main() {
	flow := logic.Flow{
		Graph: logic.Graph{
			Nodes: []logic.NodeID{
				logic.BotUpdateEvent.ID,
				logic.BotSendMessageAction.ID,
			},
			Connections: []logic.Connection{
				{
					Source: logic.ConnectionIO{
						NodeID:   logic.BotUpdateEvent.ID,
						DataPath: "output.message.chat.id",
					},
					Destination: logic.ConnectionIO{
						NodeID:   logic.BotSendMessageAction.ID,
						DataPath: "input.chatID",
					},
				},
				{
					Source: logic.ConnectionIO{
						NodeID:   logic.NilNodeID,
						DataPath: "userValue",
					},
					Destination: logic.ConnectionIO{
						NodeID:   logic.BotSendMessageAction.ID,
						DataPath: "input.text",
					},
				},
			},
		},
		UserValues: logic.UserValues{
			logic.BotSendMessageAction.ID: map[string]any{
				"input.text": "Test text",
			},
		},
	}

	flowData, err := json.Marshal(flow)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(flowData))

	fmt.Println()

	nodeData, err := json.Marshal(logic.Nodes)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(nodeData))
}
