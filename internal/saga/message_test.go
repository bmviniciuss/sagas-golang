package saga

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventType_String(t *testing.T) {
	type fields struct {
		SagaName string
		StepName string
		Action   ActionType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should return string if event type",
			fields: fields{
				SagaName: "saga",
				StepName: "step",
				Action:   SuccessActionType,
			},
			want: "saga.step.success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EventType{
				SagaName: tt.fields.SagaName,
				StepName: tt.fields.StepName,
				Action:   tt.fields.Action,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("EventType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventType_MarshalJSON(t *testing.T) {
	type fields struct {
		SagaName string
		StepName string
		Action   ActionType
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "should return json bytes of event type",
			fields: fields{
				SagaName: "saga",
				StepName: "step",
				Action:   SuccessActionType,
			},
			want:    []byte(`"saga.step.success"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EventType{
				SagaName: tt.fields.SagaName,
				StepName: tt.fields.StepName,
				Action:   tt.fields.Action,
			}
			got, err := e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("EventType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventType_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		compare EventType
	}{
		{
			name: "should return err if event type is not quoted",
			args: args{
				data: []byte(`saga.step.success`),
			},
			wantErr: true,
		},
		{
			name: "should return err if event type is empty",
			args: args{
				data: []byte(`""`),
			},
			wantErr: true,
		},
		{
			name: "should return err if event type is malformed",
			args: args{
				data: []byte(`"saga.step"`),
			},
			wantErr: true,
		},
		{
			name: "should return err if event action type is unknown",
			args: args{
				data: []byte(`"saga.step.unknown"`),
			},
			wantErr: true,
		},
		{
			name: "should return nil and parse event type",
			args: args{
				data: []byte(`"saga.step.success"`),
			},
			wantErr: false,
			compare: EventType{
				SagaName: "saga",
				StepName: "step",
				Action:   SuccessActionType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event EventType
			err := event.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Nil(t, err)
			assert.Equal(t, tt.compare, event)
		})
	}
}
