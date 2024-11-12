package document

import (
	"encoding/csv"
	"errors"
	"io"
	"log"

	"github.com/QubelyLabs/bedrock/pkg/util"
)

var (
	DefaultCSVDocument = NewCSVDocument()
)

type csvDocument struct {
	header  []string
	records [][]string
	len     int64
}

func (d *csvDocument) Validate(r io.Reader, headerValidator func(header []string) error, recordValidator func(header []string, record []string) error) (bool, string) {
	d.reset()

	csvReader := csv.NewReader(r)

	d.len = 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Unable to read file")
			return false, "unable to read file"
		}

		if d.len == 0 {
			err := headerValidator(record)
			record = append(record, "errorStatus", "errorMessage")
			if err != nil {
				return false, err.Error()
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

func (d *csvDocument) Export(w io.Writer) error {
	csvWriter := csv.NewWriter(w)

	if len(d.header) > 0 {
		err := csvWriter.Write(d.header)

		if err != nil {
			return err
		}
	}

	if len(d.records) > 0 {
		for _, record := range d.records {
			err := csvWriter.Write(record)

			// reset header and record if any record line
			// could not be written
			// this way, we won't end up with a partial document
			if err != nil {
				d.reset()

				return err
			}
		}
	}

	csvWriter.Flush()

	return nil
}

func (d *csvDocument) Import(r io.Reader) error {
	d.reset()

	csvReader := csv.NewReader(r)

	d.len = 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if d.len == 0 {
			d.header = record
		} else {
			d.records = append(d.records, record)
		}

		d.len++
	}

	return nil
}

func (d *csvDocument) ToSlice() [][]string {
	result := make([][]string, len(d.records)+1)
	result[0] = d.header
	copy(result[1:], d.records)
	return result
}

func (d *csvDocument) ToMap() ([]map[string]any, error) {
	results := []map[string]any{}

	for _, record := range d.records {
		result := map[string]any{}

		for i, field := range record {
			title := ""

			if i >= len(d.header) {
				return nil, errors.New("header vs records title mismatch")
			} else {
				title = d.header[i]
			}

			result[title] = field
		}

		results = append(results, result)
	}

	return results, nil
}

func (d *csvDocument) FromSlice(data [][]string) error {
	if len(data) > 0 {
		d.header = data[0]
	}

	if len(data) > 1 {
		d.records = data[1:]
	}

	return nil
}

func (d *csvDocument) FromMap(data []map[string]any) error {
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

func (d *csvDocument) reset() {
	d.header = []string{}
	d.records = [][]string{}
}

func NewCSVDocument() *csvDocument {
	header := []string{}
	records := [][]string{}
	return &csvDocument{
		header:  header,
		records: records,
	}
}
