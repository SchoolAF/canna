package validator

import "github.com/asaskevich/govalidator"

func ValidateData(data interface{}) error {
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return err
	}
	return nil
}
