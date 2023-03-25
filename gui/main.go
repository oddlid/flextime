package main

/*
This package is mainly just to play around with the fyne library.
I'm not even sure if a GUI is very useful for this. But ok as a PoC, I guess.

I'm not very used to GUI programming, so I feel this is getting quite messy very fast.
But at least I'll maybe learn something in the process.
*/

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/oddlid/flextime/flex"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _db *flex.DB
var _currentCustomer *flex.Customer

func init() {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999-07:00"

	dbfile, present := os.LookupEnv("FLEXTIME_FILE")
	if !present {
		log.Debug().Msg("$FLEXTIME_FILE not set, not attempting load")
		return
	}
	log.Debug().Str("FLEXTIME_FILE", dbfile).Msg("DB file was specified")
	db, err := openDB(dbfile)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	log.Debug().Str("FLEXTIME_FILE", dbfile).Msg("DB file was loaded successfully")
	_db = db
}

type addEntryRowWidget struct {
	txtDateBinding    binding.String
	txtAmountBinding  binding.String
	txtCommentBinding binding.String
	btnAdd            *widget.Button
	layoutContainer   *fyne.Container
}

type customerListWidget struct {
	customerNames        []string
	customerNamesBinding binding.ExternalStringList
	list                 *widget.List
	// don't know yet if we need a layout for this struct
}

type customerHeaderWidget struct {
	nameBinding      binding.String
	flexTotalBinding binding.String
	layoutContainer  *fyne.Container
}

func getAddEntryRowWidget(btnAddFunc func()) *addEntryRowWidget {
	w := &addEntryRowWidget{
		btnAdd:            widget.NewButton("Add", btnAddFunc),
		txtDateBinding:    binding.NewString(),
		txtAmountBinding:  binding.NewString(),
		txtCommentBinding: binding.NewString(),
	}
	lblDate := widget.NewLabel("Date:")
	lblAmount := widget.NewLabel("Amount:")
	lblComment := widget.NewLabel("Comment")
	txtDate := widget.NewEntryWithData(w.txtDateBinding)
	txtAmount := widget.NewEntryWithData(w.txtAmountBinding)
	txtComment := widget.NewEntryWithData(w.txtCommentBinding)
	layoutContainer := container.NewBorder(
		nil, nil, nil, w.btnAdd,
		w.btnAdd,
		container.NewGridWithColumns(
			3,
			container.NewBorder(
				nil, nil, lblDate, nil,
				lblDate,
				txtDate,
			),
			container.NewBorder(
				nil, nil, lblAmount, nil,
				lblAmount,
				txtAmount,
			),
			container.NewBorder(
				nil, nil, lblComment, nil,
				lblComment,
				txtComment,
			),
		),
	)
	w.layoutContainer = layoutContainer

	return w
}

func longestEntry(slice []string) int {
	maxLen := 0
	for _, entry := range slice {
		entryLen := len(entry)
		if entryLen > maxLen {
			maxLen = entryLen
		}
	}
	return maxLen
}

func getListTemplateString(names []string) string {
	return strings.Repeat("#", longestEntry(names)+2)
}

func getCustomerListWidget() *customerListWidget {
	w := &customerListWidget{
		customerNames: []string{"Customer1", "Customer2", "Customer3"}, // filler values to have before populating for real
	}
	w.customerNamesBinding = binding.BindStringList(&w.customerNames)

	list := widget.NewListWithData(
		w.customerNamesBinding,
		func() fyne.CanvasObject {
			return widget.NewLabel(getListTemplateString(w.customerNames))
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			co.(*widget.Label).Bind(di.(binding.String))
		},
	)
	w.list = list

	return w
}

func (clw *customerListWidget) sync(customers flex.Customers) {
	clw.customerNames = clw.customerNames[:0]
	for _, c := range customers {
		clw.customerNames = append(clw.customerNames, c.Name)
	}
	clw.customerNamesBinding.Reload()
}

func getCustomerHeaderWidget() *customerHeaderWidget {
	w := &customerHeaderWidget{
		nameBinding:      binding.NewString(),
		flexTotalBinding: binding.NewString(),
	}
	lblName := widget.NewLabelWithData(w.nameBinding)
	lblFlexTotal := widget.NewLabelWithData(w.flexTotalBinding)

	layoutContainer := container.NewBorder(
		nil, nil, lblName, lblFlexTotal,
		lblName,
		lblFlexTotal,
	)
	w.layoutContainer = layoutContainer

	return w
}

func (chw *customerHeaderWidget) sync(c *flex.Customer) {
	chw.nameBinding.Set(c.Name)
	chw.flexTotalBinding.Set(c.GetTotalFlex().String())
}

func getWindowContainer(
	customerList *customerListWidget,
	customerHeader *customerHeaderWidget,
	addEntryWidget *addEntryRowWidget,
) *fyne.Container {
	rightPane := container.NewBorder(
		customerHeader.layoutContainer, addEntryWidget.layoutContainer, nil, nil,
		customerHeader.layoutContainer, addEntryWidget.layoutContainer,
	)
	return container.NewBorder(
		nil, nil, customerList.list, rightPane,
		customerList.list, rightPane,
	)
}

func main() {
	a := app.New()
	w := a.NewWindow("FlexTime GUI test")

	btnAddFunc := func() {
		if _currentCustomer != nil {
			log.Debug().Str("customer", _currentCustomer.Name).Msg("Add to this customer")
			return
		}
		log.Debug().Msg("Add button clicked")
	}
	clw := getCustomerListWidget()
	chw := getCustomerHeaderWidget()
	aerw := getAddEntryRowWidget(btnAddFunc)
	if _db != nil {
		clw.sync(_db.Customers)
		clw.list.OnSelected = func(id widget.ListItemID) {
			customerName := clw.customerNames[int(id)]
			customer, err := _db.GetCustomer(customerName)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}
			chw.sync(customer)
			_currentCustomer = customer
		}
	}
	w.SetContent(
		getWindowContainer(clw, chw, aerw),
	)

	w.ShowAndRun()
}
