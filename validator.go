package ice

import (
	validator "github.com/asaskevich/govalidator"
)

type RequestValidator interface {
	Validate(req interface{})
}

type Validator struct {
	IsValid bool              `json:"-"`
	Errors  map[string]string `json:"-"`
}

func (v *Validator) Validate(req interface{}) {
	_, err := validator.ValidateStruct(req)
	if err == nil {
		v.IsValid = true
		v.Errors = make(map[string]string)
	} else {
		v.Errors = validator.ErrorsByField(err)
	}
}

func (v Validator) Execute(conn Conn) {
	if v.Errors != nil {
		conn.SendErrors("", v.Errors)
	}
}

func (v *Validator) AddError(field, error string) *Validator {
	v.Errors[field] = error
	return v
}
