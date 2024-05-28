package saga

const (
	REQUEST_ACTION_TYPE             ActionType = "request"
	COMPESATION_REQUEST_ACTION_TYPE ActionType = "compensation_request"
)

type ActionType string

func (at ActionType) String() string {
	return string(at)
}

func (at ActionType) IsRequest() bool {
	return at == REQUEST_ACTION_TYPE
}

func (at ActionType) IsCompensationRequest() bool {
	return at == COMPESATION_REQUEST_ACTION_TYPE
}
