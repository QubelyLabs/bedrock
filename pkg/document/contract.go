package document

import "io"

type Document interface {
	// Validate should append two new fields to the header: status, message
	// every validated record should have a status of true or false
	// while the message field show contain the error message if available or an empty string
	Validate(r io.Reader, headerValidator func(header []string) error, recordValidator func(header []string, record []string) error) (bool, string)
	Export(w io.Writer) error
	Import(r io.Reader) error
	ToSlice() [][]string
	ToMap() ([]map[string]any, error)
	FromSlice(data [][]string) error
	FromMap(data []map[string]any) error
}
