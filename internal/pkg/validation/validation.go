package validation

//
//import "fmt"
//
//type VaidationError struct {
//	Field   string
//	Message string
//}
//
//func (e VaidationError) Error() string {
//
//	return fmt.Sprintf("%s: %s", e.Field, e.Message)
//}
//
//type Validator struct {
//	errors []VaidationError
//}
//
//func NewValidator() *Validator {
//	return &Validator{errors: make([]VaidationError, 0)}
//}
//func (v *Validator) ValidateRequired(fieldName string, value interface{}) *Validator {
//	if isEmpty(value) {
//		v.errors = append(v.errors, VaidationError{Field: fieldName, Message: "value is empty"})
//	}
//	return v
//}
//func (v *Validator ) ValidatePositiv
