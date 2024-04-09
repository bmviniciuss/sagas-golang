package streaming

import (
	"strings"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
)

type EventTypeSuffix string

const (
	SuccessEventTypeSuffix      EventTypeSuffix = ".success"
	FailureEventTypeSuffix      EventTypeSuffix = ".failure"
	CompensationEventTypeSuffix EventTypeSuffix = ".compensation"
	CompensatedEventTypeSuffix  EventTypeSuffix = ".compensated"
)

type Message struct {
	Metadata  Metadata               `json:"metadata"`
	EventID   string                 `json:"event_id"`
	EventType string                 `json:"event_type"`
	EventData map[string]interface{} `json:"event_data"`
	Saga      Saga                   `json:"saga"`
}

func (m Message) GetActionType() (saga.ActionType, bool) {
	eventTypeSplit := strings.Split(m.EventID, ".")
	if len(eventTypeSplit) == 0 {
		return saga.ActionType(""), false
	}
	actionStr := eventTypeSplit[len(eventTypeSplit)-1]
	switch actionStr {
	case saga.RequestActionType.String():
		return saga.RequestActionType, true
	case saga.SuccessActionType.String():
		return saga.SuccessActionType, true
	case saga.FailureActionType.String():
		return saga.FailureActionType, true
	case saga.CompensateActionType.String():
		return saga.CompensateActionType, true
	case saga.CompensatedActionType.String():
		return saga.CompensatedActionType, true
	default:
		return saga.ActionType(""), false
	}

}

type Metadata struct {
	GlobalID string `json:"global_id"`
	ClientID string `json:"client_id"`
}

type Saga struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	ReplyChannel string                 `json:"reply_channel"`
	Step         SagaStep               `json:"step"`
	Metadata     map[string]interface{} `json:"saga"`
}

type SagaStep struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
