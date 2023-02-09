package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/checks"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/dialogs"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/newt"
	"github.com/rivo/tview"
)

const (
	CONFIGURE_NETWORK_LABEL string = "Configure Networking"
	RENDEZVOUS_IP_LABEL     string = "Rendezvous IP Address"
	RELEASE_IMAGE_LABEL     string = "Release Image"
)

func (u *UI) markCheckSuccess(row int, col int) {
	u.checks.SetCell(row, col, &tview.TableCell{Text: " ✓", Color: tcell.ColorGreen})
}

func (u *UI) markCheckFail(row int, col int) {
	u.checks.SetCell(row, col, &tview.TableCell{Text: " ✖", Color: tcell.ColorRed})
}

func (u *UI) appendNewErrorToDetails(heading string, errorString string) {
	u.appendToDetails(fmt.Sprintf("%s%s:%s\n%s", "[red]", heading, "[white]", errorString))
}

func (u *UI) appendToDetails(newLines string) {
	current := u.details.GetText(false)
	if len(current) > 10000 {
		// if details run more than 10000 characters, reset
		current = ""
	}
	u.details.SetText(current + newLines)
}

func (u *UI) createCheckPage(config checks.Config) {
	u.envVars = tview.NewTable()
	// TODO: remove if rendezvous IP does not need to be checked
	// u.config.SetCell(0, 0, &tview.TableCell{Text: " Rendezvous IP", Color: tcell.ColorWhite})
	// u.config.SetCell(0, 1, &tview.TableCell{Text: config.RendezvousHostIP, Color: newt.ColorBlue})
	u.envVars.SetCell(1, 0, &tview.TableCell{Text: " Release image URL", Color: tcell.ColorWhite})
	u.envVars.SetCell(1, 1, &tview.TableCell{Text: config.ReleaseImageURL, Color: newt.ColorBlue})

	releaseImageHostName, err := checks.ParseHostnameFromURL(config.ReleaseImageURL)
	if err != nil {
		dialogs.PanicDialog(u.app, err)
	}
	u.checks = tview.NewTable()
	u.checks.SetBorder(true)
	u.checks.SetTitle("Checks")
	u.checks.SetCell(0, 0, &tview.TableCell{Text: " ?", Color: tcell.ColorWhite})
	u.checks.SetCell(0, 1, &tview.TableCell{Text: "podman pull release image", Color: tcell.ColorWhite})
	u.checks.SetCell(1, 0, &tview.TableCell{Text: " ?", Color: tcell.ColorWhite})
	u.checks.SetCell(1, 1, &tview.TableCell{Text: "nslookup " + releaseImageHostName, Color: tcell.ColorWhite})
	u.checks.SetCell(2, 0, &tview.TableCell{Text: " ?", Color: tcell.ColorWhite})
	if releaseImageHostName == "quay.io" {
		u.checks.SetCell(2, 1, &tview.TableCell{Text: "quay.io does not respond to ping, ping skipped", Color: tcell.ColorWhite})
	} else {
		u.checks.SetCell(2, 1, &tview.TableCell{Text: "ping " + releaseImageHostName, Color: tcell.ColorWhite})
	}

	u.details = tview.NewTextView()
	u.details.SetBorder(true)
	u.details.SetTitle("Check Errors")
	u.details.SetTitleColor(tcell.ColorWhite)
	u.details.SetDynamicColors(true)

	u.form = tview.NewForm()
	u.form.SetBorder(false)
	u.form.AddButton("Configure network", u.NMTUIRunner(u.app, u.pages, nil))
	u.form.SetButtonsAlign(tview.AlignCenter)
	// u.buttons.SetBackgroundColor(tcell.ColorBlack)

	u.grid = tview.NewGrid().SetRows(3, 5, 0, 3).SetColumns(0).
		AddItem(u.envVars, 0, 0, 1, 1, 0, 0, false).
		AddItem(u.checks, 1, 0, 1, 1, 0, 0, false).
		AddItem(u.details, 2, 0, 1, 1, 0, 0, false).
		AddItem(u.form, 3, 0, 1, 1, 0, 0, false)
	u.grid.SetTitle("Agent installer network boot setup")
	u.grid.SetTitleColor(tcell.ColorWhite)
	u.grid.SetBorder(true)

	u.pages.AddPage("checkScreen", u.grid, true, true)

	u.app.SetRoot(u.pages, true).SetFocus(u.form).EnableMouse(true)
}