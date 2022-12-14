.PHONY: generate

generate:
	@go generate ./...
	@echo "[OK] WASM added!"

security:
	@gosec -exclude-dir=.cache -exclude-generated ./...
	@echo "[OK] Go security check was completed!"
