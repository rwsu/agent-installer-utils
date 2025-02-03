package ui

import (
	"fmt"
	
	"github.com/openshift/agent-installer-utils/tools/agent_tui/newt"
	"github.com/rivo/tview"
)

const (
	OK_BUTTON               string = "<Ok>"
	PAGE_INVALID_IP_ADDRESS        = "invalidIPAddress"

	invalidIPText = "The IP address %s is not a valid IPv4 or IPv6 address."
)

// Creates the invalid IP address modal but does not add the modal
// to pages. The rendezvousIPForm does that when it validates the IP
// address entered.
func (u *UI) createInvalidIPAddressModal() {
	// view is the modal asking the user if they would still
	// like to change their network configuration.
	u.invalidIPAddressModal = tview.NewModal().
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == OK_BUTTON {
				u.setFocusToRendezvousIP()
			}
		}).
		SetBackgroundColor(newt.ColorBlack)
	u.invalidIPAddressModal.
		SetBorderColor(newt.ColorBlack).
		SetBorder(true)
	u.invalidIPAddressModal.
		SetButtonBackgroundColor(newt.ColorGray).
		SetButtonTextColor(newt.ColorRed)
	userPromptButtons := []string{OK_BUTTON}
	u.invalidIPAddressModal.AddButtons(userPromptButtons)

	//u.invalidIPAddressModal.SetText(invalidIPText)
	u.pages.AddPage(PAGE_INVALID_IP_ADDRESS, u.invalidIPAddressModal, true, false)
}

func (u *UI) ShowInvalidIPAddressDialog(invalidIPAddress string) {
	//u.setIsTimeoutDialogActive(true)
	u.invalidIPAddressModal.SetText(fmt.Sprintf(invalidIPText, invalidIPAddress))
	u.app.SetFocus(u.invalidIPAddressModal)
	u.pages.ShowPage(PAGE_INVALID_IP_ADDRESS)
}
