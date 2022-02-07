PROJECT := github.com/juju/packaging/v2

.PHONY: check-licence check-go check

check: check-licence check-go
	go test $(PROJECT)/...

check-licence:
	@(fgrep -rl "Licensed under the LGPLv3" .;\
		fgrep -rl "MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT" .;\
		find . -name "*.go") | sed -e 's,\./,,' | sort | uniq -u | \
		xargs -I {} echo FAIL: licence missed: {}

check-go:
	$(eval GOFMT := $(strip $(shell gofmt -l .| sed -e "s/^/ /g")))
	@(if [ "$(GOFMT)" != "" ]; then \
		echo go fmt is sad: $(GOFMT); \
		exit 1; \
	fi )
	@(go vet -all -composites=false -copylocks=false .)

static-analysis:
	golangci-lint run -c .github/golangci-lint.config.yaml

test:
	go test -v ./...
