package ui

import (
	"sync/atomic"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/checks"
	"github.com/rivo/tview"
)

type UI struct {
	app                 *tview.Application
	pages               *tview.Pages
	mainFlex, innerFlex *tview.Flex
	primaryCheck        *tview.Table
	checks              *tview.Table    // summary of all checks
	details             *tview.TextView // where errors from checks are displayed
	netConfigForm       *tview.Form     // contains "Configure network" button
	timeoutModal        *tview.Modal    // popup window that times out
	splashScreen        *tview.Modal    // display initial waiting message
	nmtuiActive         atomic.Value
	timeoutDialogActive atomic.Value
	timeoutDialogCancel chan bool
	dirty               atomic.Value // dirty flag set if the user interacts with the ui

	rendezvousIPForm       *tview.Form
	invalidIPAddressModal  *tview.Modal
	rendezvousIPFormActive atomic.Value

	focusableItems []tview.Primitive // the list of widgets that can be focused
	focusedItem    int               // the current focused widget
}

func NewUI(app *tview.Application, config checks.Config) *UI {
	ui := &UI{
		app:                 app,
		timeoutDialogCancel: make(chan bool),
	}
	ui.nmtuiActive.Store(false)
	ui.timeoutDialogActive.Store(false)
	ui.dirty.Store(false)
	ui.create(config)
	return ui
}

func (u *UI) GetApp() *tview.Application {
	return u.app
}

func (u *UI) GetPages() *tview.Pages {
	return u.pages
}

func (u *UI) returnFocusToChecks() {
	// reset u.focusableItems to those on the checks page
	u.focusableItems = []tview.Primitive{
		u.netConfigForm.GetButton(0),
		u.netConfigForm.GetButton(1),
	}
	u.pages.SwitchToPage(PAGE_CHECKSCREEN)
	// shifting focus back to the "Configure network"
	// button requires setting focus in this sequence
	// form -> form-button
	u.app.SetFocus(u.netConfigForm)
	u.app.SetFocus(u.netConfigForm.GetButton(0))
}

func (u *UI) setFocusToRendezvousIP() {
	u.setIsRendezousIPFormActive(true)
	// reset u.focusableItems to those on the rendezvous IP page
	u.focusableItems = []tview.Primitive{
		u.rendezvousIPForm.GetButton(0),
		u.rendezvousIPForm.GetFormItemByLabel(FIELD_RENDEZVOUS_HOST_IP),
	}

	u.pages.SwitchToPage(PAGE_RENDEZVOUS_IP)
	// shifting focus back to the "Configure network"
	// button requires setting focus in this sequence
	// form -> form-button
	u.app.SetFocus(u.rendezvousIPForm)
	u.app.SetFocus(u.rendezvousIPForm.GetFormItemByLabel(FIELD_RENDEZVOUS_HOST_IP))
}

func (u *UI) IsNMTuiActive() bool {
	return u.nmtuiActive.Load().(bool)
}

func (u *UI) setIsTimeoutDialogActive(isActive bool) {
	u.timeoutDialogActive.Store(isActive)
}

func (u *UI) IsTimeoutDialogActive() bool {
	return u.timeoutDialogActive.Load().(bool)
}

func (u *UI) setIsRendezousIPFormActive(isActive bool) {
	u.rendezvousIPFormActive.Store(isActive)
}

func (u *UI) IsRendezvousIPFormActive() bool {
	return u.rendezvousIPFormActive.Load().(bool)
}

func (u *UI) IsDirty() bool {
	return u.dirty.Load().(bool)
}

func (u *UI) create(config checks.Config) {
	u.pages = tview.NewPages()
	u.createCheckPage(config)
	u.createTimeoutModal(config)
	u.createSplashScreen()
	u.createRendezvousIPPage(config)
	u.createInvalidIPAddressModal()
	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !u.IsRendezvousIPFormActive() {
			u.dirty.Store(true)
		}
		return event
	})
}
