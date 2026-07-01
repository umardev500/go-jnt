check_awb:
	go run cmd/cli/check_zero_awb/main.go

api:
	go run cmd/api/scan/main.go

report:
	go run cmd/cli/og_reporting/main.go

sj:
	go run cmd/cli/sj/main.go $(filter-out $@,$(MAKECMDGOALS))

manifest:
	go run cmd/cli/manifest/main.go $(filter-out $@,$(MAKECMDGOALS))

info:
	go run cmd/cli/main.go $(filter-out $@,$(MAKECMDGOALS))

cf:
	cloudflared tunnel --url http://localhost:8080

get_sj:
	go run cmd/cli/get_sj/main.go

%:
	@: