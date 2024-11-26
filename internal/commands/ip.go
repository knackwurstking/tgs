package commands

import (
	"fmt"
	"io"
	"net/http"
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
}

func NewIP(api tgs.API) *IP {
	return &IP{
		RequestSendMessage: nil,
	}
}

func (this *IP) Run(chatID int) error {
	if this.RequestSendMessage == nil {
		return fmt.Errorf("missing sendMessage request")
	}

	address, err := this.fetchAddressFromURL()
	if err != nil {
		return err
	}

	this.RequestSendMessage.ParseMode = data.ParseModeMarkdownV2
	this.RequestSendMessage.Text = fmt.Sprintf("`%s`", address)
	this.RequestSendMessage.ChatID = chatID

	_, err = this.RequestSendMessage.Send()
	if err != nil {
		return err
	}

	return nil
}

func (this *IP) fetchAddressFromURL() (address string, err error) {
	resp, err := http.Get(defaultIPAddress)
	if err != nil {
		return address, err
	}
	if resp.StatusCode != http.StatusOK {
		return address, fmt.Errorf("request to %s: %d (%s)",
			defaultIPAddress, resp.StatusCode, resp.Status,
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return address, err
	}

	return strings.Trim(string(data), "\n\r\t "), nil
}
