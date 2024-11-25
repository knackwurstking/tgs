package commands

import (
	"fmt"
	"strings"

	"github.com/knackwurstking/tgs/pkg/data"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

var (
	defaultIPAddress string = "ifconfig.io"
)

type IP struct {
	RequestSendMessage *tgs.RequestSendMessage `json:"-"`

	Address string `json:"address"`

	URL *string `json:"-"` // Defaults to "ifconfig.io"
}

func NewIP(api tgs.API, url *string) *IP {
	if url == nil {
		url = &defaultIPAddress
	}

	return &IP{
		RequestSendMessage: nil,
		URL:                url,
	}
}

func (ip *IP) Run(chatID int) error {
	if ip.RequestSendMessage == nil {
		return fmt.Errorf("missing sendMessage request")
	}

	address, err := ip.fetchAddressFromURL()
	if err != nil {
		return err
	}

	ip.RequestSendMessage.ParseMode = data.ParseModeMarkdownV2
	ip.RequestSendMessage.Text = fmt.Sprintf("`%s`", address)
	ip.RequestSendMessage.ChatID = chatID

	_, err = ip.RequestSendMessage.Send()
	if err != nil {
		return err
	}

	return fmt.Errorf("under construction")
}

func (ip *IP) fetchAddressFromURL() (string, error) {
	var ret string

	switch *ip.URL {
	case defaultIPAddress:
		// TODO: Fetch ip address from this URL
		break
	}

	return strings.Trim(ret, "\n\t "), fmt.Errorf("under construction")
}
