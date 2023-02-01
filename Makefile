commit = $(shell cd /Users/jorrit/Documents/master-software-engineering/thesis/GoLib && git log -1 | head -1 | cut -d ' ' -f 2)

randomize:
	cd randomize_service && go get -u "github.com/Jorrit05/GoLib@$(commit)"
	docker build -t randomize_service randomize_service/

query:
	cd query_service && go get -u "github.com/Jorrit05/GoLib@$(commit)"
	docker build -t query_service query_service/

gateway:
	cd gateway_service && go get -u "github.com/Jorrit05/GoLib@$(commit)"
	docker build -t gateway_service gateway_service/

all: randomize query gateway

.PHONY: all
.PHONY: randomize
.PHONY: query