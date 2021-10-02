package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type billingCompany struct {
	company string
}

type billingCreatedAt struct {
	value time.Time
}

type billingType struct {
	value string
}

type billingValue struct {
	value float64
}

type billingID struct {
	value string
}

type billing struct {
	Company   billingCompany   `json:"company"`
	CreatedAt billingCreatedAt `json:"created_at"`
	Type      billingType      `json:"type"`
	Value     billingValue     `json:"value"`
	ID        billingID        `json:"id"`
	invalid   bool
}

type billingJSONFields map[string]interface{}
type billingTypeValidValue string

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
		return fmt.Errorf("billingCompany: %w", err)
	}
	if len(strCompany) == 0 {
		return fmt.Errorf("billingCompany: json: unable to unmarshal, value is empty string")
	}
	b.company = strCompany
	return nil
}

func (b *billingCreatedAt) UnmarshalJSON(data []byte) error {
	strDate, err := unmarshalToString(data)
	if err != nil {
		return fmt.Errorf("billingCreatedAt: %w", err)
	}
	t, err := time.Parse(time.RFC3339, strDate)
	if err != nil {
		return fmt.Errorf("billingCreditAt: %w", err)
	}
	b.value = t
	return nil
}

func (b *billingType) UnmarshalJSON(data []byte) error {
	strType, err := unmarshalToString(data)
	if err != nil {
		return fmt.Errorf("billingType: %w", err)
	}
	switch billingTypeValidValue(strType) {
	case income, outcome, plus, minus:
		b.value = strType
	default:
		return fmt.Errorf("billingType: invalid billing type passed - %s", strType)
	}
	return nil
}

func (b *billingValue) UnmarshalJSON(data []byte) error {
	var floatValue float64
	err := unmarshalToNumber(data, &floatValue)
	if err != nil {
		return fmt.Errorf("billingValue: %w", err)
	}
	b.value = floatValue
	return nil
}

func (b *billingID) UnmarshalJSON(data []byte) error {
	var strID string
	var intID int
	if data[0] == QuoteByte {
		err := json.Unmarshal(data, &strID)
		if err != nil {
			return fmt.Errorf("unmarshallToNumber: %w", err)
		}
		b.value = strID
	} else {
		err := json.Unmarshal(data, &intID)
		if err != nil {
			return fmt.Errorf("unmarshallToNumber: %w", err)
		}
		b.value = strconv.Itoa(intID)
	}
	return nil
}

func (b *billing) UnmarshalJSON(data []byte) error {
	var billingJSONMap billingJSONFields
	if err := json.Unmarshal(data, &billingJSONMap); err != nil {
		return fmt.Errorf("billing unmarshal: %w", err)
	}
	billingJSONMap.unwrapOperationField()
	*b = billingJSONMap.parseFields()
	return nil
}

func (b *billing) validate() error {
	if b.Company == (billingCompany{}) {
		return fmt.Errorf("company field is invalid")
	}
	if b.ID == (billingID{}) {
		return fmt.Errorf("id field is invalid")
	}
	if b.CreatedAt == (billingCreatedAt{}) {
		return fmt.Errorf("created_at field is invalid")
	}
	if b.Type == (billingType{}) {
		b.invalid = true
	}
	if b.Value == (billingValue{}) {
		b.invalid = true
	}
	return nil
}

func (fields billingJSONFields) unwrapOperationField() {
	if operation, ok := fields[operationField].(map[string]interface{}); ok {
		for key, value := range operation {
			fields[key] = value
		}
		delete(fields, operationField)
	}
}

func (fields billingJSONFields) parseFields() billing {
	var b billing
	fields.parseField(companyField, &b.Company)
	fields.parseField(typeField, &b.Type)
	fields.parseField(valueField, &b.Value)
	fields.parseField(idField, &b.ID)
	fields.parseField(createdAtFiled, &b.CreatedAt)
	return b
}

func (fields billingJSONFields) parseField(fieldName string, billingField interface{}) {
	if data, err := fields.getByteField(fieldName); err == nil {
		unmarshalValue(data, billingField)
	}
}

func (fields billingJSONFields) getByteField(name string) ([]byte, error) {
	if value, ok := fields[name]; ok {
		bytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("get byte field: %w", err)
		}
		return bytes, nil
	}
	return nil, fmt.Errorf("get byte field: no such field: %s", name)
}

func unmarshalValue(data []byte, fieldType interface{}) {
	if err := json.Unmarshal(data, fieldType); err != nil {
		log.Println(err)
	}
}

func unmarshalToString(data []byte) (string, error) {
	var strType string
	err := json.Unmarshal(data, &strType)
	if err != nil {
		return "", fmt.Errorf("unmarshalToString: %w", err)
	}
	strType = strings.TrimSpace(strType)
	return strType, nil
}

func unmarshalToNumber(data []byte, number interface{}) error {
	switch number.(type) {
	case *int:
	case *float64:
	default:
		return fmt.Errorf("unmarshallToNumber: passed value should be a pointer")
	}
	if data[0] == QuoteByte {
		err := json.Unmarshal(data[1:len(data)-1], number)
		if err != nil {
			return fmt.Errorf("unmarshallToNumber: %w", err)
		}
	} else {
		err := json.Unmarshal(data, number)
		if err != nil {
			return fmt.Errorf("unmarshallToNumber: %w", err)
		}
	}
	return nil
}
