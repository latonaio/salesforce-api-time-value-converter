package salesforce_api_time_value_converter

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ConvertToSalesforceTimeFormat(t time.Time) string {
	return fmt.Sprintf(`/Date(%d)/`, t.UnixMilli())
}
func ConvertToSalesforceTimeDurationFormat(t time.Time) string {
	return fmt.Sprintf(`PT%02dH%02dM%02dS`, t.UTC().Hour(), t.UTC().Minute(), t.UTC().Second())
}

func ChangeFormatToReadableDateTime(salesforceTime string) string {
	if salesforceTime == "" {
		return ""
	}
	t := ConvertToTimeFormat(salesforceTime)
	if t.Year() <= 1 {
		return t.Format("15:04:05")
	}

	if t.UTC().Hour() == 0 && t.UTC().Minute() == 0 && t.UTC().Second() == 0 && t.UTC().Nanosecond() == 0 {
		return t.Format("2006-01-02")
	}
	if t.Year() >= 10000 {
		return "9999-12-31"
	}

	return t.Format(time.RFC3339)
}

func ChangeFormatToReadableTime(salesforceTime string) string {
	if salesforceTime == "" {
		return ""
	}

	t, err := time.Parse("PT15H04M05S", salesforceTime)
	if err != nil {
		return salesforceTime
	}

	return t.Format("15:04:05")
}

func ChangeFormatToReadableTimeFromConsecutiveFormat(salesforceTime string) string {
	if salesforceTime == "" {
		return ""
	}

	t, err := time.Parse("20060102150405", salesforceTime)
	if err != nil {
		return salesforceTime
	}
	if t.Year() <= 1 {
		return t.Format("15:04:05")
	}
	if t.UTC().Hour() == 0 && t.UTC().Minute() == 0 && t.UTC().Second() == 0 && t.UTC().Nanosecond() == 0 {
		return t.Format("2006-01-02")
	}
	if t.Year() >= 10000 {
		return "9999-12-31T23:59:59+00:00"
	}

	return t.Format(time.RFC3339)
}

func ChangeFormatToSalesforceFormat(readableTime string) string {
	if readableTime == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, readableTime)
	if err == nil {
		return ConvertToSalesforceTimeFormat(t)
	}

	t, err = time.Parse("2006-01-02", readableTime)
	if err == nil {
		return ConvertToSalesforceTimeFormat(t)
	}

	t, err = time.Parse("15:04:05", readableTime)
	if err == nil {
		return ConvertToSalesforceTimeDurationFormat(t)
	}

	return readableTime
}

func ChangeTimeFormatToReadableForStruct(str interface{}) {
	rv := reflect.ValueOf(str)
	pickStringToReadable(rv)
}

func ChangeTimeFormatToSalesforceFormatStruct(str interface{}) {
	rv := reflect.ValueOf(str)
	pickStringToSalesforceFormat(rv)
}

func pickStringToSalesforceFormat(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		pickStringToSalesforceFormat(rv.Elem())
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			pickStringToSalesforceFormat(rv.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			pickStringToSalesforceFormat(rv.Field(i))
		}
	case reflect.Map:
		for _, e := range rv.MapKeys() {
			pickStringToSalesforceFormat(rv.MapIndex(e))
		}
	}

	if rv.Kind() == reflect.String {
		changeValueToSAPFormat(rv)
	}
}

func changeValueToSAPFormat(rv reflect.Value) {
	if rv.Kind() != reflect.String {
		return
	}
	if !rv.CanSet() {
		return
	}

	strValue := rv.String()
	if isReadableTimeFormat(strValue) {
		rv.SetString(ChangeFormatToSalesforceFormat(strValue))
	}
	if isReadableDateFormat(strValue) {
		rv.SetString(ChangeFormatToSalesforceFormat(strValue))
	}

}

func pickStringToReadable(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		pickStringToReadable(rv.Elem())
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			pickStringToReadable(rv.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			pickStringToReadable(rv.Field(i))
		}

	}
	if rv.Kind() == reflect.String {
		changeValueToReadable(rv)
	}
}

func changeValueToReadable(rv reflect.Value) {
	if rv.Kind() != reflect.String {
		return
	}
	if !rv.CanSet() {
		return
	}

	strValue := rv.String()
	if isSalesforceDateFormat(strValue) {
		rv.SetString(ChangeFormatToReadableDateTime(strValue))
		return
	}

	if isSalesforceDurationFormat(strValue) {
		rv.SetString(ChangeFormatToReadableTime(strValue))
		return
	}
	if isSalesforceDateTimeFormat(strValue) {
		rv.SetString(ChangeFormatToReadableTimeFromConsecutiveFormat(strValue))
		return
	}

}

func ConvertToTimeFormat(salesforceTime string) time.Time {
	err := validateSalesforceDateFormat(salesforceTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return time.Time{}
	}

	milli, err := getUnixmilli(salesforceTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return time.Time{}
	}
	return time.UnixMilli(milli)
}

// getUnixmilli unixミリ秒を返す
func getUnixmilli(salesforceTime string) (int64, error) {
	fixedString := strings.Join(strings.Split(salesforceTime, `\`), "")
	num := fixedString[len(`/Date(`) : len(fixedString)-len(`)/`)]
	num = strings.Split(num, "+")[0]
	milli, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("given word '%s' can not be converted to number: %w", salesforceTime, err)
	}
	return milli, nil
}
