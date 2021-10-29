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
	startPipeline(ctx, pg.Prices(ctx))
	logger.Info("program successfully finished.")
}

func startPipeline(ctx context.Context, prices <-chan domain.Price) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	generate10mCandles(&wg,
		generate2mCandles(&wg,
			generate1mCandles(ctx, &wg, prices)))
	wg.Wait()
}

func generate1mCandles(ctx context.Context, wg *sync.WaitGroup, inPrices <-chan domain.Price) <-chan domain.Candle {
	outCandles := make(chan domain.Candle)
	writerChan := make(chan domain.Candle)
	candleMap := domain.CandleMap{}

	go func() {
		defer wg.Done()
		defer close(outCandles)
		defer close(writerChan)

		go fileWriter(writerChan, domain.CandlePeriod1m)

		for {
			select {
			case price := <-inPrices:
				log.Info(price)
				newCandle, err := domain.NewCandleFromPrice(price, domain.CandlePeriod1m)
				if err != nil {
					log.Warn(err)
					continue
				}
				closedCandle, err := candleMap.Update(newCandle, domain.CandlePeriod1m)
				if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
					outCandles <- closedCandle
					writerChan <- closedCandle
				}
			case <-ctx.Done():
				candles := candleMap.FlushMap()
				for _, val := range candles {
					outCandles <- val
					writerChan <- val
				}
				return
			}
		}
	}()

	return outCandles
}

func generate2mCandles(wg *sync.WaitGroup, inCandles <-chan domain.Candle) <-chan domain.Candle {
	writerChan := make(chan domain.Candle)
	outCandles := make(chan domain.Candle)
	candleMap := domain.CandleMap{}

	go func() {
		defer wg.Done()
		defer close(outCandles)
		defer close(writerChan)

		go fileWriter(writerChan, domain.CandlePeriod2m)

		for candle := range inCandles {
			closedCandle, err := candleMap.Update(candle, domain.CandlePeriod2m)
			if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
				outCandles <- closedCandle
				writerChan <- closedCandle
			}
		}

		candles := candleMap.FlushMap()
		for _, val := range candles {
			outCandles <- val
			writerChan <- val
		}
	}()

	return outCandles
}

func generate10mCandles(wg *sync.WaitGroup, inCandles <-chan domain.Candle) {
	writerChan := make(chan domain.Candle)
	candleMap := domain.CandleMap{}

	go func() {
		defer wg.Done()
		defer close(writerChan)

		go fileWriter(writerChan, domain.CandlePeriod10m)

		for candle := range inCandles {
			closedCandle, err := candleMap.Update(candle, domain.CandlePeriod10m)
			if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
				writerChan <- closedCandle
			}
		}

		candles := candleMap.FlushMap()
		for _, val := range candles {
			writerChan <- val
		}
	}()
}

func fileWriter(candles <-chan domain.Candle, period domain.CandlePeriod) {
	file, err := createFile(period)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()
	for candle := range candles {
		if err := writeToFile(file, candle); err != nil {
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
