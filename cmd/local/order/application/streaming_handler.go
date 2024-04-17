package application

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg *saga.Message) (*saga.Message, error)
}
