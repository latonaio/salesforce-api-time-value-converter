package salesforce_api_time_value_converter

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func isSalesforceDateFormat(s string) bool {
	if err := validateSalesforceDateFormat(s); err != nil {
		return false
	}
	return true
}
func isSalesforceDurationFormat(salesforceTime string) bool {
	ok, _ := regexp.MatchString(`PT[0-2]\dH[0-6]\dM[0-6]\dS`, salesforceTime)
	return ok
}
func isSalesforceDateTimeFormat(salesforceTime string) bool {
	_, err := time.Parse("20060102150405", salesforceTime)
	return err == nil
}

func isReadableTimeFormat(s string) bool {
	if _, err := time.Parse(time.RFC3339, s); err != nil {
		return false
	}
	return true
}

func isReadableDateFormat(s string) bool {
	if _, err := time.Parse("2006-01-02", s); err != nil {
		return false
	}
	return true
}

func validateSalesforceDateFormat(salesforceTime string) error {
	err := validatePrefix(salesforceTime)
	if err != nil {
		return err
	}
	err = validateSuffix(salesforceTime)
	if err != nil {
		return err
	}
	return nil
}

func validatePrefix(salesforceTime string) error {
	if !(strings.HasPrefix(salesforceTime, `\/Date(`) || strings.HasPrefix(salesforceTime, `/Date(`)) {
		return fmt.Errorf(
			"%s is not type of SalesforceTime timestamp", salesforceTime,
		)
	}
	return nil
}

func validateSuffix(salesforceTime string) error {
	if !(strings.HasSuffix(salesforceTime, `)\/`) || strings.HasSuffix(salesforceTime, `)/`)) {
		return fmt.Errorf(
			"%s is not type of SalesforceTime timestamp", salesforceTime,
		)
	}
	return nil
}
