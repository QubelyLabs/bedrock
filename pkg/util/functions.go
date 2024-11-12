package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"net/url"

	"golang.org/x/crypto/bcrypt"
)

// Get attempts to retrieve a value from a KVStore and assert it to a specific type
func GetFromMap[T any](m Object, key string) T {
	value, ok := m[key].(T)
	if !ok {
		return *new(T)
	}

	return value
}

// SetToMap attempts to store a value in a KVStore
func SetToMap[T any](m Object, key string, value T) {
	m[key] = value
}

// RemoveFromMap removes a key-value pair from a KVStore
func RemoveFromMap(m Object, key string) {
	delete(m, key)
}

func ToString(value any) (string, bool) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.String:
		return v.String(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), true
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), true
	default:
		return "", false
	}
}

func ToFloat64(value any) (float64, bool) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	default:
		return 0, false
	}
}

func ToBool(value any) (bool, bool) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Bool:
		return v.Bool(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0, true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() != 0, true
	case reflect.String:
		s := v.String()
		b, err := strconv.ParseBool(s)
		return b, err == nil
	default:
		return false, false
	}
}

func ToTime(value any) (time.Time, bool) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.String:
		t, err := time.Parse(time.RFC3339, v.String())
		return t, err == nil
	default:
		return time.Time{}, false
	}
}

func ConsistentHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyHash(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func Random(len int) string {
	num := make([]string, len)
	for i := 0; i <= len-1; i++ {
		num[i] = strconv.Itoa(rand.Intn(9))
	}
	return strings.Join(num, "")
}

func ToBase64(data Object) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	base64String := base64.StdEncoding.EncodeToString(jsonData)

	return base64String
}

func FromBase64(str string) Object {
	jsonData, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Unmarshal the JSON bytes
	var data Object
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil
	}

	return data
}

func ToBase64String(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func FromBase64String(str string) (string, error) {
	value, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(value), nil
}

// func getStringProperty(properties map[string]interface{}, key string) (string, error) {
//     value, ok := properties[key]
//     if !ok || value == "" {
//         return "", &PropertyError{PropertyName: key}
//     }
//     strValue, ok := value.(string)
//     if !ok {
//         return "", &PropertyError{PropertyName: key}
//     }
//     return strValue, nil
// }

// func getMapProperty(properties map[string]interface{}, key string) (map[string]util.Object, error) {
//     value, ok := properties[key]
//     if !ok {
//         return nil, &PropertyError{PropertyName: key}
//     }
//     mapValue, ok := value.(map[string]util.Object)
//     if !ok {
//         return nil, & PropertyError{PropertyName: key}
//     }
//     return mapValue, nil
// }

func Min(items ...uint64) uint64 {
	if len(items) == 0 {
		return uint64(0)
	}

	min := items[0]
	for _, i := range items[1:] {
		if i < min {
			min = i
		}
	}

	return min
}

func FormatNumberToString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case int:
		if v == 0 {
			return ""
		}
		return strconv.Itoa(v)
	case int64:
		if v == 0 {
			return ""
		}
		return strconv.FormatInt(v, 10)
	case float64:
		if v == 0 {
			return ""
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		log.Printf("FormatCustomerNumber - unable to convert value (%v) of type %T to a string equivalent", v, v)
		return ""
	}
}

func UTCTime() time.Time {
	return time.Now().UTC()
}

func LocalTime(name string) time.Time {
	if name == "" {
		name = "Africa/Lagos"
	}

	// Load the desired location
	location, err := time.LoadLocation(name)
	if err != nil {
		fmt.Println("Error loading location", name, err)

		// try local location
		location, err = time.LoadLocation("Local")

		if err != nil {
			fmt.Println("Error loading local location", err)

			// we can live with UTC at this point
			return UTCTime()
		}
	}

	return time.Now().In(location)
}

func ParseUintWithDefault(value string, defaultValue uint64) uint64 {
	parsedValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

func MakeFilter(queries url.Values, model any) (map[string]any, uint64, uint64) {
	page := ParseUintWithDefault(queries.Get("page"), 1)
	perPage := ParseUintWithDefault(queries.Get("per_page"), 12)

	filterType := queries.Get("type")
	if filterType == "" || !InArray(FilterTypes, filterType) {
		filterType = FilterTypePartialMatch
	}

	from := queries.Get("from")
	to := queries.Get("to")
	timeKey := queries.Get("timeKey")
	if timeKey == "" {
		timeKey = "createdAt"
	}

	exactFilter := map[string]any{}
	partialFilter := []map[string]any{}
	for key, values := range queries {
		if len(values) < 1 {
			continue
		}

		value := values[0]

		// ignore for predefined queries
		if InArray([]string{"page", "per_page", "type", "from", "to", "timeKey"}, value) {
			continue
		}

		// ignore for queries not in the provided struct/model i.e mongo collection
		if !InArray(GetStructTags(model, "bson"), value) {
			continue
		}

		if filterType == FilterTypeExactMatch {
			exactFilter[key] = value
		} else {
			partialFilter = append(partialFilter, map[string]any{
				key: map[string]any{
					"$regex":   value,
					"$options": "i",
				},
			})
		}
	}

	filter := map[string]any{}
	if filterType == FilterTypeExactMatch {
		filter = exactFilter
	} else if len(filter) > 0 {
		filter = map[string]any{
			"$or": partialFilter,
		}
	}

	if from != "" && to != "" {
		f, fErr := time.Parse(DATE_TIME_FORMAT, from)
		t, tErr := time.Parse(DATE_TIME_FORMAT, to)
		if fErr == nil && tErr == nil {
			if len(filter) > 0 {
				filter = map[string]any{
					"$and": []map[string]any{
						filter,
						{
							timeKey: map[string]any{
								"$gte": f,
								"$lte": t,
							},
						},
					},
				}
			} else {
				filter = map[string]any{
					"$and": []map[string]any{
						{
							timeKey: map[string]any{
								"$gte": f,
								"$lte": t,
							},
						},
					},
				}
			}
		}
	} else if from != "" && to == "" {
		f, err := time.Parse(DATE_TIME_FORMAT, from)
		if err == nil {
			if len(filter) > 0 {
				filter = map[string]any{
					"$and": []map[string]any{
						filter,
						{
							timeKey: map[string]any{
								"$gte": f,
							},
						},
					},
				}
			} else {
				filter = map[string]any{
					"$and": []map[string]any{
						{
							timeKey: map[string]any{
								"$gte": f,
							},
						},
					},
				}
			}
		}
	} else if from == "" && to != "" {
		t, err := time.Parse(DATE_TIME_FORMAT, to)
		if err == nil {
			if len(filter) > 0 {
				filter = map[string]any{
					"$and": []map[string]any{
						filter,
						{
							timeKey: map[string]any{
								"$lte": t,
							},
						},
					},
				}
			} else {
				filter = map[string]any{
					"$and": []map[string]any{
						{
							timeKey: map[string]any{
								"$lte": t,
							},
						},
					},
				}
			}
		}
	}

	return filter, page, perPage
}

func ArraySum(arr []float64) float64 {
	res := 0.0
	for i := 0; i < len(arr); i++ {
		res += arr[i]
	}
	return res
}

func InArray[T comparable](arr []T, pin T) bool {

	for i := 0; i < len(arr); i++ {

		if arr[i] == pin {
			return true
		}
	}

	return false
}

func ArrayMap[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

// Get retrieves tags from an arbitrary struct and return then as an array
// group can be any of "bson", "json" etc
func GetStructTags(obj any, group string) []string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var tags []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(group)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

func FilterNonZeroFields(input any) map[string]any {
	result := map[string]any{}
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldTag := t.Field(i).Tag.Get("json")

		var keyName string
		if fieldTag != "" {
			parts := strings.Split(fieldTag, ",")
			keyName = parts[0]
		} else {
			keyName = t.Field(i).Name
		}

		if field.Kind() == reflect.Struct {
			nestedResult := FilterNonZeroFields(field.Interface())
			if len(nestedResult) > 0 {
				result[keyName] = nestedResult
			}
		} else {
			if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				result[keyName] = field.Interface()
			}
		}
	}
	return result
}

func AddDeleteFilter(filter map[string]any) map[string]any {
	if len(filter) > 0 {
		return map[string]any{
			"$and": []map[string]any{
				filter,
				{
					"deletedAt": nil,
				},
			},
		}
	}

	return map[string]any{
		"$and": []map[string]any{
			{
				"deletedAt": nil,
			},
		},
	}
}

func ParsePercentageString(input string) (bool, float64, error) {
	endsWithPercent := strings.HasSuffix(input, "%")
	if endsWithPercent {
		input = input[:len(input)-1]
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return endsWithPercent, 0, fmt.Errorf("invalid input: %w", err)
	}

	return endsWithPercent, value, nil
}

func IsPartOfList(list, item string) bool {
	items := strings.Split(list, ",")
	for _, v := range items {
		if strings.TrimSpace(v) == item {
			return true
		}
	}
	return false
}

// [TODO] refactor, optimize
func HydrateStructFromMap(m any, s any) error {
	// Convert map to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Convert JSON to struct
	err = json.Unmarshal(jsonData, &s)
	if err != nil {
		return err
	}

	return nil
}

// resolveValue resolves a dotted path like "ledger.updatedAt" in a map
func resolveValue(m map[string]any, path string) (any, bool) {
	keys := strings.Split(path, ".")
	var current any = m

	for _, key := range keys {
		if reflect.TypeOf(current).Kind() == reflect.Map {
			currentMap := current.(map[string]any)
			val, exists := currentMap[key]
			if !exists {
				return nil, false
			}
			current = val
		} else {
			return nil, false
		}
	}
	return current, true
}

// formatTime formats the time based on the provided Go format
func formatTime(input any, format string) (string, bool) {
	t, ok := input.(time.Time)
	if !ok {
		return "", false
	}
	return t.Format(format), true
}

// transformMap applies the template to the input map and returns a new map
func TransformMapForReport(records []map[string]any, template map[string]string) []map[string]any {
	results := []map[string]any{}

	for _, record := range records {
		result := map[string]any{}
		for newKey, templateValue := range template {
			parts := strings.Split(templateValue, ",")
			fieldPath := strings.TrimSpace(parts[0])

			value, found := resolveValue(record, fieldPath)
			if !found {
				result[newKey] = strings.TrimSpace(templateValue)
			} else {
				if len(parts) > 1 {
					format := strings.TrimSpace(parts[1])
					formattedValue, ok := formatTime(value, format)
					if ok {
						result[newKey] = formattedValue
					} else {
						result[newKey] = fmt.Sprintf("%v", value)
					}
				} else {
					result[newKey] = fmt.Sprintf("%v", value)
				}
			}
		}

		results = append(results, result)
	}

	return results
}
