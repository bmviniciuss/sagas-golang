package saga

type ActionType string

func (at ActionType) String() string {
	return string(at)
}

func (at ActionType) IsSuccess() bool {
	return at == SuccessActionType
}

func (at ActionType) IsFailure() bool {
	return at == FailureActionType
}

func (at ActionType) IsCompensated() bool {
	return at == CompensatedActionType
}

// func (at ActionType) IsResponseType() bool {

// }

const (
	// requests
	RequestActionType    ActionType = "request"
	CompensateActionType ActionType = "compensate"
	// responses
	SuccessActionType     ActionType = "success"
	FailureActionType     ActionType = "failure"
	CompensatedActionType ActionType = "compensated"
)

var (
	consumerActionTypes = map[ActionType]struct{}{
		SuccessActionType:     struct{}{},
		FailureActionType:     struct{}{},
		CompensatedActionType: struct{}{},
	}
)
