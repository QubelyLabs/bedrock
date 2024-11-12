package util

type Object = map[string]any

const (
	FilterTypeExactMatch   = "exact"
	FilterTypePartialMatch = "partial"
)

var (
	FilterTypes = []string{FilterTypeExactMatch, FilterTypePartialMatch}
)

const (
	DATE_TIME_FORMAT  = "2006-01-02T15:04:05"
	DATE_FORMAT       = "2006-01-02"
	TIME_FORMAT       = "15:04:05"
	SHORT_TIME_FORMAT = "15:04"
	WEEK_DAY_FORMAT   = "Mon"
)
