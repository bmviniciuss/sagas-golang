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
		Metadata  map[string]string      `json:"metadata"`
	}

	Saga struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		ReplyChannel string    `json:"reply_channel"`
		Step         SagaStep  `json:"step"`
	}

	SagaStep struct {
		ID     uuid.UUID  `json:"id"`
		Name   string     `json:"name"`
		Action ActionType `json:"action"`
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

func NewMessage(
	globalID uuid.UUID,
	eventData map[string]interface{},
	metadata map[string]string,
	workflow *Workflow,
	step *Step,
	action ActionType,
) *Message {
	return &Message{
		EventID:  uuid.New(),
		GlobalID: globalID,
		EventType: EventType{
			SagaName: workflow.Name,
			StepName: step.Name,
			Action:   action,
		},
		Saga: Saga{
			ID:           workflow.ID,
			Name:         workflow.Name,
			ReplyChannel: workflow.ReplyChannel,
			Step: SagaStep{
				ID:     step.ID,
				Name:   step.Name,
				Action: action,
			},
		},
		EventData: eventData,
		Metadata:  metadata,
	}
}
