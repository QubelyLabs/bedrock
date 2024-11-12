package document

import (
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/QubelyLabs/bedrock/pkg/util"
	excel "github.com/xuri/excelize/v2"
)

var (
	DefaultExcelDocument = NewExcelDocument()
)

type excelDocument struct {
	header  []string
	records [][]string
	len     int64
}

func (d *excelDocument) Validate(r io.Reader, headerValidator func(header []string) error, recordValidator func(header []string, record []string) error) (bool, string) {
	d.reset()

	excelFile, err := excel.OpenReader(r)
	if err != nil {
		log.Printf("Unable to read file", err)
		return false, "unable to read file"
	}
	defer excelFile.Close()

	sheetName := excelFile.GetSheetName(1)
	if sheetName == "" {
		log.Printf("Unable to read file")
		return false, "unable to read file"
	}

	d.len = 0
	records, err := excelFile.GetRows(sheetName)
	if err != nil {
		log.Printf("Unable to read file", err)
		return false, "unable to read file"
	}

	for _, record := range records {
		if d.len == 0 {
			err := headerValidator(record)
			record = append(record, "errorStatus", "errorMessage")
			if err != nil {
				return false, "unable to read file"
			}

			d.header = record
		} else {
			err := recordValidator(d.header, record)
			if err != nil {
				record = append(record, "false", err.Error())
			} else {
				record = append(record, "true", "")
			}

			d.records = append(d.records, record)
		}

		d.len++
	}

	return true, ""

}

func (d *excelDocument) Export(w io.Writer) error {
	sheetName := "Sheet1"

	excelFile := excel.NewFile()
	defer func() {
		if err := excelFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	index, err := excelFile.NewSheet(sheetName)
	if err != nil {
		return err
	}

	if len(d.header) > 0 {
		err := excelFile.SetSheetRow(sheetName, "A"+"1", &d.header)

		if err != nil {
			return err
		}
	}

	if len(d.records) > 0 {
		for i, record := range d.records {
			err := excelFile.SetSheetRow(sheetName, "A"+strconv.Itoa(i+2), &record)

			// reset header and record if any record line
			// could not be written
			// this way, we won't end up with a partial document
			if err != nil {
				d.reset()

				return err
			}
		}
	}

	excelFile.SetActiveSheet(index)

	err = excelFile.Write(w)
	if err != nil {
		return err
	}

	return nil
}

func (d *excelDocument) Import(r io.Reader) error {
	d.reset()

	excelFile, err := excel.OpenReader(r)
	if err != nil {
		return err
	}
	defer excelFile.Close()

	sheetName := excelFile.GetSheetName(1)
	if sheetName == "" {
		return errors.New("invalid sheet, the workbook seems to have no sheet")
	}

	d.len = 0
	records, err := excelFile.GetRows(sheetName)
	if err != nil {
		return err
	}
	for _, record := range records {
		if d.len == 0 {
			d.header = record
		} else {
			d.records = append(d.records, record)
		}

		d.len++
	}

	return nil
}

func (d *excelDocument) ToSlice() [][]string {
	result := make([][]string, len(d.records)+1)
	result[0] = d.header
	copy(result[1:], d.records)
	return result
}

func (d *excelDocument) ToMap() ([]map[string]any, error) {
	results := []map[string]any{}

	for _, record := range d.records {
		result := map[string]any{}

		for i, field := range record {
			title := ""

			if len(d.header) > i {
				title = d.header[i+1]
			} else {
				return nil, errors.New("header vs records title mismatch")
			}

			result[title] = field
		}

		results = append(results, result)
	}

	return results, nil
}

func (d *excelDocument) FromSlice(data [][]string) error {
	if len(data) > 0 {
		d.header = data[0]
	}

	if len(data) > 1 {
		d.records = data[1:]
	}

	return nil
}

func (d *excelDocument) FromMap(data []map[string]any) error {
	for _, datum := range data {
		for key := range datum {
			if !util.InArray(d.header, key) {
				d.header = append(d.header, key)
			}
		}
	}

	if len(data) > 1 {
		for _, datum := range data[1:] {
			record := []string{}
			for i, key := range d.header {
				record[i] = datum[key].(string)
			}
			d.records = append(d.records, record)
		}
	}

	return nil
}

func (d *excelDocument) reset() {
	d.header = []string{}
	d.records = [][]string{}
}

func NewExcelDocument() *excelDocument {
	header := []string{}
	records := [][]string{}
	return &excelDocument{
		header:  header,
		records: records,
	}
}
