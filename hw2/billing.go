package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type validator interface {
	isValid() error
}

type billingCompany string
type billingType string
type billingValue float64
type billingID string
type billingCreatedAt time.Time

type operation struct {
	Value     billingValue     `json:"value"`
	Type      billingType      `json:"type"`
	ID        billingID        `json:"id"`
	CreatedAt billingCreatedAt `json:"created_at"`
}

type billingRaw struct {
	Company   billingCompany `json:"company"`
	Operation operation      `json:"operation"`
	operation
}

type billing struct {
	company   billingCompany
	bType     billingType
	value     billingValue
	id        billingID
	createdAt billingCreatedAt
	invalid   bool
}

type billingJSONFields map[string]interface{}
type billingTypeValidValue string

var (
	ErrBillingCompany        = errors.New("billing company")
	ErrBillingType           = errors.New("billing type")
	ErrInvalidBillingType    = errors.New("invalid billing type passed")
	ErrBillingValue          = errors.New("billing value")
	ErrBillingID             = errors.New("billing ID")
	ErrBillingCreatedAt      = errors.New("billing created at")
	ErrJSONEmptyString       = errors.New("json: unable to unmarshal, value is empty string")
	ErrSetupBillingField     = errors.New("setup billing field")
	ErrUnmarshalToNumber     = errors.New("unmarshal to number")
	ErrUnmarshalToString     = errors.New("unmarshal to string")
	ErrPassedValueNotPointer = errors.New("passed value is not a pointer")
	ErrSkipBilling           = errors.New("skip billing")
	ErrInvalidBilling        = errors.New("invalid error")
)

const (
	companyField   = "company"
	typeField      = "type"
	valueField     = "value"
	idField        = "id"
	createdAtFiled = "created_at"
	operationField = "operation"
)

const QuoteByte = 34

const (
	income  billingTypeValidValue = "income"
	outcome billingTypeValidValue = "outcome"
	plus    billingTypeValidValue = "+"
	minus   billingTypeValidValue = "-"
)

func (b *billingCompany) UnmarshalJSON(data []byte) error {
	strCompany, err := unmarshalToString(data)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBillingCompany, err)
	}
	if len(strCompany) == 0 {
		return fmt.Errorf("%s: %w", ErrBillingCompany, ErrJSONEmptyString)
	}
	*b = billingCompany(strCompany)
	return nil
}

func (b billingCompany) isValid() error {
	if b == "" {
		return ErrSkipBilling
	}
	return nil
}

func (b *billingCreatedAt) UnmarshalJSON(data []byte) error {
	strDate, err := unmarshalToString(data)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBillingCreatedAt, err)
	}
	t, err := time.Parse(time.RFC3339, strDate)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBillingCreatedAt, err)
	}
	*b = billingCreatedAt(t)
	return nil
}

func (b billingCreatedAt) isZero() bool {
	return time.Time(b) == time.Time{}
}

func (b billingCreatedAt) isValid() error {
	if b.isZero() {
		return ErrSkipBilling
	}
	return nil
}

func (b *billingType) UnmarshalJSON(data []byte) error {
	strType, err := unmarshalToString(data)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBillingType, err)
	}
	switch billingTypeValidValue(strType) {
	case income, outcome, plus, minus:
		*b = billingType(strType)
	default:
		return fmt.Errorf("%s: %w - %s", ErrBillingType, ErrInvalidBillingType, strType)
	}
	return nil
}

func (b billingType) isValid() error {
	if b == "" {
		return ErrInvalidBilling
	}
	return nil
}

func (b *billingValue) UnmarshalJSON(data []byte) error {
	var floatValue float64
	err := unmarshalToNumber(data, &floatValue)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBillingValue, err)
	}
	*b = billingValue(floatValue)
	return nil
}

func (b billingValue) isValid() error {
	if b == 0 {
		return ErrInvalidBilling
	}
	return nil
}

func (b *billingID) UnmarshalJSON(data []byte) error {
	var strID string
	var intID int
	if data[0] == QuoteByte {
		err := json.Unmarshal(data, &strID)
		if err != nil {
			return fmt.Errorf("%s: %w", ErrBillingID, err)
		}
		*b = billingID(strID)
	} else {
		err := json.Unmarshal(data, &intID)
		if err != nil {
			return fmt.Errorf("%s: %w", ErrBillingID, err)
		}
		*b = billingID(strconv.Itoa(intID))
	}
	return nil
}

func (b billingID) isValid() error {
	if b == "" {
		return ErrSkipBilling
	}
	return nil
}

func newBilling(fields billingJSONFields) (billing, error) {
	var bRaw billingRaw
	fields.setupBillingField(companyField, &bRaw.Company)
	fields.setupBillingField(valueField, &bRaw.Value)
	fields.setupBillingField(typeField, &bRaw.Type)
	fields.setupBillingField(idField, &bRaw.ID)
	fields.setupBillingField(createdAtFiled, &bRaw.CreatedAt)
	fields.setupBillingField(operationField, &bRaw.Operation)
	b := bRaw.toBilling()
	err := b.validate()
	if err != nil {
		errors.Is(err, ErrInvalidBilling)
		b.invalid = true
	}
	return b, err
}

func (b billing) validate() error {
	validators := []validator{b.company, b.bType, b.value, b.id, b.createdAt}
	var lastError error
	for _, val := range validators {
		if err := val.isValid(); err != nil {
			switch err {
			case ErrSkipBilling:
				return err
			case ErrInvalidBilling:
				lastError = err
			}
		}
	}
	return lastError
}

func (fields billingJSONFields) setupBillingField(fieldName string, billingField interface{}) {
	val, ok := fields[fieldName]
	if !ok {
		return
	}
	data, err := json.Marshal(val)
	if err != nil {
		log.Warn(fmt.Errorf("%v: %w", ErrSetupBillingField, err))
	}
	if err := json.Unmarshal(data, billingField); err != nil {
		log.Warn(fmt.Errorf("%v: %w", ErrSetupBillingField, err))
	}
}

func (bRaw billingRaw) toBilling() billing {
	var b billing
	b.company = bRaw.Company
	if bRaw.Type != "" {
		b.bType = bRaw.Type
	} else if bRaw.Operation.Type != "" {
		b.bType = bRaw.Operation.Type
	}
	if bRaw.Value != 0 {
		b.value = bRaw.Value
	} else if bRaw.Operation.Value != 0 {
		b.value = bRaw.Operation.Value
	}
	if bRaw.ID != "" {
		b.id = bRaw.ID
	} else if bRaw.Operation.ID != "" {
		b.id = bRaw.Operation.ID
	}
	if !bRaw.CreatedAt.isZero() {
		b.createdAt = bRaw.CreatedAt
	} else if !bRaw.Operation.CreatedAt.isZero() {
		b.createdAt = bRaw.Operation.CreatedAt
	}
	return b
}

func unmarshalToString(data []byte) (string, error) {
	var strType string
	err := json.Unmarshal(data, &strType)
	if err != nil {
		return "", fmt.Errorf("%v: %w", ErrUnmarshalToString, err)
	}
	strType = strings.TrimSpace(strType)
	return strType, nil
}

func unmarshalToNumber(data []byte, number interface{}) error {
	switch number.(type) {
	case *int:
	case *float64:
	default:
		return fmt.Errorf("%v: %v", ErrUnmarshalToNumber, ErrPassedValueNotPointer)
	}
	if data[0] == QuoteByte {
		err := json.Unmarshal(data[1:len(data)-1], number)
		if err != nil {
			return fmt.Errorf("%v: %w", ErrUnmarshalToNumber, err)
		}
	} else {
		err := json.Unmarshal(data, number)
		if err != nil {
			return fmt.Errorf("%v: %w", ErrUnmarshalToNumber, err)
		}
	}
	return nil
}
