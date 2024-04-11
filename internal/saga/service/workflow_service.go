package service

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, destination string, data []byte) error
}

type Workflow struct {
	logger    *zap.SugaredLogger
	publisher Publisher
}

func NewWorkflow(
	logger *zap.SugaredLogger,
	publisher Publisher,
) *Workflow {
	return &Workflow{
		logger:    logger,
		publisher: publisher,
	}
}

// TODO: add unit tests
func (w *Workflow) ProcessMessage(ctx context.Context, message *saga.Message, workflow *saga.Workflow) error {
	l := w.logger
	l.Infof("Processing message: %s", message.EventType.Action.String())
	nextStep, err := workflow.GetNextStep(ctx, *message)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while getting next step")
		return err
	}
	if nextStep == nil {
		l.Info("There are no more steps to process. Successfully finished workflow.")
		return nil
	}
	l.Infof("Next step: %s", nextStep.Name)
	nextActionType, err := message.EventType.Action.Next()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while getting next action type")
		return err
	}
	l.Infof("Next step action type: %s", nextStep.Name)
	payload, err := nextStep.PayloadBuilder.Build(ctx, message.EventData, nextActionType)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while building payload")
		return err
	}
	nextMsg := saga.NewMessage(message.GlobalID, payload, message.Metadata, workflow, nextStep, nextActionType)
	jsonMsg, err := json.Marshal(nextMsg)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while marshalling message")
		return err
	}

	err = w.publisher.Publish(ctx, nextStep.DestinationTopic(nextActionType), jsonMsg)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error publishing message to destination")
		return err
	}

	l.Infof("Successfully processed message and produce")
	return nil
}
