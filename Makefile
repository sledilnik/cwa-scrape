BASEURL = https://svc90.cwa.gov.si/version/v1/
COUNTRY = SI

download:
	mkdir -p data/$(COUNTRY) || true

	for DAY in $$(curl -s $(BASEURL)diagnosis-keys/country/$(COUNTRY)/date | jq -r '.[]');\
	do \
		echo -n 'Downloading $$DAY:	';\
		wget $(BASEURL)diagnosis-keys/country/$(COUNTRY)/date/$$DAY -O data/$(COUNTRY)/$$DAY.zip;\
	done

	wget $(BASEURL)configuration/country/$(COUNTRY)/app_config -O data/$(COUNTRY)/app_config.zip

analyze:
	# go install github.com/google/exposure-notifications-server/tools/export-analyzer
	for file in $$(ls data/$(COUNTRY)/????-??-??.zip);\
	do \
		BASENAME=$$(basename $$file .zip);\
		echo "Analyzing $$file:	";\
		export-analyzer -q --file="$$file" 2>&1 |tail -n +2 >"data/$(COUNTRY)/$$BASENAME.json" 2>&1;\
		export-analyzer -json=false --file="$$file" >"data/$(COUNTRY)/$$BASENAME.log" 2>&1;\
	done
