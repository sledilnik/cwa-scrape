# Temporary Exposure Key Export scraper

[![Periodic update status](https://github.com/stefanb/cwa-scrape/workflows/Periodic%20update/badge.svg)](https://github.com/stefanb/cwa-scrape/actions)

Periodically scrapes the daily data from Slovenian contact tracing app #OstaniZdrav and publishes it into this git repository into `data/SI` directory.

It aggregates some statistics into [keycount.csv](data/SI/keycount.csv) and [.json](data/SI/keycount.json), which can be used to produce charts like:

![Chart of new and active keys on the Corona Warn App server](data/SI/keycount.png)
