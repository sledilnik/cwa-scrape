BASEURL = https://svc90.cwa.gov.si/version/v1/
COUNTRY = SI

all: download analyze

download:
	./download.sh $(COUNTRY) $(BASEURL)

analyze:
	go install github.com/google/exposure-notifications-server/tools/export-analyzer
	for file in $$(ls data/$(COUNTRY)/????-??-??.zip);\
	do \
		BASENAME=$$(basename $$file .zip);\
		echo "Analyzing $$file:	";\
		export-analyzer -q -sig=false --file="$$file" 2>&1 |tail -n +2 >"data/$(COUNTRY)/$$BASENAME.json" 2>&1;\
		export-analyzer -json=false --file="$$file" >"data/$(COUNTRY)/$$BASENAME.log" 2>&1 || true;\
	done

	go run export-aggregate.go chart.go --path=data/$(COUNTRY)
