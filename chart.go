package main

import (
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func writeChart(dailyKeyCounts []DailyKeyCount, filename string) {

	dates := make([]time.Time, 0)
	newKeys := make([]float64, 0)
	newKeys14Days := make([]float64, 0)

	for _, v := range dailyKeyCounts {
		d, _ := time.Parse(isoDateFormat, v.Date)
		dates = append(dates, d)
		// dates = append(dates, float64(i))

		newKeys = append(newKeys, float64(v.NewKeysCount))
		newKeys14Days = append(newKeys14Days, float64(v.KeysInLast14Days))
	}

	graph := chart.Chart{
		Width:      1000,
		Height:     600,
		Title:      "Ključi na #OstaniZdrav strežniku",
		TitleStyle: chart.StyleShow(),

		XAxis: chart.XAxis{
			Name:         "Dan",
			TickPosition: chart.TickPositionBetweenTicks,
			Style: chart.Style{
				Show:                true,
				TextRotationDegrees: 90,
			},
		},
		YAxis: chart.YAxis{
			Name:      "Aktivni ključi (14 dni)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxisSecondary: chart.YAxis{
			Name:      "Novi ključi (na dan)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},

		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Aktivni ključi (14 dni)",
				XValues: dates,
				YValues: newKeys14Days,
				YAxis:   chart.YAxisPrimary,
				Style: chart.Style{
					Show:        true,
					StrokeWidth: 2,
					StrokeColor: drawing.Color{R: 248, G: 198, B: 45, A: 255},
					FillColor:   drawing.Color{R: 248, G: 198, B: 45, A: 50},
				},
			},
			chart.TimeSeries{
				Name:    "Novi ključi (na dan)",
				XValues: dates,
				YValues: newKeys,
				YAxis:   chart.YAxisSecondary,
				Style: chart.Style{
					Show:        true,
					StrokeWidth: 2,
					StrokeColor: drawing.Color{R: 78, G: 126, B: 245, A: 255},
					FillColor:   drawing.Color{R: 78, G: 126, B: 245, A: 50},
				},
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, _ := os.Create(filename)
	defer f.Close()
	graph.Render(chart.PNG, f)
}
