TEST_DIRS := $(shell find test/domain -type d)

test:
	@echo "Running tests in $(TEST_DIRS)..."
	@go test -v $(TEST_DIRS)
