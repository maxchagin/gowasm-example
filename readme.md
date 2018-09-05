# Golang 1.11 webassembly experiment

## Search users by Github on gowasm

### Build
`go get github.com/maxchagin/gowasm-example`    
`GOARCH=wasm GOOS=js go build -o web/test.wasm main.go`

### Run server
`go run server.go`

### Browser
Open page http://localhost:8080/web/wasm_exec.html
