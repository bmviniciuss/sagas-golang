package streaming

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/internal/saga/service"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type MessageHandler struct {
	logger               *zap.SugaredLogger
	eventTypeWorkflowMap map[string]saga.Workflow
	workflowService      *service.Workflow
}

func NewMessageHandler(logger *zap.SugaredLogger, workflows []saga.Workflow, workflowService *service.Workflow) *MessageHandler {
	eventTypeWorkflowMap := make(map[string]saga.Workflow)
	for _, workflow := range workflows {
		for eventType := range workflow.ConsumerEventTypes() {
			eventTypeWorkflowMap[eventType] = workflow
		}
	}

	return &MessageHandler{
		logger:               logger,
		eventTypeWorkflowMap: eventTypeWorkflowMap,
		workflowService:      workflowService,
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

	// TODO: add idempotence check
	// Get workflow
	workflow, ok := h.eventTypeWorkflowMap[message.EventType.String()]
	if !ok {
		l.Infof("Got unknown event type %s", message.EventType.String())
		err = commitFn()
		if err != nil {
			l.With(zap.Error(err)).Error("Got error committing message")
			return err
		}
		return nil
	}

	err = h.workflowService.ProcessMessage(ctx, &message, &workflow)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error processing workflow message")
		return err
	}

	// TODO: add idempotence set
	err = commitFn()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error committing message")
		return err
	}

	return nil
}
