package accounting

import (
	"bytes"
	"github.com/jinzhu/now"
	"github.com/jung-kurt/gofpdf"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"
	"path/filepath"
	"psutils/pkg/config"
	"text/template"
	"time"
)

var t *template.Template

type ConfirmationConfig struct {
	Date      time.Time
	Month     string
	Seller    config.Company
	Customer  config.Company
	StartDate string
	EndDate   string
	Items     []string
}

var payload = ConfirmationConfig{
	Date: time.Time{},
}

func init() {
	t = template.Must(template.ParseFiles(filepath.Join("pkg", "accounting", "templates", "confirmation.tmpl")))
}

func CreateConfirmation(c *cli.Context) {
	conf := config.Load()

	date, err := time.Parse(conf.Other.MontDateFormat, c.String("date"))
	if err != nil {
		log.With("err", err).Fatalf("can't parse date parameter")
	}

	startDate := now.New(date).BeginningOfMonth().Format("2006-01-02")
	endDate := now.New(date).EndOfMonth().Format("2006-01-02")
	payload.Date = date
	payload.StartDate = startDate
	payload.EndDate = endDate
	payload.Month = date.Format("01/2006")
	payload.Customer = conf.Accounting.Customer
	payload.Seller = conf.Accounting.Seller
	payload.Items = conf.Accounting.Items

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetCompression(true)
	pdf.SetFont("Arial", "", 20)
	//_, lineHt := pdf.GetFontSize()

	tmpl := pdf.CreateTemplate(func(tpl *gofpdf.Tpl) {
		tpl.SetFontSize(14)
		_, lineHt := pdf.GetFontSize()

		reader := bytes.NewBufferString("")
		err := t.Execute(reader, payload)
		if err != nil {
			log.With("err", err).Fatalf("can't parse tmpl")
			return
		}

		str := reader.String()
		html := tpl.HTMLBasicNew()
		html.Write(lineHt, str)
		return
	})

	pdf.AddPage()
	pdf.UseTemplate(tmpl)

	err = pdf.OutputFileAndClose("./test.pdf")
	if err != nil {
		log.With("err", err).Fatal("creating output file")
	}
}