package events

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/bmviniciuss/sagas-golang/pkg/utc"
	"github.com/google/uuid"
)

type Event struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Origin        string                 `json:"origin"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	Date          string                 `json:"date"`
	Data          map[string]interface{} `json:"data"`
}

func NewEvent(eventType string, origin string, data map[string]interface{}) *Event {
	now := utc.Now()
	payload := make(map[string]interface{})
	if data != nil {
		payload = data
	}
	return &Event{
		ID:            uuid.NewString(),
		Type:          eventType,
		Origin:        origin,
		CorrelationID: uuid.NewString(),
		Date:          now.Time().Format(utc.ISO8601Layout),
		Data:          payload,
	}
}

func (e *Event) WithCorrelationID(correlationID string) *Event {
	e.CorrelationID = correlationID
	return e
}

func (m *Event) Hash() (string, error) {
	dataBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	sha256 := sha256.New()
	hash := sha256.Sum(dataBytes)
	return fmt.Sprintf("%x", hash), nil
}

func (m *Event) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
