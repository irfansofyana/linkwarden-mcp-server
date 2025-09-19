package linkwardenmcp

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/mcpgo"
)

// Validator provides a fluent interface for validating parameters
// and collecting errors
type Validator struct {
	request *mcpgo.CallToolRequest
	errors  []error
}

// NewValidator creates a new validator for the given request
func NewValidator(r *mcpgo.CallToolRequest) *Validator {
	return &Validator{
		request: r,
		errors:  []error{},
	}
}

// addError adds a non-nil error to the collection
func (v *Validator) addError(err error) *Validator {
	if err != nil {
		v.errors = append(v.errors, err)
	}
	return v
}

// HasErrors returns true if there are any validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// HandleErrorsIfAny formats all errors and returns an appropriate tool result
func (v *Validator) HandleErrorsIfAny() (*mcpgo.ToolResult, error) {
	if v.HasErrors() {
		messages := make([]string, 0, len(v.errors))
		for _, err := range v.errors {
			messages = append(messages, err.Error())
		}
		errorMsg := "Validation errors:\n- " + strings.Join(messages, "\n- ")
		return mcpgo.NewToolResultError(errorMsg), nil
	}
	return nil, nil
}

// extractValueGeneric is a standalone generic function to extract a parameter
// of type T
func extractValueGeneric[T any](
	request *mcpgo.CallToolRequest,
	name string,
	required bool,
) (*T, error) {
	// Type assert Arguments from any to map[string]interface{}
	args, ok := request.Arguments.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid arguments type")
	}

	val, ok := args[name]
	if !ok || val == nil {
		if required {
			return nil, errors.New("missing required parameter: " + name)
		}
		return nil, nil // Not an error for optional params
	}

	var result T
	data, err := json.Marshal(val)
	if err != nil {
		return nil, errors.New("invalid parameter type: " + name)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.New("invalid parameter type: " + name)
	}

	return &result, nil
}

// Generic validation functions

// validateAndAddRequired validates and adds a required parameter of any type
func validateAndAddRequired[T any](
	v *Validator,
	params map[string]interface{},
	name string,
) *Validator {
	value, err := extractValueGeneric[T](v.request, name, true)
	if err != nil {
		return v.addError(err)
	}

	if value == nil {
		return v
	}

	params[name] = *value
	return v
}

// validateAndAddOptional validates and adds an optional parameter of any type
// if not empty
func validateAndAddOptional[T any](
	v *Validator,
	params map[string]interface{},
	name string,
) *Validator {
	value, err := extractValueGeneric[T](v.request, name, false)
	if err != nil {
		return v.addError(err)
	}

	if value == nil {
		return v
	}

	params[name] = *value

	return v
}

// validateAndAddToPath is a generic helper to extract a value and write it into
// `target[targetKey]` if non-empty
func validateAndAddToPath[T any](
	v *Validator,
	target map[string]interface{},
	paramName string,
	targetKey string,
) *Validator {
	value, err := extractValueGeneric[T](v.request, paramName, false)
	if err != nil {
		return v.addError(err)
	}

	if value == nil {
		return v
	}

	target[targetKey] = *value

	return v
}

// ValidateAndAddOptionalStringToPath validates an optional string
// and writes it into target[targetKey]
func (v *Validator) ValidateAndAddOptionalStringToPath(
	target map[string]interface{},
	paramName, targetKey string,
) *Validator {
	return validateAndAddToPath[string](v, target, paramName, targetKey) // nolint:lll
}

// ValidateAndAddOptionalBoolToPath validates an optional bool
// and writes it into target[targetKey]
// only if it was explicitly provided in the request
func (v *Validator) ValidateAndAddOptionalBoolToPath(
	target map[string]interface{},
	paramName, targetKey string,
) *Validator {
	// Now validate and add the parameter
	value, err := extractValueGeneric[bool](v.request, paramName, false)
	if err != nil {
		return v.addError(err)
	}

	if value == nil {
		return v
	}

	target[targetKey] = *value
	return v
}

// ValidateAndAddOptionalIntToPath validates an optional integer
// and writes it into target[targetKey]
func (v *Validator) ValidateAndAddOptionalIntToPath(
	target map[string]interface{},
	paramName, targetKey string,
) *Validator {
	return validateAndAddToPath[int64](v, target, paramName, targetKey)
}

// Type-specific validator methods

// ValidateAndAddRequiredString validates and adds a required string parameter
func (v *Validator) ValidateAndAddRequiredString(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddRequired[string](v, params, name)
}

// ValidateAndAddOptionalString validates and adds an optional string parameter
func (v *Validator) ValidateAndAddOptionalString(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddOptional[string](v, params, name)
}

// ValidateAndAddRequiredMap validates and adds a required map parameter
func (v *Validator) ValidateAndAddRequiredMap(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddRequired[map[string]interface{}](v, params, name)
}

// ValidateAndAddOptionalMap validates and adds an optional map parameter
func (v *Validator) ValidateAndAddOptionalMap(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddOptional[map[string]interface{}](v, params, name)
}

// ValidateAndAddRequiredArray validates and adds a required array parameter
func (v *Validator) ValidateAndAddRequiredArray(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddRequired[[]interface{}](v, params, name)
}

// ValidateAndAddOptionalArray validates and adds an optional array parameter
func (v *Validator) ValidateAndAddOptionalArray(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddOptional[[]interface{}](v, params, name)
}

// ValidateAndAddPagination validates and adds pagination parameters
// (count and skip)
func (v *Validator) ValidateAndAddPagination(
	params map[string]interface{},
) *Validator {
	return v.ValidateAndAddOptionalInt(params, "count").
		ValidateAndAddOptionalInt(params, "skip")
}

// ValidateAndAddExpand validates and adds expand parameters
func (v *Validator) ValidateAndAddExpand(
	params map[string]interface{},
) *Validator {
	expand, err := extractValueGeneric[[]string](v.request, "expand", false)
	if err != nil {
		return v.addError(err)
	}

	if expand == nil {
		return v
	}

	if len(*expand) > 0 {
		for _, val := range *expand {
			params["expand[]"] = val
		}
	}
	return v
}

// ValidateAndAddRequiredInt validates and adds a required integer parameter
func (v *Validator) ValidateAndAddRequiredInt(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddRequired[int64](v, params, name)
}

// ValidateAndAddOptionalInt validates and adds an optional integer parameter
func (v *Validator) ValidateAndAddOptionalInt(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddOptional[int64](v, params, name)
}

// ValidateAndAddRequiredFloat validates and adds a required float parameter
func (v *Validator) ValidateAndAddRequiredFloat(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddRequired[float64](v, params, name)
}

// ValidateAndAddOptionalFloat validates and adds an optional float parameter
func (v *Validator) ValidateAndAddOptionalFloat(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddOptional[float64](v, params, name)
}

// ValidateAndAddRequiredBool validates and adds a required boolean parameter
func (v *Validator) ValidateAndAddRequiredBool(
	params map[string]interface{},
	name string,
) *Validator {
	return validateAndAddRequired[bool](v, params, name)
}

// ValidateAndAddOptionalBool validates and adds an optional boolean parameter
// Note: This adds the boolean value only
// if it was explicitly provided in the request
func (v *Validator) ValidateAndAddOptionalBool(
	params map[string]interface{},
	name string,
) *Validator {
	// Now validate and add the parameter
	value, err := extractValueGeneric[bool](v.request, name, false)
	if err != nil {
		return v.addError(err)
	}

	if value == nil {
		return v
	}

	params[name] = *value
	return v
}

// Parameter conversion utility functions

// ExtractOptionalString extracts an optional string parameter from a map
// and returns a pointer to it if present
func ExtractOptionalString(params map[string]interface{}, key string) *string {
	if val, ok := params[key]; ok {
		if strVal, ok := val.(string); ok {
			return &strVal
		}
	}
	return nil
}

// ExtractOptionalInt extracts an optional integer parameter from a map
// and returns a pointer to it if present, converting from int64 to int
func ExtractOptionalInt(params map[string]interface{}, key string) *int {
	if val, ok := params[key]; ok {
		if intVal, ok := val.(int64); ok {
			intPtr := int(intVal)
			return &intPtr
		}
	}
	return nil
}

// ExtractOptionalInt64 extracts an optional int64 parameter from a map
// and returns a pointer to it if present
func ExtractOptionalInt64(params map[string]interface{}, key string) *int64 {
	if val, ok := params[key]; ok {
		if intVal, ok := val.(int64); ok {
			return &intVal
		}
	}
	return nil
}

// ExtractOptionalFloat64 extracts an optional float64 parameter from a map
// and returns a pointer to it if present
func ExtractOptionalFloat64(params map[string]interface{}, key string) *float64 {
	if val, ok := params[key]; ok {
		if floatVal, ok := val.(float64); ok {
			return &floatVal
		}
	}
	return nil
}

// ExtractOptionalBool extracts an optional boolean parameter from a map
// and returns a pointer to it if present
func ExtractOptionalBool(params map[string]interface{}, key string) *bool {
	if val, ok := params[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return &boolVal
		}
	}
	return nil
}

// SetOptionalString sets a string pointer field if the value exists in params
func SetOptionalString(params map[string]interface{}, key string, target **string) {
	*target = ExtractOptionalString(params, key)
}

// SetOptionalInt sets an int pointer field if the value exists in params
func SetOptionalInt(params map[string]interface{}, key string, target **int) {
	*target = ExtractOptionalInt(params, key)
}

// SetOptionalInt64 sets an int64 pointer field if the value exists in params
func SetOptionalInt64(params map[string]interface{}, key string, target **int64) {
	*target = ExtractOptionalInt64(params, key)
}

// SetOptionalFloat64 sets a float64 pointer field if the value exists in params
func SetOptionalFloat64(params map[string]interface{}, key string, target **float64) {
	*target = ExtractOptionalFloat64(params, key)
}

// SetOptionalBool sets a bool pointer field if the value exists in params
func SetOptionalBool(params map[string]interface{}, key string, target **bool) {
	*target = ExtractOptionalBool(params, key)
}

// ParameterMapping represents a mapping between a parameter key and a target field
type ParameterMapping struct {
	Key    string
	Target interface{} // Should be **T where T is the target type
	Type   string      // "string", "int", "int64", "float64", "bool"
}

// SetOptionalParameters sets multiple optional parameters using a mapping configuration
func SetOptionalParameters(params map[string]interface{}, mappings []ParameterMapping) {
	for _, mapping := range mappings {
		switch mapping.Type {
		case "string":
			if target, ok := mapping.Target.(**string); ok {
				SetOptionalString(params, mapping.Key, target)
			}
		case "int":
			if target, ok := mapping.Target.(**int); ok {
				SetOptionalInt(params, mapping.Key, target)
			}
		case "int64":
			if target, ok := mapping.Target.(**int64); ok {
				SetOptionalInt64(params, mapping.Key, target)
			}
		case "float64":
			if target, ok := mapping.Target.(**float64); ok {
				SetOptionalFloat64(params, mapping.Key, target)
			}
		case "bool":
			if target, ok := mapping.Target.(**bool); ok {
				SetOptionalBool(params, mapping.Key, target)
			}
		}
	}
}
