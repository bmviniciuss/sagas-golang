package saga

import (
	"strings"

	"github.com/google/uuid"
)

type eventTypeSuffix string

const (
	SuccessEventTypeSuffix      eventTypeSuffix = ".success"
	FailureEventTypeSuffix      eventTypeSuffix = ".failure"
	CompensationEventTypeSuffix eventTypeSuffix = ".compensation"
	CompensatedEventTypeSuffix  eventTypeSuffix = ".compensated"
)

type Message struct {
	EventID    string                 `json:"event_id"`
	EventType  string                 `json:"event_type"`
	Metadata   Metadata               `json:"metadata"`
	EventData  map[string]interface{} `json:"event_data"`
	Saga       Saga                   `json:"saga"`
	ActionType ActionType             `json:"-"`
}

type Metadata struct {
	GlobalID string `json:"global_id"`
	ClientID string `json:"client_id"`
}

type Saga struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ReplyChannel string    `json:"reply_channel"`
	Step         SagaStep  `json:"step"`
}

type SagaStep struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (m Message) GetActionType() (ActionType, bool) {
	eventTypeSplit := strings.Split(m.EventID, ".")
	if len(eventTypeSplit) == 0 {
		return ActionType(""), false
	}
	actionStr := eventTypeSplit[len(eventTypeSplit)-1]
	switch actionStr {
	case RequestActionType.String():
		return RequestActionType, true
	case SuccessActionType.String():
		return SuccessActionType, true
	case FailureActionType.String():
		return FailureActionType, true
	case CompensateActionType.String():
		return CompensateActionType, true
	case CompensatedActionType.String():
		return CompensatedActionType, true
	default:
		return ActionType(""), false
	}

}
