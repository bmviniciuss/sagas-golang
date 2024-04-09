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

const (
	RequestActionType      ActionType = "request"
	CompensationActionType ActionType = "compensation"
	SuccessActionType      ActionType = "success"
	FailureActionType      ActionType = "failure"
	CompensatedActionType  ActionType = "compensated"
)
