# Temporary Exposure Key Export scraper

[![Periodic update status](https://github.com/sledilnik/cwa-scrape/workflows/Periodic%20update/badge.svg)](https://github.com/sledilnik/cwa-scrape/actions)

Periodically scrapes the daily data from [Slovenian contact tracing app #OstaniZdrav](https://www.gov.si/ostanizdrav) and publishes it into this git repository into [data/SI](data/SI) directory.

After scraping all the data, it aggregates some statistics into [keycount.csv](data/SI/keycount.csv) and [.json](data/SI/keycount.json), which is then plotted onto the [keycount.png](data/SI/keycount.png) chart:

![Chart of new and active keys on the Corona Warn App server](data/SI/keycount.png)

## Requirements

GNU Make, [Go](https://golang.org) 1.15 or later, [jq](https://stedolan.github.io/jq/) 1.6 or later.

```bash
$ sudo apt install golang jq
```

## Running

It is ran hourly on GitHub servers using github [action](.github/workflows/make.yml). To run it locally just do:

```bash
$ make
```

## Tools & Libraries used

* [Export analyzer tool](https://github.com/google/exposure-notifications-server/tree/main/tools/export-analyzer) from [Google exposure notification server](https://github.com/google/exposure-notifications-server)
* [github.com/wcharczuk/go-chart](https://github.com/wcharczuk/go-chart) charting library for Go
* [github.com/gocarina/gocsv](https://github.com/gocarina/gocsv) CSV library for Go

## Acknowledgements

This project was inspired by a devoted team of [sledilnik.org](https://sledilnik.org), a covid-19 tracker for Slovenia. Other related projects can be found at [:octocat: GitHub](https://github.com/sledilnik).
