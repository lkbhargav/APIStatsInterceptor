run: 
	. .env
	go run main.go

production:
	env GOOS=linux GOARCH=amd64 go build -o dist/APIStatsInterceptor_linux main.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/APIStatsInterceptor_mac main.go
	env GOOS=linux GOARCH=arm go build -o dist/APIStatsInterceptor_linux_arm main.go
	env GOOS=windows GOARCH=amd64 go build -o dist/APIStatsInterceptor_win.exe main.go
	zip -r dist.zip dist