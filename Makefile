.PHONY: proxy servers requests request-app request-app-test request-api

PROXY_URL := http://localhost:9090
CURL := curl --silent --show-error --fail-with-body

proxy:
	go run .

servers:
	python3 test/mock_servers.py

requests: request-app request-app-test request-api

request-app:
	$(CURL) --header 'Host: app.example.com' '$(PROXY_URL)/'
	@printf '\n'

request-app-test:
	$(CURL) --header 'Host: app.example.com' '$(PROXY_URL)/test'
	@printf '\n'

request-api:
	$(CURL) --header 'Host: api.example.com' '$(PROXY_URL)/'
	@printf '\n'
