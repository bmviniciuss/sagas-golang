package saga

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type EventType struct {
	SagaName string
	StepName string
	Action   ActionType
}

func (e EventType) String() string {
	return fmt.Sprintf("%s.%s.%s", e.SagaName, e.StepName, e.Action.String())
}

// TODO: rename to NewResponseEventType
func NewReplyEventType(ev EventType, action ActionType) EventType {
	return EventType{
		SagaName: ev.SagaName,
		StepName: ev.StepName,
		Action:   action,
	}
}

func (e EventType) MarshalJSON() ([]byte, error) {
	val := e.String()
	return []byte(strconv.Quote(val)), nil
}

func (e *EventType) UnmarshalJSON(data []byte) error {
	val, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	if val == "" {
		return fmt.Errorf("got empty event type")
	}
	split := strings.Split(val, ".")
	if len(split) != 3 {
		return fmt.Errorf("event type malformed")
	}
	e.SagaName = split[0]
	e.StepName = split[1]
	actionType, err := NewActionType(split[2])
	if err != nil {
		return err
	}
	e.Action = actionType
	return nil
}

type (
	Message struct {
		GlobalID  uuid.UUID              `json:"global_id"`
		EventID   uuid.UUID              `json:"event_id"`
		EventType EventType              `json:"event_type"`
		Saga      Saga                   `json:"saga"`
		EventData map[string]interface{} `json:"event_data"`
		Metadata  map[string]string      `json:"metadata,omitempty"`
	}
	Saga struct {
		ReplyChannel string `json:"reply_channel"`
	}
)

func (m *Message) Hash() (string, error) {
	dataBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	sha256 := sha256.New()
	hash := sha256.Sum(dataBytes)
	return fmt.Sprintf("%x", hash), nil
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func NewMessage(
	globalID uuid.UUID,
	eventData map[string]interface{},
	metadata map[string]string,
	workflow *Workflow,
	step *Step,
	action ActionType,
) *Message {
	data := make(map[string]interface{})
	if eventData != nil {
		data = eventData
	}

	return &Message{
		EventID:  uuid.New(),
		GlobalID: globalID,
		EventType: EventType{
			SagaName: workflow.Name,
			StepName: step.Name,
			Action:   action,
		},
		Saga: Saga{
			ReplyChannel: workflow.ReplyChannel,
		},
		EventData: data,
		Metadata:  metadata,
	}
}

func NewParticipantMessage(
	globalID uuid.UUID,
	eventData map[string]interface{},
	metadata map[string]string,
	eventType EventType,
	replyChannel string,
) *Message {
	data := make(map[string]interface{})
	if eventData != nil {
		data = eventData
	}

	return &Message{
		EventID:   uuid.New(),
		GlobalID:  globalID,
		EventType: eventType,
		Saga: Saga{
			ReplyChannel: replyChannel,
		},
		EventData: data,
		Metadata:  metadata,
	}
}
