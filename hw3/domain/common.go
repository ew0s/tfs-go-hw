package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Price struct {
	Ticker string
	Value  float64
	TS     time.Time
}

var (
	ErrUnknownPeriod                = errors.New("unknown period")
	ErrCreateNewCandle              = errors.New("new candle")
	ErrUpdateCandlesMap             = errors.New("update candles map")
	ErrUpdateCandle                 = errors.New("update candle")
	ErrUpdateCandleMismatchedPeriod = errors.New("update candle: mismatch period")
)

type CandlePeriod string
type CandleMap map[string]*Candle

const (
	CandlePeriod1m  CandlePeriod = "1m"
	CandlePeriod2m  CandlePeriod = "2m"
	CandlePeriod10m CandlePeriod = "10m"
)

func PeriodTS(period CandlePeriod, ts time.Time) (time.Time, error) {
	switch period {
	case CandlePeriod1m:
		return ts.Truncate(time.Minute), nil
	case CandlePeriod2m:
		return ts.Truncate(2 * time.Minute), nil
	case CandlePeriod10m:
		return ts.Truncate(10 * time.Minute), nil
	default:
		return time.Time{}, ErrUnknownPeriod
	}
}

type Candle struct {
	Ticker string
	Period CandlePeriod // Интервал
	Open   float64      // Цена открытия
	High   float64      // Максимальная цена
	Low    float64      // Минимальная цена
	Close  float64      // Цена закрытие
	TS     time.Time    // Время начала интервала
}

func (cm CandleMap) UpdateFromPrice(price Price, period CandlePeriod) (Candle, error) {
	val, ok := cm[price.Ticker]
	if !ok {
		if newCandle, err := NewCandleFromPrice(period, price); err != nil {
			return Candle{}, fmt.Errorf("%v: %v", ErrUpdateCandlesMap, err)
		} else {
			cm[price.Ticker] = &newCandle
		}
		return Candle{}, nil
	}
	if err := val.UpdateFromPrice(price); err != nil {
		if errors.Is(err, ErrUpdateCandleMismatchedPeriod) {
			closedCandle := *val
			if newCandle, err := NewCandleFromPrice(period, price); err != nil {
				return closedCandle, fmt.Errorf("%v: %w", ErrUpdateCandlesMap, err)
			} else {
				cm[price.Ticker] = &newCandle
				return closedCandle, ErrUpdateCandleMismatchedPeriod
			}
		}
	}
	return Candle{}, nil
}

func (cm CandleMap) UpdateFromCandle(candle Candle, period CandlePeriod) (Candle, error) {
	val, ok := cm[candle.Ticker]
	if !ok {
		if newCandle, err := NewCandleFromCandle(period, candle); err != nil {
			return Candle{}, fmt.Errorf("%v: %w", ErrUpdateCandlesMap, err)
		} else {
			cm[candle.Ticker] = &newCandle
			return Candle{}, nil
		}
	}
	if err := val.UpdateFromCandle(candle); err != nil {
		if errors.Is(err, ErrUpdateCandleMismatchedPeriod) {
			closedCandle := *val
			if newCandle, err := NewCandleFromCandle(period, candle); err != nil {
				return closedCandle, fmt.Errorf("%v: %w", ErrUpdateCandlesMap, err)
			} else {
				cm[candle.Ticker] = &newCandle
				return closedCandle, ErrUpdateCandleMismatchedPeriod
			}
		}
	}
	return Candle{}, nil
}

func (cm CandleMap) FlushMap() []Candle {
	var candles []Candle
	for _, val := range cm {
		candles = append(candles, *val)
	}
	return candles
}

func NewCandleFromPrice(period CandlePeriod, p Price) (Candle, error) {
	candleTS, err := PeriodTS(period, p.TS)
	if err != nil {
		return Candle{}, fmt.Errorf("%e: %w", ErrCreateNewCandle, err)
	}
	return Candle{
		Ticker: p.Ticker,
		Period: period,
		Open:   p.Value,
		High:   p.Value,
		Low:    p.Value,
		Close:  p.Value,
		TS:     candleTS,
	}, nil
}

func NewCandleFromCandle(period CandlePeriod, c Candle) (Candle, error) {
	newCandle := c
	candleTS, err := PeriodTS(period, c.TS)
	if err != nil {
		return Candle{}, fmt.Errorf("%v: %w", ErrCreateNewCandle, err)
	}
	newCandle.TS = candleTS
	newCandle.Period = period
	return newCandle, nil
}

func (c *Candle) UpdateFromPrice(p Price) error {
	pTS, err := PeriodTS(c.Period, p.TS)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrUpdateCandle, err)
	}
	if pTS != c.TS {
		return ErrUpdateCandleMismatchedPeriod
	}
	if p.Value > c.High {
		c.High = p.Value
	}
	if p.Value < c.Low {
		c.Low = p.Value
	}
	c.Close = p.Value
	return nil
}

func (c *Candle) UpdateFromCandle(otherCandle Candle) error {
	candleTS, err := PeriodTS(c.Period, otherCandle.TS)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrUpdateCandle, err)
	}
	if c.TS != candleTS {
		return ErrUpdateCandleMismatchedPeriod
	}
	if otherCandle.High > c.High {
		c.High = otherCandle.High
	}
	if otherCandle.Low < c.Low {
		c.Low = otherCandle.Low
	}
	c.Close = otherCandle.Close
	return nil
}

func (c Candle) String() string {
	strFields := []string{
		c.Ticker, c.TS.Format(time.RFC3339),
		fmt.Sprintf("%f", c.Open), fmt.Sprintf("%f", c.High),
		fmt.Sprintf("%f", c.Low), fmt.Sprintf("%f", c.Close)}
	return strings.Join(strFields, ",")
}
