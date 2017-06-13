.PHONY: test

GO_TEST=go test -v -cover -covermode=count -coverprofile=single.coverprofile

test:
	@echo Running tests
	@$(eval PKGS := $(shell go list ./... | grep -v /vendor/))
	@echo "mode: count" >  merged.coverpofile
	@$(foreach PKG, $(PKGS), $(GO_TEST) $(PKG) || exit 1 ; cat single.coverprofile | grep -v "mode:" >> merged.coverpofile ;)
	@go tool cover --html ./merged.coverpofile -o coverage.html
	@rm *.coverpofile
