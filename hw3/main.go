package main

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"hw-async/domain"
	"hw-async/generator"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	ErrCreateFile        = errors.New("create file")
	ErrWriteCandleToFile = errors.New("writeCandleToFile")
)

const (
	candle1mFileName  = "candle1m.csv"
	candle2mFileName  = "candle2m.csv"
	candle10mFileName = "candle10m.csv"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func main() {
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	prices := pg.Prices(ctx)

	wg := sync.WaitGroup{}
	wg.Add(3)
	generate10mCandles(ctx, &wg,
		generate2mCandles(ctx, &wg,
			generate1mCandles(ctx, &wg, prices)))
	wg.Wait()
	logger.Info("program successfully finished.")
}

func generate1mCandles(ctx context.Context, wg *sync.WaitGroup, inPrices <-chan domain.Price) <-chan domain.Candle {
	outCandles := make(chan domain.Candle)
	candleMap := domain.CandleMap{}

	go func() {
		defer close(outCandles)
		defer wg.Done()

		file, err := createFile(domain.CandlePeriod1m)
		if err != nil {
			log.Error(err)
		}
		defer file.Close()

		for {
			select {
			case price := <-inPrices:
				log.Info(price)
				closedCandle, err := candleMap.UpdateFromPrice(price, domain.CandlePeriod1m)
				if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
					outCandles <- closedCandle
					if err := writeToFile(file, closedCandle); err != nil {
						log.Warningln(err)
					}
				}
			case <-ctx.Done():
				lastCandles := candleMap.FlushMap()
				for _, val := range lastCandles {
					if err := writeToFile(file, val); err != nil {
						log.Warningln(err)
					}
				}
				return
			}
		}
	}()

	return outCandles
}

func generate2mCandles(ctx context.Context, wg *sync.WaitGroup, inCandles <-chan domain.Candle) <-chan domain.Candle {
	outCandles := make(chan domain.Candle)
	candleMap := domain.CandleMap{}

	go func() {
		defer close(outCandles)
		defer wg.Done()

		file, err := createFile(domain.CandlePeriod2m)
		if err != nil {
			log.Error(err)
		}
		defer file.Close()

		for {
			select {
			case candle := <-inCandles:
				closedCandle, err := candleMap.UpdateFromCandle(candle, domain.CandlePeriod2m)
				if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
					outCandles <- closedCandle
					if err := writeToFile(file, closedCandle); err != nil {
						log.Warningln(err)
					}
				}
			case <-ctx.Done():
				candles := candleMap.FlushMap()
				for _, val := range candles {
					if err := writeToFile(file, val); err != nil {
						log.Warningln(err)
					}
				}
				return
			}
		}
	}()

	return outCandles
}

func generate10mCandles(ctx context.Context, wg *sync.WaitGroup, inCandles <-chan domain.Candle) {
	candleMap := domain.CandleMap{}

	go func() {
		defer wg.Done()

		file, err := createFile(domain.CandlePeriod10m)
		if err != nil {
			log.Error(err)
		}
		defer file.Close()

		for {
			select {
			case candle := <-inCandles:
				closedCandle, err := candleMap.UpdateFromCandle(candle, domain.CandlePeriod10m)
				if errors.Is(err, domain.ErrUpdateCandleMismatchedPeriod) {
					if err := writeToFile(file, closedCandle); err != nil {
						log.Warningln(err)
					}
				}
			case <-ctx.Done():
				candles := candleMap.FlushMap()
				for _, val := range candles {
					if err := writeToFile(file, val); err != nil {
						log.Warningln(err)
					}
				}
				return
			}
		}
	}()
}

func createFile(period domain.CandlePeriod) (*os.File, error) {
	var fileName string
	switch period {
	case domain.CandlePeriod1m:
		fileName = candle1mFileName
	case domain.CandlePeriod2m:
		fileName = candle2mFileName
	case domain.CandlePeriod10m:
		fileName = candle10mFileName
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
