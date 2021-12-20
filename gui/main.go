package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type addEntryRowWidget struct {
	lblDate         *widget.Label
	lblAmount       *widget.Label
	lblComment      *widget.Label
	txtDate         *widget.Entry
	txtAmount       *widget.Entry
	txtComment      *widget.Entry
	btnAdd          *widget.Button
	layoutContainer *fyne.Container
}

func getAddEntryRowWidget(btnAddFunc func()) *addEntryRowWidget {
	w := &addEntryRowWidget{
		lblDate:    widget.NewLabel("Date:"),
		lblAmount:  widget.NewLabel("Amount:"),
		lblComment: widget.NewLabel("Comment"),
		txtDate:    widget.NewEntry(),
		txtAmount:  widget.NewEntry(),
		txtComment: widget.NewEntry(),
		btnAdd:     widget.NewButton("Add", btnAddFunc),
	}
	layoutContainer := container.New(
		layout.NewBorderLayout(nil, nil, nil, w.btnAdd),
		w.btnAdd,
		container.New(
			layout.NewGridLayout(3),
			container.New(
				layout.NewBorderLayout(nil, nil, w.lblDate, nil),
				w.lblDate,
				w.txtDate,
			),
			container.New(
				layout.NewBorderLayout(nil, nil, w.lblAmount, nil),
				w.lblAmount,
				w.txtAmount,
			),
			container.New(
				layout.NewBorderLayout(nil, nil, w.lblComment, nil),
				w.lblComment,
				w.txtComment,
			),
		),
	)
	w.layoutContainer = layoutContainer

	return w
}

func main() {
	a := app.New()
	w := a.NewWindow("FlexTime GUI test")

	w.SetContent(
		getAddEntryRowWidget(func() { log.Print("btnAdd clicked") }).layoutContainer,
	)

	w.ShowAndRun()
}
