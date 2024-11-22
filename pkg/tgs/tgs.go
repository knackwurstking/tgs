package tgs

import "fmt"

type Requests interface {
	*RequestGetMe | *RequestGetUpdates
}

type Responses interface {
	*ResponseGetMe | *ResponseGetUpdates
}

func Get[Request Requests, Response Responses](request Request) (Response, error) {
	if request == nil {
		return nil, fmt.Errorf("missing request")
	}

	// TODO: ...

	return nil, nil
}

func Post[Request Requests, Response Responses](request Request) (Response, error) {
	if request == nil {
		return nil, fmt.Errorf("missing request")
	}

	// TODO: ...

	return nil, nil
}
