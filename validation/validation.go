/*
 * MIT License
 *
 * Copyright (c) 2023 WeCoding.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package validation

import (
	"fmt"
	"os"
	"reflect"

	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"

	"github.com/wecoding.top/common/validation/field"
)

const (
	maxDescriptionLength = 255
)

// Validator is a custom validator for configs.
type Validator struct {
	val   *validator.Validate
	data  interface{}
	trans ut.Translator
}

// NewValidator creates a new Validator.
func NewValidator(data interface{}) *Validator {
	result := validator.New()

	// independent validators
	result.RegisterValidation("dir", validateDir)                 // nolint: errcheck // no need
	result.RegisterValidation("file", validateFile)               // nolint: errcheck // no need
	result.RegisterValidation("description", validateDescription) // nolint: errcheck // no need
	result.RegisterValidation("name", validateName)               // nolint: errcheck // no need

	// default translations
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	err := en.RegisterDefaultTranslations(result, trans)
	if err != nil {
		panic(err)
	}

	// additional translations
	translations := []struct {
		tag         string
		translation string
	}{
		{
			tag:         "dir",
			translation: "{0} must point to an existing directory, but found '{1}'",
		},
		{
			tag:         "file",
			translation: "{0} must point to an existing file, but found '{1}'",
		},
		{
			tag:         "description",
			translation: fmt.Sprintf("must be less than %d", maxDescriptionLength),
		},
		{
			tag:         "name",
			translation: "is not a invalid name",
		},
	}
	for _, t := range translations {
		err = result.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation), translateFunc)
		if err != nil {
			panic(err)
		}
	}

	return &Validator{
		val:   result,
		data:  data,
		trans: trans,
	}
}

func registrationFunc(tag string, translation string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, true); err != nil {
			return
		}

		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), reflect.ValueOf(fe.Value()).String())
	if err != nil {
		return fe.(error).Error()
	}

	return t
}

// Validate validates config for errors and returns an error (it can be casted to
// ValidationErrors, containing a list of errors inside). When error is printed as string, it will
// automatically contains the full list of validation errors.
func (v *Validator) Validate() field.ErrorList {
	// validate policy
	err := v.val.Struct(v.data)
	if err == nil {
		return nil
	}

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return field.ErrorList{field.Invalid(field.NewPath(""), err.Error(), "")}
	}

	allErrs := field.ErrorList{}

	// collect human-readable errors
	vErrors, _ := err.(validator.ValidationErrors)
	for _, vErr := range vErrors {
		allErrs = append(allErrs, field.Invalid(field.NewPath(vErr.Namespace()), vErr.Translate(v.trans), ""))
	}

	return allErrs
}

// validateDir checks if a given string is an existing directory.
func validateDir(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

// validateFile checks if a given string is an existing file.
func validateFile(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
		return true
	}

	return false
}

// validateDescription checks if a given description is illegal.
func validateDescription(fl validator.FieldLevel) bool {
	description := fl.Field().String()

	return len(description) <= maxDescriptionLength
}

// validateName checks if a given name is illegal.
func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if errs := IsQualifiedName(name); len(errs) > 0 {
		return false
	}

	return true
}
