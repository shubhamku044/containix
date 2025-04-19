APP_NAME=containix

.PHONY: dev run build clean

dev:
	@echo "🚀 Starting ${APP_NAME} in dev mode with air..."
	air

run:
	@echo " Running ${APP_NAME}..."
	go run ./main.go

