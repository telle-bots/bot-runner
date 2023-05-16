package logic_2

import (
	"fmt"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

func HandleTrigger(event Event, eventData any) (bool, error) {
	return false, nil
}

func HandleCondition(condition SourceCondition) (bool, error) {
	return false, nil
}

func HandleAction(action Action, eventType EventType, eventData any) error {
	switch action.Type {
	case ActionTypeSendMessage:
		switch eventType {
		case EventTypeUpdate:

		default:
			return fmt.Errorf("unsupported event type in action %q: %q", action.Type, eventType)
		}
	default:
		return fmt.Errorf("unsupported action type: %q", action.Type)
	}

	return nil
}

func ProcessWorkflows(
	log *zap.SugaredLogger, wfs Workflows, triggerTypes []TriggerType, eventTypes []EventType, eventData any,
) {
	for _, wf := range wfs {
		for _, tr := range wf.Triggers {
			if !lo.Contains(triggerTypes, tr.Type) {
				continue
			}

			for _, ev := range tr.Events {
				if !lo.Contains(eventTypes, ev.Type) {
					continue
				}

				var triggered bool
				triggered, err := HandleTrigger(ev, eventData)
				if err != nil {
					log.Errorw("Handle trigger", "error", err, "trigger", tr.Type, "event", ev.Type)
					continue
				}
				if !triggered {
					continue
				}

				var pass bool
				conditionsPass := true
				for _, cds := range wf.Conditions {
					if !conditionsPass {
						break
					}
					for _, cd := range cds.Conditions {
						pass, err = HandleCondition(cd)
						if err != nil {
							log.Errorw("Handle condition", "error", err, "trigger", tr.Type,
								"event", ev.Type, "condition-source", cds.Type, "condition", cd.Type)
							conditionsPass = false
							break
						}
						if !pass {
							conditionsPass = false
							break
						}
					}
				}
				if !conditionsPass {
					continue
				}

				actionPass := true
				for _, acs := range wf.Actions {
					if !actionPass {
						return
					}
					for _, ac := range acs.Actions {
						if err = HandleAction(ac, ev.Type, eventData); err != nil {
							log.Errorw("Handle action", "error", err, "trigger", tr.Type,
								"event", ev.Type, "action-source", acs.Type, "action", ac.Type)
							actionPass = false
							break
						}
					}
				}
			}
		}
	}
}
