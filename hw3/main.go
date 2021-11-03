package main

import (
	"context"
	"errors"
	"fmt"
	"hw-async/domain"
	"hw-async/generator"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrCreateFile        = errors.New("create file")
	ErrWriteCandleToFile = errors.New("writeCandleToFile")
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		osCall := <-signals
		fmt.Println()
		log.Infof("system call: %+v", osCall)
		cancel()
	}()

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	startPipeline(pg.Prices(ctx))
	logger.Info("program successfully finished.")
}

func startPipeline(prices <-chan domain.Price) {
	wg := sync.WaitGroup{}
	wg.Add(4)

	candlesFromPrices := generateCandleFromPrice(&wg, domain.CandlePeriod1m, prices)
	candles1m := generateCandleFromCandle(&wg, domain.CandlePeriod1m, candlesFromPrices)
	candles2m := generateCandleFromCandle(&wg, domain.CandlePeriod2m, candles1m)
	generateLastCandleFromCandle(&wg, domain.CandlePeriod10m, candles2m)

	wg.Wait()
}

func generateCandleFromPrice(wg *sync.WaitGroup, period domain.CandlePeriod, inPrices <-chan domain.Price) <-chan domain.Candle {
	outCandles := make(chan domain.Candle)

	go func() {
		defer wg.Done()
		defer close(outCandles)

		for price := range inPrices {
			log.Info(price)
			newCandle, err := domain.NewCandleFromPrice(price, period)
			if err != nil {
				log.Warn(err)
				continue
			}
			outCandles <- newCandle
		}
	}()

	return outCandles
}

func generateCandleFromCandle(wg *sync.WaitGroup, period domain.CandlePeriod, inCandles <-chan domain.Candle) <-chan domain.Candle {
	outCandles := make(chan domain.Candle)
	candleMap := domain.CandleMap{}

	wg.Add(1)
	saveChan := save(wg, period, outCandles)

	go func() {
		defer wg.Done()
		defer close(outCandles)

		for candle := range inCandles {
			updateCandle(candleMap, candle, period, outCandles)
		}
		flushLastCandles(candleMap, outCandles)
	}()

	return saveChan
}

func generateLastCandleFromCandle(wg *sync.WaitGroup, period domain.CandlePeriod, inCandles <-chan domain.Candle) {
	candles := generateCandleFromCandle(wg, period, inCandles)
	for range candles {
	}
}

func save(wg *sync.WaitGroup, period domain.CandlePeriod, c <-chan domain.Candle) <-chan domain.Candle {
	writerChan := make(chan domain.Candle)

	go func() {
		defer wg.Done()
		defer close(writerChan)

		file, err := createFile(period)
		if err != nil {
			log.Warn(err)
		}
		defer file.Close()

		for candle := range c {
			writerChan <- candle
			if err := writeToFile(file, candle); err != nil {
				log.Warn(err)
			}
		}
	}()

	return writerChan
}

func flushLastCandles(cm domain.CandleMap, outCandles chan<- domain.Candle) {
	candles := cm.FlushMap()
	for _, candle := range candles {
		outCandles <- candle
	}
}

func updateCandle(cm domain.CandleMap, c domain.Candle, p domain.CandlePeriod, outCh chan<- domain.Candle) {
	closedCandle, err := cm.Update(c, p)
	if err != nil {
		if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
			outCh <- closedCandle
		} else {
			log.Warn(err)
		}
	}
}

func createFile(period domain.CandlePeriod) (*os.File, error) {
	var fileName string
	switch period {
	case domain.CandlePeriod1m, domain.CandlePeriod2m, domain.CandlePeriod10m:
		fileName = fmt.Sprintf("candles_%s.csv", period)
	default:
		return nil, fmt.Errorf("%v: %w", ErrCreateFile, domain.ErrUnknownPeriod)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrCreateFile, err)
	}
	return file, nil
}

func writeToFile(file *os.File, candle domain.Candle) error {
	_, err := file.WriteString(candle.String() + "\n")
	if err != nil {
		return fmt.Errorf("%v: %w", ErrWriteCandleToFile, err)
	}
	return nil
}
