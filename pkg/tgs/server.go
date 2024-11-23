package tgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	Token string `json:"-"`
}

func (s *Server) SetToken(token string) {
	s.Token = token
}

func (s *Server) URL(command Command) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", s.Token, command)
}

func (s *Server) Send(request Request) ([]byte, error) {
	if s.Token == "" {
		return nil, fmt.Errorf("missing token")
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("GET", s.URL(request.Command()), body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)

}
