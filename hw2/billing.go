package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type operation struct {
	ID        interface{} `json:"id,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	Type      string      `json:"type,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
}

type billingRaw struct {
	Company   string    `json:"company,omitempty"`
	Operation operation `json:"operation,omitempty"`
	operation
}

type unwrappedBillingRaw struct {
	company   string
	id        interface{}
	value     interface{}
	bType     string
	createdAt string
}

type billing struct {
	company   string
	bType     string
	value     float64
	id        interface{}
	createdAt time.Time
	invalid   bool
}

var (
	ErrInvalidBillingType      = errors.New("invalid billing type")
	ErrInvalidBillingValue     = errors.New("invalid billing value")
	ErrInvalidBillingValueType = errors.New("invalid billing value type")
	ErrUnknownBillingIDType    = errors.New("unknown billing id type")
	ErrUnknownBillingTime      = errors.New("unknown type of billing created_at field")
	ErrEmptyCompany            = errors.New("empty company")
	ErrValidateParsedBilling   = errors.New("validate parsed billing")
)

const (
	plus    = "+"
	minus   = "-"
	income  = "income"
	outcome = "outcome"
)

func (bRaw billingRaw) toUnwrappedBillingRaw() unwrappedBillingRaw {
	var parsedB unwrappedBillingRaw
	parsedB.company = bRaw.Company
	if bRaw.ID != nil {
		parsedB.id = bRaw.ID
	} else if bRaw.Operation.ID != nil {
		parsedB.id = bRaw.Operation.ID
	}
	if bRaw.Type != "" {
		parsedB.bType = bRaw.Type
	} else if bRaw.Operation.Type != "" {
		parsedB.bType = bRaw.Operation.Type
	}
	if bRaw.Value != nil {
		parsedB.value = bRaw.Value
	} else if bRaw.Operation.Value != nil {
		parsedB.value = bRaw.Operation.Value
	}
	if bRaw.CreatedAt != "" {
		parsedB.createdAt = bRaw.CreatedAt
	} else if bRaw.Operation.CreatedAt != "" {
		parsedB.createdAt = bRaw.Operation.CreatedAt
	}
	return parsedB
}

func (b unwrappedBillingRaw) validateCompany() (string, error) {
	if b.company == "" {
		return "", ErrEmptyCompany
	}
	return b.company, nil
}

func (b unwrappedBillingRaw) validateID() (interface{}, error) {
	switch id := b.id.(type) {
	case string:
		return id, nil
	case float64:
		if id == float64(int(id)) {
			return id, nil
		}
		return nil, ErrUnknownBillingIDType
	default:
		return nil, ErrUnknownBillingIDType
	}
}

func (b unwrappedBillingRaw) validateType() (string, error) {
	switch b.bType {
	case plus, minus, income, outcome:
		return b.bType, nil
	default:
		return "", ErrInvalidBillingType
	}
}

func (b unwrappedBillingRaw) validateTime() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, b.createdAt)
	if err != nil {
		return time.Time{}, ErrUnknownBillingTime
	}
	return t, nil
}

func (b unwrappedBillingRaw) validateValue() (float64, error) {
	switch val := b.value.(type) {
	case string:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, ErrInvalidBillingValue
		}
		return floatVal, nil
	case float64:
		return val, nil
	default:
		return 0, ErrInvalidBillingValueType
	}
}

func (b unwrappedBillingRaw) validate() (billing, error) {
	var returnBill billing
	company, err := b.validateCompany()
	if err != nil {
		return billing{}, fmt.Errorf("%v: %w", ErrValidateParsedBilling, err)
	}
	returnBill.company = company

	id, err := b.validateID()
	if err != nil {
		return billing{}, fmt.Errorf("%v: %w", ErrValidateParsedBilling, err)
	}
	returnBill.id = id

	createdAt, err := b.validateTime()
	if err != nil {
		return billing{}, fmt.Errorf("%v: %w", ErrValidateParsedBilling, err)
	}
	returnBill.createdAt = createdAt

	value, err := b.validateValue()
	if err != nil {
		returnBill.invalid = true
	}
	returnBill.value = value

	bType, err := b.validateType()
	if err != nil {
		returnBill.invalid = true
	}
	returnBill.bType = bType

	return returnBill, nil
}
