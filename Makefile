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
