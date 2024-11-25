package commands

import (
	"fmt"
	"strings"

	"github.com/knackwurstking/tgs/pkg/tgs"
)

var (
	DefaultIPAddress string = "ifconfig.io"
)

type IP struct {
	RequestSendMessage *tgs.RequestSendMessage `json:"-"`

	Address string `json:"address"`

	URL *string `json:"-"` // Defaults to "ifconfig.io"
}

func NewIP(api tgs.API, url *string) *IP {
	if url == nil {
		url = &DefaultIPAddress
	}

	return &IP{
		RequestSendMessage: nil,
		URL:                url,
	}
}

func (ip *IP) Run() error {
	if ip.RequestSendMessage == nil {
		return fmt.Errorf("missing sendMessage request")
	}

	if ip.URL == nil {
		ip.URL = &DefaultIPAddress
	}

	address, err := ip.fetchAddressFromURL()
	if err != nil {
		return err
	}

	// TODO: Send response with address back to client, need a new request type here
	_, err = ip.RequestSendMessage.Send()
	if err != nil {
		return err
	}

	return fmt.Errorf("under construction")
}

func (ip *IP) fetchAddressFromURL() (string, error) {
	if ip.URL == nil {
		ip.URL = &DefaultIPAddress
	}

	var ret string

	switch *ip.URL {
	case DefaultIPAddress:
		// TODO: Fetch ip address from this URL
		break
	}

	return strings.Trim(ret, "\n\t "), fmt.Errorf("under construction")
}
