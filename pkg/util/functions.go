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

	// Unmarshal the JSON bytes to a map[string]any
	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil
	}

	return data
}
