package main

import (
	"encoding/json"
	"sort"
	"time"
)

type invalidOperation struct {
	ID   interface{}
	date time.Time
}

type balance struct {
	value float64
}

type companyStatistic struct {
	Company              string             `json:"company"`
	ValidOperationsCount int                `json:"valid_operations_count"`
	Balance              balance            `json:"balance"`
	InvalidOperations    []invalidOperation `json:"invalid_operations,omitempty"`
}

func (o invalidOperation) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	jsonData, err := json.Marshal(o.ID)
	return jsonData, err
}

func (b balance) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(int(b.value))
	return jsonData, err
}

func (c *companyStatistic) sortInvalidOperations() {
	sort.SliceStable(c.InvalidOperations, func(i, j int) bool {
		return c.InvalidOperations[i].date.Before(c.InvalidOperations[j].date)
	})
}

func (c *companyStatistic) updateStatistic(companyBilling billing) {
	if companyBilling.invalid {
		c.InvalidOperations = append(c.InvalidOperations, invalidOperation{ID: companyBilling.id})
	} else {
		c.ValidOperationsCount++
		c.updateBalance(companyBilling)
	}
}

func (c *companyStatistic) updateBalance(companyBilling billing) {
	switch companyBilling.bType {
	case income, plus:
		c.Balance.value += companyBilling.value
	default:
		c.Balance.value -= companyBilling.value
	}
}

func calculateCompaniesStatistic(billingsStatistic billings) []companyStatistic {
	companies := make(map[string]*companyStatistic)
	for _, billing := range billingsStatistic {
		if _, ok := companies[billing.company]; !ok {
			companies[billing.company] = &companyStatistic{
				Company:              billing.company,
				ValidOperationsCount: 0,
				Balance:              balance{value: 0},
				InvalidOperations:    []invalidOperation{},
			}
		}
		companies[billing.company].updateStatistic(billing)
	}
	return sortStatistic(companies)
}

func sortStatistic(companies map[string]*companyStatistic) []companyStatistic {
	statistic := make([]companyStatistic, 0, len(companies))
	for _, company := range companies {
		company.sortInvalidOperations()
		statistic = append(statistic, *company)
	}
	sort.SliceStable(statistic, func(i, j int) bool {
		return statistic[i].Company < statistic[j].Company
	})
	return statistic
}
