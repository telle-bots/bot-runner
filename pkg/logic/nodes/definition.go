package nodes

import (
	"encoding/json"

	"github.com/telle-bots/bot-runner/pkg/logic"
)

type Empty struct{}

type ContextValue string

var Nodes = map[logic.NodeID]json.Marshaler{
	BotUpdateEvent.ID:       BotUpdateEvent,
	BotSendMessageAction.ID: BotSendMessageAction,
}
