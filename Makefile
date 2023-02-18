commit = $(shell cd /Users/jorrit/Documents/master-software-engineering/thesis/GoLib && git log -1 | head -1 | cut -d ' ' -f 2)

targets := anonymize query gateway

$(targets): %:
	cp Dockerfile $*_service/
	cd $*_service && go get -u "github.com/Jorrit05/GoLib@$(commit)"
	docker build --build-arg NAME='$*' -t $*_service $*_service/
	rm $*_service/Dockerfile

all: $(targets)

.PHONY: all $(targets)
