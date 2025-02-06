.PHONY: all build-claude build-backend up-db start-claude start-backend kill-ports-3000 kill-ports-9000 kill-ports

# Kill processes on ports 3000 and 9000
kill-ports-3000:
	@echo "Killing processes on port 3000..."
	lsof -ti:3000 | xargs kill -9 || echo "No process found on port 3000"

kill-ports-9000:
	@echo "Killing processes on port 9000..."
	lsof -ti:9000 | xargs kill -9 || echo "No process found on port 9000"

kill-ports:
	@echo "Killing processes on ports 3000 and 9000..."
	$(MAKE) kill-ports-3000
	$(MAKE) kill-ports-9000

# STEP1: lambda関数のビルド
build-claude:
	@echo "Building lambda functions in packages/claude..."
	cd packages/claude && rm -rf .aws-sam && sam build

build-backend:
	@echo "Building lambda functions in backend..."
	cd backend && rm -rf .aws-sam && sam build

# STEP2: DB立ち上げ
up-db:
	@echo "Starting the database using docker-compose in backend..."
	cd backend && docker compose up -d

# STEP3: lambda関数の立ち上げ
start-claude:
	@echo "Starting lambda functions in packages/claude... (in background)"
	cd packages/claude && sam local start-api --env-vars env.json --port 9000 &

start-backend:
	@echo "Starting lambda functions in backend... (in background)"
	cd backend && sam local start-api --env-vars env.json --port 3000 &

# Combined targets
build-all: build-claude build-backend

start-all: start-claude start-backend

# setupAll: 3000と9000をkillしてから立ち上げ
setupAll: kill-ports build-all up-db
	@echo "Waiting for the database to be ready..."
	sleep 3
	$(MAKE) start-all

# setupGolang: 3000をkillしてから立ち上げ
setupGolang: kill-ports-3000 build-backend up-db
	@echo "Waiting for the database to be ready..."
	sleep 3
	$(MAKE) start-backend

# setupClaude: 9000をkillしてから立ち上げ
setupClaude: kill-ports-9000 build-claude
	$(MAKE) start-claude
