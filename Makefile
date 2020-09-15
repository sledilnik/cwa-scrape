BASEURL = https://svc90.cwa.gov.si/version/v1/
COUNTRY = SI
TODAY=$$(date --date="today" --iso-8601 2>&- || gdate --date="today" --iso-8601)
YESTERDAY=$$(date --date="yesterday" --iso-8601 2>&- || gdate --date="yesterday" --iso-8601)

all: download analyze

download:
	mkdir -p data/$(COUNTRY) || true
	wget $(BASEURL)diagnosis-keys/country/$(COUNTRY)/date/$(TODAY) -O data/$(COUNTRY)/$(TODAY).zip
	wget $(BASEURL)diagnosis-keys/country/$(COUNTRY)/date/$(YESTERDAY) -O data/$(COUNTRY)/$(YESTERDAY).zip
	wget $(BASEURL)configuration/country/$(COUNTRY)/app_config -O data/$(COUNTRY)/app_config.zip

analyze:
	# go install github.com/google/exposure-notifications-server/tools/export-analyzer
	for file in $$(ls data/$(COUNTRY)/????-??-??.zip);\
	do \
		BASENAME=$$(basename $$file .zip);\
		echo "Analyzing $$file:	";\
		export-analyzer -q --file="$$file" 2>&1 |tail -n +2 >"data/$(COUNTRY)/$$BASENAME.json" 2>&1;\
		export-analyzer -json=false --file="$$file" >"data/$(COUNTRY)/$$BASENAME.log" 2>&1 || true;\
	done

	go run export-aggregate.go chart.go
