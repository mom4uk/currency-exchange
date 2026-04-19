.PHONY: tests

tests:
	gotestsum --format=short-verbose ./tests/...