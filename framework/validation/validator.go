package validation

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Errors maps field names to error messages.
type Errors map[string][]string

// Validator validates data against rules.
type Validator struct {
	ctx      context.Context
	data     interface{}
	rules    map[string]string
	data_map map[string]interface{}
	errors   Errors
}

// New creates a new validator instance.
func New(ctx context.Context, data interface{}, rules map[string]string) *Validator {
	return &Validator{
		ctx:      ctx,
		data:     data,
		rules:    rules,
		data_map: make(map[string]interface{}),
		errors:   make(Errors),
	}
}

// Validate performs validation and returns errors.
func (v *Validator) Validate() (Errors, error) {
	// Convert data to map (simplified approach)
	switch d := v.data.(type) {
	case map[string]interface{}:
		v.data_map = d
	default:
		// For struct types, would need reflection here
		// Simplified: treat as empty map
	}

	// Validate each field
	for field, rulesStr := range v.rules {
		if rulesStr == "" {
			continue
		}

		rules := strings.Split(rulesStr, "|")
		value := v.getFieldValue(field)

		v.validateField(field, value, rules)
	}

	if len(v.errors) == 0 {
		return nil, nil
	}

	return v.errors, nil
}

// validateField validates a single field against multiple rules.
func (v *Validator) validateField(field string, value interface{}, rules []string) {
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		// Parse rule and parameters
		parts := strings.SplitN(rule, ":", 2)
		ruleName := parts[0]
		var param string
		if len(parts) > 1 {
			param = parts[1]
		}

		if !v.applyRule(field, value, ruleName, param) {
			v.addError(field, v.getRuleMessage(ruleName, param))
		}
	}
}

// applyRule applies a validation rule to a value.
func (v *Validator) applyRule(field string, value interface{}, rule, param string) bool {
	switch rule {
	case "required":
		return v.isNotEmpty(value)
	case "nullable":
		return true // Always passes, means field can be empty
	case "string":
		_, ok := value.(string)
		return ok || value == nil
	case "integer", "int":
		return v.isInteger(value)
	case "numeric":
		return v.isNumeric(value)
	case "email":
		return v.isEmail(v.toString(value))
	case "url":
		return v.isURL(v.toString(value))
	case "regex":
		return v.matchRegex(v.toString(value), param)
	case "min":
		return v.minValue(value, param)
	case "max":
		return v.maxValue(value, param)
	case "between":
		return v.betweenValue(value, param)
	case "confirmed":
		return v.isConfirmed(field, value)
	case "unique":
		// Placeholder: would need database access
		return true
	case "exists":
		// Placeholder: would need database access
		return true
	case "starts_with":
		return strings.HasPrefix(v.toString(value), param)
	case "ends_with":
		return strings.HasSuffix(v.toString(value), param)
	case "array":
		_, ok := value.([]interface{})
		return ok || value == nil
	case "present":
		_, ok := v.data_map[field]
		return ok
	default:
		return true
	}
}

// getFieldValue retrieves the value of a field.
func (v *Validator) getFieldValue(field string) interface{} {
	val, ok := v.data_map[field]
	if !ok {
		return nil
	}
	return val
}

// isNotEmpty checks if a value is not empty.
func (v *Validator) isNotEmpty(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case []interface{}:
		return len(v) > 0
	default:
		return true
	}
}

// isInteger checks if a value is an integer.
func (v *Validator) isInteger(value interface{}) bool {
	switch value := value.(type) {
	case int, int32, int64:
		return true
	case string:
		_, err := strconv.ParseInt(value, 10, 64)
		return err == nil
	default:
		return false
	}
}

// isNumeric checks if a value is numeric.
func (v *Validator) isNumeric(value interface{}) bool {
	switch value := value.(type) {
	case int, int32, int64, float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	default:
		return false
	}
}

// isEmail checks if a value is a valid email.
func (v *Validator) isEmail(value string) bool {
	if value == "" {
		return false
	}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// isURL checks if a value is a valid URL.
func (v *Validator) isURL(value string) bool {
	if value == "" {
		return false
	}
	pattern := `^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// matchRegex checks if a value matches a regex pattern.
func (v *Validator) matchRegex(value, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// minValue checks if a value meets the minimum requirement.
func (v *Validator) minValue(value interface{}, param string) bool {
	min, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return false
	}

	switch v := value.(type) {
	case string:
		return int64(len(v)) >= min
	case int, int32, int64:
		return toInt64(v) >= min
	default:
		return false
	}
}

// maxValue checks if a value meets the maximum requirement.
func (v *Validator) maxValue(value interface{}, param string) bool {
	max, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return false
	}

	switch v := value.(type) {
	case string:
		return int64(len(v)) <= max
	case int, int32, int64:
		return toInt64(v) <= max
	default:
		return false
	}
}

// betweenValue checks if a value is between min and max.
func (v *Validator) betweenValue(value interface{}, param string) bool {
	parts := strings.Split(param, ",")
	if len(parts) != 2 {
		return false
	}

	min, errMin := strconv.ParseInt(parts[0], 10, 64)
	max, errMax := strconv.ParseInt(parts[1], 10, 64)

	if errMin != nil || errMax != nil {
		return false
	}

	switch v := value.(type) {
	case string:
		length := int64(len(v))
		return length >= min && length <= max
	case int, int32, int64:
		val := toInt64(v)
		return val >= min && val <= max
	default:
		return false
	}
}

// isConfirmed checks if a field matches its confirmation field.
func (v *Validator) isConfirmed(field string, value interface{}) bool {
	confirmField := field + "_confirmation"
	confirmValue, ok := v.data_map[confirmField]
	if !ok {
		return false
	}

	return v.toString(value) == v.toString(confirmValue)
}

// toString converts a value to a string.
func (v *Validator) toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// addError adds an error message for a field.
func (v *Validator) addError(field, message string) {
	if v.errors[field] == nil {
		v.errors[field] = make([]string, 0)
	}
	v.errors[field] = append(v.errors[field], message)
}

// getRuleMessage returns the error message for a rule.
func (v *Validator) getRuleMessage(rule, param string) string {
	messages := map[string]string{
		"required":    "This field is required",
		"email":       "This field must be a valid email",
		"string":      "This field must be a string",
		"integer":     "This field must be an integer",
		"numeric":     "This field must be numeric",
		"min":         fmt.Sprintf("This field must be at least %s", param),
		"max":         fmt.Sprintf("This field must not exceed %s", param),
		"between":     fmt.Sprintf("This field must be between %s", param),
		"confirmed":   "This field confirmation does not match",
		"starts_with": fmt.Sprintf("This field must start with %s", param),
		"ends_with":   fmt.Sprintf("This field must end with %s", param),
		"url":         "This field must be a valid URL",
		"regex":       "This field format is invalid",
	}

	if msg, ok := messages[rule]; ok {
		return msg
	}

	return "Validation failed"
}

// toInt64 converts a value to int64.
func toInt64(value interface{}) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}
