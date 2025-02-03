package ui

import (
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/checks"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/newt"
	"github.com/rivo/tview"
)

const (
	PAGE_RENDEZVOUS_IP        = "rendezvousIPScreen"
	FIELD_CURRENT_HOST_IP     = "Current Host IP"
	FIELD_RENDEZVOUS_HOST_IP  = "Rendezvous Host IP"
	SAVE_RENDEZVOUS_IP_BUTTON = "<Save Rendezvous IP Address>"
)

func (u *UI) createRendezvousIPPage(config checks.Config) {
	primaryCheck := tview.NewTable()
	primaryCheck.SetBorder(true)
	primaryCheck.SetTitle("  Current Host IPs  ")
	primaryCheck.SetBorderColor(newt.ColorBlack)
	primaryCheck.SetBackgroundColor(newt.ColorGray)
	primaryCheck.SetTitleColor(newt.ColorBlack)

	u.rendezvousIPForm = tview.NewForm()
	u.rendezvousIPForm.SetBorder(false)
	u.rendezvousIPForm.SetBackgroundColor(newt.ColorGray)
	u.rendezvousIPForm.SetButtonsAlign(tview.AlignCenter)
	u.rendezvousIPForm.AddTextView(FIELD_CURRENT_HOST_IP, "192.168.x.x", 40, 1, false, false)
	u.rendezvousIPForm.AddInputField(FIELD_RENDEZVOUS_HOST_IP, "", 40, nil, nil)
	u.rendezvousIPForm.AddButton(SAVE_RENDEZVOUS_IP_BUTTON, func() {
		// save rendezvous IP address and switch to checks page
		ipAddress := u.rendezvousIPForm.GetFormItemByLabel(FIELD_RENDEZVOUS_HOST_IP).(*tview.InputField).GetText()
		validationError := validateIP(ipAddress)
		if validationError != "" {
			u.ShowInvalidIPAddressDialog(ipAddress)
		} else {
			// set focus to checks page and let controller know rendezvousIP is set
			u.setIsRendezousIPFormActive(false)
			u.returnFocusToChecks()
		}
	})
	u.rendezvousIPForm.SetButtonActivatedStyle(tcell.StyleDefault.Background(newt.ColorRed).
		Foreground(newt.ColorGray))
	u.rendezvousIPForm.SetButtonStyle(tcell.StyleDefault.Background(newt.ColorGray).
		Foreground(newt.ColorBlack))

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(u.rendezvousIPForm, 8, 0, false)
	mainFlex.SetTitle("  Rendezvous Host IP Setup  ").
		SetTitleColor(newt.ColorRed).
		SetBorder(true).
		SetBackgroundColor(newt.ColorGray).
		SetBorderColor(tcell.ColorBlack)

	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(mainFlex, mainFlexHeight, 0, false).
		AddItem(nil, 0, 1, false)

	// Allow the user to cycle the focus only over the configured items
	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyUp:
			u.focusedItem++
			if u.focusedItem > len(u.focusableItems)-1 {
				u.focusedItem = 0
			}

		case tcell.KeyBacktab, tcell.KeyDown:
			u.focusedItem--
			if u.focusedItem < 0 {
				u.focusedItem = len(u.focusableItems) - 1
			}

		default:
			// forward the event to the default handler
			return event
		}

		u.app.SetFocus(u.focusableItems[u.focusedItem])
		return nil
	})

	width := 80
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(innerFlex, width, 1, false).
		AddItem(nil, 0, 1, false)

	u.pages.SetBackgroundColor(newt.ColorBlue)
	u.pages.AddPage(PAGE_RENDEZVOUS_IP, flex, true, true)
}

func validateIP(ipAddress string) string {
	if net.ParseIP(ipAddress) == nil {
		return fmt.Sprintf("%s is not a valid IP address", ipAddress)
	}
	return ""
}

func saveIPAddress(ipAddress string) string {
	return ""
}
