.PHONY: dev/server
dev/server: 
	go run github.com/goware/rerun/cmd/rerun@latest \
		-watch cmd \
		-watch internal \
		-ignore app \
		-run "go run ./cmd/api"

.PHONY: dev/app
dev/app: 
	cd ./app && bash -c "pnpm dev"

.PHONY: dev
dev:
	make -j2 dev/server dev/app
