commit = $(shell cd /Users/jorrit/Documents/master-software-engineering/thesis/GoLib && git log -1 | head -1 | cut -d ' ' -f 2)

anonymize:
	cd anonymize_service && go get -u "github.com/Jorrit05/GoLib@$(commit)"
	docker build -t anonymize_service anonymize_service/

query:
	cd query_service && go get -u "github.com/Jorrit05/GoLib@$(commit)"
	docker build -t query_service query_service/

gateway:
	cd gateway_service && go get -u "github.com/Jorrit05/GoLib@$(commit)"
	docker build -t gateway_service gateway_service/

all: anonymize query gateway

.PHONY: all
.PHONY: anonymize
.PHONY: query