package application

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/pkg/events"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg *events.Event) (*events.Event, error)
}
