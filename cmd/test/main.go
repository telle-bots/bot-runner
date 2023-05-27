package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/telle-bots/bot-runner/pkg/logic"
	"github.com/telle-bots/bot-runner/pkg/logic/nodes"
)

func main() {
	updateNode := uuid.New()
	sendNode := uuid.New()

	flow := logic.Flow{
		Graph: logic.Graph{
			Nodes: logic.GraphNodes{
				updateNode: nodes.BotUpdateEvent.ID,
				sendNode:   nodes.BotSendMessageAction.ID,
			},
			Connections: []logic.Connection{
				{
					Source: logic.ConnectionIO{
						NodeID: updateNode,
						Type:   logic.IOTypeTrigger,
					},
					Destination: logic.ConnectionIO{
						NodeID: sendNode,
						Type:   logic.IOTypeTrigger,
					},
				},
				{
					Source: logic.ConnectionIO{
						NodeID:   updateNode,
						DataPath: "message.chat.id",
					},
					Destination: logic.ConnectionIO{
						NodeID:   sendNode,
						DataPath: "chatID",
					},
				},
				{
					Source: logic.ConnectionIO{
						NodeID: logic.NilNodeID,
						Type:   logic.IOTypeUserValue,
					},
					Destination: logic.ConnectionIO{
						NodeID:   sendNode,
						DataPath: "text",
					},
				},
			},
		},
		UserValues: logic.UserValues{
			sendNode: map[string]any{
				"text": "Test text",
			},
		},
	}

	flowData, err := json.Marshal(flow)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(flowData))

	fmt.Println()

	nodeData, err := json.Marshal(nodes.Nodes)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(nodeData))
}
