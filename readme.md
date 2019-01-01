# Go wasm example

## Search users by Github on gowasm
Example of manipulation with DOM, working with the template and generating queries    


![Go wasm example](https://hsto.org/getpro/habr/post_images/9c2/020/f20/9c2020f205bbbbaeae6734360cf7bfe9.gif)

### Clone project
`cd work_dir`   
`git clone https://github.com/maxchagin/gowasm-example ./gowasm-example`   
`cd gowasm-example`

### Build app
`GOARCH=wasm GOOS=js go build -o web/test.wasm main.go`

### Run server
`go run server.go`

### Browser
Open page http://localhost:8080/web/wasm_exec.html

### Demo
http://wasm.lovefrontend.ru/web/wasm_exec.html
