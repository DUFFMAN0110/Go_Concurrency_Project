package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Personal Finance Tracker")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Username")

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Amount")

	categoryEntry := widget.NewEntry()
	categoryEntry.SetPlaceHolder("Category")

	reportOutput := widget.NewMultiLineEntry()
	reportOutput.SetMinRowsVisible(10)
	reportOutput.Disable()

	addBtn := widget.NewButton("Add Expense", func() {
		user := userEntry.Text
		amount := amountEntry.Text
		category := categoryEntry.Text
		if user == "" || amount == "" || category == "" {
			dialog.ShowError(fmt.Errorf("All fields required"), w)
			return
		}
		addFromInput(user, amount, category, reportOutput)
	})

	reportBtn := widget.NewButton("Show Report", func() {
		user := userEntry.Text
		if user == "" {
			dialog.ShowError(fmt.Errorf("Username required"), w)
			return
		}
		reportOutput.SetText("")
		sendReport(user, reportOutput)
	})

	exportBtn := widget.NewButton("Export CSV", func() {
		user := userEntry.Text
		if user == "" {
			dialog.ShowError(fmt.Errorf("Username required"), w)
			return
		}
		exportCSV(user, reportOutput)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Personal Finance CLI (GUI Edition)"),
		userEntry,
		amountEntry,
		categoryEntry,
		container.NewHBox(addBtn, reportBtn, exportBtn),
		widget.NewLabel("Report:"),
		reportOutput,
	))

	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}
