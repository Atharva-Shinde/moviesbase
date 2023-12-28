package validator

import "regexp"

// page 96

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// checks if all values in the string slice are unique
// used for genres
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(uniqueValues) == len(values)
}

// checks if the value has any match for the regexp
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// checks if the value is present in the varidiac slice of string
func In(value string, list ...string) bool {
	for i := range list {
		if list[i] == value {
			return true
		}
	}
	return false
}

// checks if there are no errors in validations above
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// checks if validator function returned false and adds error message to the validator
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// adds errors to Errors map
func (v *Validator) AddError(key, message string) {
	if _, value := v.Errors[key]; !value {
		v.Errors[key] = message
	}
}
