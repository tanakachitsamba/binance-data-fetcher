package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
)

type Prices struct {
	price []*binance.Kline
	sym   string
}

func main() {
	// Create a new Binance client
	client := binance.NewClient("API_KEY", "SECRET_KEY")

	syms := []string{"BTCUSDT", "ETHUSDT", "XRPUSDT", "DOGEUSDT", "AXSUSDT", "BNBUSDT", "LTCUSDT", "SOLUSDT", "MATICUSDT", "TKOUSDT", "LINKUSDT", "GALAUSDT"}

	var prices []Prices
	for _, i := range syms {
		price, err := client.NewKlinesService().Symbol(i).Limit(1000).Interval("5m").Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		var p Prices = Prices{price: price, sym: i}

		prices = append(prices, p)
	}

	var priceSet []PriceData

	for _, price := range prices {
		for _, i := range price.price {
			var convPrice PriceData

			convPrice.High = i.High
			convPrice.Close = i.Close

			convPrice.Low = i.Low
			convPrice.Volume = i.Volume
			convPrice.Open = i.Open
			convPrice.QuoteAssetVolume = i.QuoteAssetVolume
			convPrice.TakerBuyBaseAssetVolume = i.TakerBuyBaseAssetVolume
			convPrice.TakerBuyQuoteAssetVolume = i.TakerBuyQuoteAssetVolume

			convPrice.CloseTime = strconv.Itoa(int(i.CloseTime))
			convPrice.OpenTime = strconv.Itoa(int(i.OpenTime))
			convPrice.TradeNum = strconv.Itoa(int(i.TradeNum))

			convPrice.sym = price.sym

			priceSet = append(priceSet, convPrice)

		}
	}

	percentageGain := func(buyPrice, sellPrice decimal.Decimal) decimal.Decimal {
		return sellPrice.Sub(buyPrice).Div(buyPrice).Mul(decimal.NewFromInt32(100))
	}

	var (
		h, c  decimal.Decimal
		index int
	)
	for idx, i := range priceSet {
		close, err := decimal.NewFromString(i.Close)
		if err != nil {
			fmt.Println(err)
			return
		}

		high, err := decimal.NewFromString(i.High)
		if err != nil {
			fmt.Println(err)
			return
		}

		if index == 0 && idx == 0 {
			h = high
			c = close
		}

		priceSet[idx].prof = "0"

		if idx > 0 {
			if percentageGain(h, high).GreaterThan(decimal.NewFromFloat32(0.5)) || percentageGain(c, close).GreaterThan(decimal.NewFromFloat32(0.5)) {
				priceSet[idx-1].prof = "1"
			} else {
				priceSet[idx-1].prof = "0"
			}
		}

		if index != idx {
			h = high
			c = high
			index = idx

		}
	}

	var groups [][]string
	for idx, _ := range priceSet {
		if idx >= 4 && idx < 994 {
			priv := priceSet[idx : idx+5]

			var sets []string
			for _, i := range priv {
				sets = append(sets, i.High)
				sets = append(sets, i.Close)
				sets = append(sets, i.Open)
				sets = append(sets, i.Low)
				sets = append(sets, i.Volume)
				sets = append(sets, i.QuoteAssetVolume)
				sets = append(sets, i.TakerBuyBaseAssetVolume)
				sets = append(sets, i.TakerBuyQuoteAssetVolume)
				sets = append(sets, i.TradeNum)
				sets = append(sets, i.OpenTime)
				sets = append(sets, i.CloseTime)
				sets = append(sets, i.prof)
				sets = append(sets, i.sym)
			}
			groups = append(groups, sets)
		}

	}

	// Open a file for writing
	f, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a new CSV writer
	w := csv.NewWriter(f)

	labels := []string{"High", "Close", "Open", "Low", "Volume", "QuoteAssetVolume", "TakerBuyBaseAssetVolume", "TakerBuyQuoteAssetVolume", "TradeNum", "OpenTime", "CloseTime", "prof", "sym"}
	var processedLabels []string
	for i := 0; i < 5; i++ {
		for _, item := range labels {
			processedLabels = append(processedLabels, item+strconv.Itoa(i))
		}
	}

	// Write the row labels as the first record in the CSV file
	err = w.Write(processedLabels)
	if err != nil {
		panic(err)
	}

	// Loop through the rows of data
	for _, row := range groups {
		// Loop through the values in the row
		err := w.Write(row)
		if err != nil {
			panic(err)
		}
	}

	// Flush the writer to ensure that all data is written to the output stream
	w.Flush()
}

type PriceData struct {
	High, Close, Open, Low, Volume, QuoteAssetVolume, TakerBuyBaseAssetVolume, TakerBuyQuoteAssetVolume, TradeNum, OpenTime, CloseTime, prof, sym string
}
