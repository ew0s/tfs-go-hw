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
	ErrCreateFile              = errors.New("create file")
	ErrWriteCandleToFile       = errors.New("writeCandleToFile")
	ErrGenerateCandleFromPrice = errors.New("generate candle from price")
	ErrUpdateCandle            = errors.New("update candle")
	ErrSave                    = errors.New("save")
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
	var errs []<-chan error

	candlesFromPrices, err := generateCandleFromPrice(domain.CandlePeriod1m, prices)
	errs = append(errs, err)

	candles1m, errCh, saveErrCh := generateCandleFromCandle(domain.CandlePeriod1m, candlesFromPrices)
	errs = append(errs, errCh, saveErrCh)

	candles2m, errCh, saveErrCh := generateCandleFromCandle(domain.CandlePeriod2m, candles1m)
	errs = append(errs, errCh, saveErrCh)

	errCh, saveErrCh = generateLastCandleFromCandle(domain.CandlePeriod10m, candles2m)
	errs = append(errs, errCh, saveErrCh)

	waitForPipeline(errs...)
}

func waitForPipeline(errs ...<-chan error) {
	errCh := mergeErrors(errs...)
	for err := range errCh {
		if err != nil {
			log.Warn(err)
		}
	}
}

func mergeErrors(errs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(errs))

	output := func(c <-chan error) {
		defer wg.Done()
		for val := range c {
			out <- val
		}
	}

	wg.Add(len(errs))
	for _, errCh := range errs {
		go output(errCh)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func generateCandleFromPrice(period domain.CandlePeriod, inPrices <-chan domain.Price) (<-chan domain.Candle, <-chan error) {
	outCandles := make(chan domain.Candle)
	errCh := make(chan error)

	go func() {
		defer close(outCandles)
		defer close(errCh)

		for price := range inPrices {
			log.Info(price)
			newCandle, err := domain.NewCandleFromPrice(price, period)
			if err != nil {
				errCh <- fmt.Errorf("%s: %w", ErrGenerateCandleFromPrice, err)
				continue
			}
			outCandles <- newCandle
		}
	}()

	return outCandles, errCh
}

func generateCandleFromCandle(period domain.CandlePeriod, inCandles <-chan domain.Candle) (<-chan domain.Candle, <-chan error, <-chan error) {
	outCandles := make(chan domain.Candle)
	errCh := make(chan error)
	candleMap := domain.CandleMap{}

	saveChan, saveErrCh := save(period, outCandles)

	go func() {
		defer close(outCandles)
		defer close(errCh)

		for candle := range inCandles {
			updateCandle(candleMap, candle, period, outCandles, errCh)
		}
		flushLastCandles(candleMap, outCandles)
	}()

	return saveChan, errCh, saveErrCh
}

func generateLastCandleFromCandle(period domain.CandlePeriod, inCandles <-chan domain.Candle) (<-chan error, <-chan error) {
	candles, errCh, saveErrCh := generateCandleFromCandle(period, inCandles)
	go func() {
		for range candles {
		}
	}()
	return errCh, saveErrCh
}

func save(period domain.CandlePeriod, c <-chan domain.Candle) (<-chan domain.Candle, <-chan error) {
	writerChan := make(chan domain.Candle)
	errCh := make(chan error)

	go func() {
		defer close(writerChan)
		defer close(errCh)

		file, err := createFile(period)
		if err != nil {
			errCh <- fmt.Errorf("%s: %w", ErrSave, err)
		}
		defer file.Close()

		for candle := range c {
			writerChan <- candle
			if err := writeToFile(file, candle); err != nil {
				errCh <- fmt.Errorf("%s: %w", ErrSave, err)
			}
		}
	}()

	return writerChan, errCh
}

func flushLastCandles(cm domain.CandleMap, outCandles chan<- domain.Candle) {
	candles := cm.FlushMap()
	for _, candle := range candles {
		outCandles <- candle
	}
}

func updateCandle(cm domain.CandleMap, c domain.Candle, p domain.CandlePeriod, outCh chan<- domain.Candle, errCh chan<- error) {
	closedCandle, err := cm.Update(c, p)
	if err != nil {
		if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
			outCh <- closedCandle
		} else {
			errCh <- fmt.Errorf("%v: %w", ErrUpdateCandle, err)
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
