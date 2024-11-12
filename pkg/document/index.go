package document

const (
	CSV   string = "csv"
	Excel string = "excel"
	JSON  string = "json"
)

func Use(documentType string) Document {
	switch documentType {
	case CSV:
		return DefaultCSVDocument
	case Excel:
		return DefaultExcelDocument
	default:
		return DefaultCSVDocument
	}
}
