PROJECTNAME := $(shell basename "$(PWD)")
STIME := $(shell date +%s)

.PHONY: build
build:
	@echo "  >  Building Program..."
	go build -ldflags="-s -w" -o bin/${PROJECTNAME} main.go; 
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"
	
.PHONY: scrape_instagram
scrape_instagram: build
	@echo "  >  Running Program..."
	./bin/${PROJECTNAME} scrape
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"

.PHONY: rest_api
rest_api: build
	@echo "  >  Running Program..."
	./bin/${PROJECTNAME} rest
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"
	
.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o internal/engine/ports/openapi_types.gen.go -package ports api/scraper.yml
	oapi-codegen -generate chi-server -o internal/engine/ports/openapi_api.gen.go -package ports api/scraper.yml
