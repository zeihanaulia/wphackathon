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
	
