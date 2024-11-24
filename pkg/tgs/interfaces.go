package tgs

type Request interface {
	Command() Command
}

type API interface {
	Token() string
	SetToken(token string)
	URL(command Command) string
	SendRequest(request Request) ([]byte, error)
}
