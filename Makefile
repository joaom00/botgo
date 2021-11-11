.EXPORT_ALL_VARIABLE:

ENTRY_FILES=main.go

run-dev:
	@echo "=>> Running locally with air"
	@air run $(ENTRY_FILES)
