BASEURL = https://svc90.cwa.gov.si/version/v1/

all: download analyze appconfig

download:
	./download.sh $(BASEURL)

analyze:
	go install github.com/google/exposure-notifications-server/tools/export-analyzer


	for COUNTRY in $(shell find data -type d -maxdepth 1 -regex 'data/[A-Z]*' -execdir echo {} ';') ; do \
		echo "Analyzing data for country: $$COUNTRY" ;\
		for file in $$(ls data/$$COUNTRY/????-??-??.zip); do \
			BASENAME=$$(basename $$file .zip);\
			echo "Analyzing $$file:	";\
			export-analyzer -sig=false -tek-age=336000h --file="$$file" 2>&1 |tail -n +3 | head -n -2 | sed 's/export: {/{/' >"data/$$COUNTRY/$$BASENAME.json" 2>&1;\
			export-analyzer -json=false -tek-age=336000h --file="$$file" >"data/$$COUNTRY/$$BASENAME.log" 2>&1 || true;\
		done;\
		go run export-aggregate.go chart.go --path=data/$$COUNTRY --country=$$COUNTRY ;\
	done;

appconfig:
	ls -la *.sh
	./check-app-config.sh $(BASEURL)
