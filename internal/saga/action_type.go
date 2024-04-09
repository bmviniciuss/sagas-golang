package saga

import "fmt"

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

func (at ActionType) IsResponseType() bool {
	_, ok := responseActionTypes[at]
	return ok
}

const (
	// requests
	RequestActionType    ActionType = "request"
	CompensateActionType ActionType = "compensate"
	// responses
	SuccessActionType     ActionType = "success"
	FailureActionType     ActionType = "failure"
	CompensatedActionType ActionType = "compensated"
)

func NewActionType(actionType string) (ActionType, error) {
	switch actionType {
	case RequestActionType.String():
		return RequestActionType, nil
	case CompensateActionType.String():
		return CompensateActionType, nil
	case SuccessActionType.String():
		return SuccessActionType, nil
	case FailureActionType.String():
		return FailureActionType, nil
	case CompensatedActionType.String():
		return CompensatedActionType, nil
	default:
		return ActionType(""), fmt.Errorf("unknown action type: %s", actionType)
	}
}

var (
	responseActionTypes = map[ActionType]struct{}{
		SuccessActionType:     {},
		FailureActionType:     {},
		CompensatedActionType: {},
	}
)
