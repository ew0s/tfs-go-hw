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
	skip      bool
}

type billingTypeValidValue string
type jsonFiled string

const QuoteByte = 34

const (
	income  billingTypeValidValue = "income"
	outcome billingTypeValidValue = "outcome"
	plus    billingTypeValidValue = "+"
	minus   billingTypeValidValue = "-"
)

const (
	companyField   jsonFiled = "company"
	createdAtField jsonFiled = "created_at"
	typeField      jsonFiled = "type"
	valueField     jsonFiled = "value"
	idField        jsonFiled = "id"
	operationField jsonFiled = "operation"
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
		return fmt.Errorf("billingType: invalid billing type passed")
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
	unwrappedData, err := b.unwrapBillingData(data)
	if err != nil {
		log.Println("billing: ", err)
	}
	b.setBillingFields(unwrappedData, companyField)
	b.setBillingFields(unwrappedData, typeField)
	b.setBillingFields(unwrappedData, valueField)
	b.setBillingFields(unwrappedData, idField)
	b.setBillingFields(unwrappedData, createdAtField)
	return nil
}

func (b *billing) setBillingFields(data map[string]*json.RawMessage, billingField jsonFiled) {
	value, ok := data[string(billingField)]
	if !ok {
		return
	}
	switch billingField {
	case companyField:
		if err := unmarshalValue(*value, &b.Company); err != nil {
			b.skip = true
			log.Println(err)
		}
	case typeField:
		if err := unmarshalValue(*value, &b.Type); err != nil {
			b.invalid = true
			log.Println(err)
		}
	case valueField:
		if err := unmarshalValue(*value, &b.Value); err != nil {
			b.invalid = true
			log.Println(err)
		}
	case idField:
		if err := unmarshalValue(*value, &b.ID); err != nil {
			b.skip = true
			log.Println(err)
		}
	case createdAtField:
		if err := unmarshalValue(*value, &b.CreatedAt); err != nil {
			b.skip = true
			log.Println(err)
		}
	}
	b.toSkip(data)
	b.toSetInvalid(data)
}

func (b *billing) toSkip(data map[string]*json.RawMessage) {
	if _, ok := data[string(companyField)]; !ok {
		b.skip = true
	}
	if _, ok := data[string(createdAtField)]; !ok {
		b.skip = true
	}
	if _, ok := data[string(idField)]; !ok {
		b.skip = true
	}
}

func (b *billing) toSetInvalid(data map[string]*json.RawMessage) {
	if _, ok := data[string(typeField)]; !ok {
		b.invalid = true
	}
	if _, ok := data[string(valueField)]; !ok {
		b.invalid = true
	}
}

func unmarshalValue(data []byte, fieldType interface{}) error {
	if err := json.Unmarshal(data, fieldType); err != nil {
		return fmt.Errorf("unmarshalValue: %w", err)
	}
	return nil
}

func (b billing) unwrapBillingData(data []byte) (map[string]*json.RawMessage, error) {
	var jsonBilling map[string]*json.RawMessage
	if err := json.Unmarshal(data, &jsonBilling); err != nil {
		return nil, fmt.Errorf("unwrapBillingData: %w", err)
	}
	if value, ok := jsonBilling[string(operationField)]; ok {
		var operationJSON map[string]*json.RawMessage
		if err := json.Unmarshal(*value, &operationJSON); err != nil {
			return nil, fmt.Errorf("unwrapBillingData: unwrap operation: %w", err)
		}
		for k, v := range operationJSON {
			jsonBilling[k] = v
		}
		delete(jsonBilling, string(operationField))
	}
	return jsonBilling, nil
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
