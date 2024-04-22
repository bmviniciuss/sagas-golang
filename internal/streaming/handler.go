package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/saga/service"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type IdempotenceService interface {
	Has(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, ttl time.Duration) error
}

type MessageHandler struct {
	logger              *zap.SugaredLogger
	executionRepository saga.ExecutionRepository
	workflowService     *service.Execution
	idempotenceService  IdempotenceService
}

func NewMessageHandler(
	logger *zap.SugaredLogger,
	executionRepository saga.ExecutionRepository,
	workflowService *service.Execution,
	idempotenceService IdempotenceService,
) *MessageHandler {
	return &MessageHandler{
		logger:              logger,
		executionRepository: executionRepository,
		workflowService:     workflowService,
		idempotenceService:  idempotenceService,
	}
}

func (h *MessageHandler) Handle(ctx context.Context, msg *kafka.Message, commitFn func() error) error {
	l := h.logger
	l.Info("Handling message")
	var message saga.Message
	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error unmarshalling message")
		return err
	}
	l.With("message", message).Info("Got message")
	msgHash, err := message.Hash()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error creating message hash")
		return err
	}
	key := fmt.Sprintf("%s:%s:%s", message.GlobalID, message.EventID, msgHash)
	idempotent, err := h.idempotenceService.Has(ctx, key)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error checking idempotence")
	}

	if idempotent {
		l.Info("Message was already processed")
		err = commitFn()
		if err != nil {
			l.With(zap.Error(err)).Error("Got error committing message")
			return err
		}
		return nil
	}

	// Get execution
	execution, err := h.executionRepository.Find(ctx, message.Saga.ID.String())
	if err != nil {
		l.With(zap.Error(err)).Error("Got error getting workflow")
		return err // TODO: handle error
	}

	if execution.IsEmpty() {
		l.Info("execution not found. Message will be ignored")
		err = commitFn()
		if err != nil {
			l.With(zap.Error(err)).Error("Got error committing message")
			return err
		}
		return nil
	}

	err = h.workflowService.ProcessMessage(ctx, &message, execution)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error processing workflow message")
		return err
	}

	err = h.idempotenceService.Set(ctx, key, time.Hour*24*30)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error setting idempotence")
	}
	err = commitFn()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error committing message")
		return err
	}

	return nil
}
