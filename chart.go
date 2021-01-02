package main

import (
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func writeChart(dailyKeyCounts []DailyKeyCount, filename string, country string) {

	dates := make([]time.Time, 0)
	newKeys := make([]float64, 0)
	newKeys14Days := make([]float64, 0)
	nonExpiredKeys := make([]float64, 0)

	for _, v := range dailyKeyCounts {
		d, _ := time.Parse(isoDateFormat, v.Date)
		dates = append(dates, d)
		// dates = append(dates, float64(i))

		newKeys = append(newKeys, float64(v.NewKeysCount))
		newKeys14Days = append(newKeys14Days, float64(v.NewKeysInLast14Days))
		nonExpiredKeys = append(nonExpiredKeys, float64(v.NonExpiredKeys))
	}

	activeKeysSeries := chart.TimeSeries{
		Name:    "Novi ključi (14 dni)",
		XValues: dates,
		YValues: newKeys14Days,
		YAxis:   chart.YAxisPrimary,
		Style: chart.Style{
			Show:        true,
			StrokeWidth: 2,
			StrokeColor: drawing.Color{R: 248, G: 198, B: 45, A: 255},
			FillColor:   drawing.Color{R: 252, G: 244, B: 213, A: 255},
			// FillColor:   drawing.Color{R: 248, G: 198, B: 45, A: 50},
		},
	}

	nonExpiredKeysSeries := chart.TimeSeries{
		Name:    "Nepretečeni ključi",
		XValues: dates,
		YValues: nonExpiredKeys,
		YAxis:   chart.YAxisPrimary,
		Style: chart.Style{
			Show:        true,
			StrokeWidth: 2,
			// StrokeDashArray: []float64{5, 2},
			StrokeColor: drawing.ColorRed, //.Color{R: 248, G: 150, B: 5, A: 255},
			FillColor:   drawing.Color{R: 252, G: 204, B: 183, A: 255},
			// FillColor: drawing.Color{R: 248, G: 198, B: 45, A: 50},
		},
	}

	newKeysSeries := chart.TimeSeries{
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
	}

	newKeysMovingAverageSeries := chart.SMASeries{
		Period:      7,
		InnerSeries: newKeysSeries,
		YAxis:       chart.YAxisSecondary,
		Name:        "Novi ključi (povprečje 7 dni)",
		Style: chart.Style{
			Show:            true,
			StrokeWidth:     2,
			StrokeDashArray: []float64{5, 2},
			StrokeColor:     drawing.Color{R: 78, G: 126, B: 245, A: 255},
		},
	}

	graph := chart.Chart{
		Width:      1000,
		Height:     600,
		Title:      fmt.Sprintf("%s ključi na #OstaniZdrav strežniku", country),
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
			Name:      "Nepretečeni/Novi ključi (14 dni)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%0.0f", vf)
				}
				return ""
			},
		},
		YAxisSecondary: chart.YAxis{
			Name:      "Novi ključi (na dan)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%0.0f", vf)
				}
				return ""
			},
		},

		Series: []chart.Series{
			activeKeysSeries,
			chart.LastValueAnnotation(activeKeysSeries),
			nonExpiredKeysSeries,
			chart.LastValueAnnotation(nonExpiredKeysSeries),
			newKeysSeries,
			newKeysMovingAverageSeries,
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, _ := os.Create(filename)
	defer f.Close()
	graph.Render(chart.PNG, f)
}
