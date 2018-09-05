package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"
	"time"

	"github.com/cbroglie/mustache"
)

const (
	ApiGitHub = "https://api.github.com"
)

type Search struct {
	Response Response
	Result   Result
}

type Response struct {
	Status string
}

// Result struct GitHub API https://api.github.com/users/maxchagin
type Result struct {
	Login            string `json:"login"`
	ID               int    `json:"id"`
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	AvatarURL        string `json:"avatar_url"`
	Name             string `json:"name"`
	PublicRepos      int    `json:"public_repos"`
	PublicGists      int    `json:"public_gists"`
	Followers        int    `json:"followers"`
}

type Element struct {
	tag    string
	params map[string]string
}

type Box struct {
	el js.Value
}

type App struct {
	inputBox  js.Value
	resultBox js.Value
	userTMPL  string
	errorTMPL string
	search    chan Search
}

func main() {
	forever := make(chan bool)
	// Get the main app block element
	box := Box{
		el: js.Global().Get("document").Call("getElementById", "box"),
	}

	a := App{
		inputBox:  box.createInputBox(),
		resultBox: box.createResultBox(),
		userTMPL:  getTMPL("user.mustache"),
		errorTMPL: getTMPL("error.mustache"),
		search:    make(chan Search),
	}
	// User input handler
	go a.userHandler()
	// Output User Information
	go a.listResults()

	<-forever
}

func getTMPL(name string) string {
	resp, err := http.Get("tmpl/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (el *Element) createEl() js.Value {
	// Create an item, similarly:
	// var e = document.createElement('input');
	e := js.Global().Get("document").Call("createElement", el.tag)
	// Set attributes and their values
	for attr, value := range el.params {
		e.Set(attr, value)
	}
	return e
}

// Create an input form
func (b *Box) createInputBox() js.Value {
	el := Element{
		tag: "input",
		params: map[string]string{
			"placeholder": "GitHub username",
		},
	}
	input := el.createEl()
	// Output to the page, adding to the div with id = app
	b.el.Call("appendChild", input)
	return input
}

// Create a block to insert search results
func (b *Box) createResultBox() js.Value {
	el := Element{
		tag: "div",
		params: map[string]string{
			"id": "search_result",
		},
	}
	div := el.createEl()
	b.el.Call("appendChild", div)
	return div
}

func (a *App) userHandler() {
	spammyChan := make(chan string, 10)
	go debounce(1000*time.Millisecond, spammyChan, func(arg string) {
		// Get Data with github api
		go a.getUserCard(arg)
	})
	a.inputBox.Call("addEventListener", "keyup", js.NewCallback(func(args []js.Value) {
		// Placeholder "Loading..."
		a.loadingResults()

		e := args[0]
		// Get the value of an element
		user := e.Get("target").Get("value").String()
		// Clear the results block
		if user == "" {
			a.clearResults()
		}
		spammyChan <- user
		println(user)
	}))
}

func (a *App) getUserCard(user string) {
	resp, err := http.Get(ApiGitHub + "/users/" + user)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var search Search
	json.Unmarshal(b, &search.Result)

	search.Response.Status = resp.Status
	a.search <- search
}

func (a *App) listResults() {
	var tmpl string
	for {
		search := <-a.search
		switch search.Result.ID {
		case 0:
			// TMPL for 404 page
			tmpl = a.errorTMPL
		default:
			tmpl = a.userTMPL
		}
		data, _ := mustache.Render(tmpl, search)
		// Output the result to a page
		a.resultBox.Set("innerHTML", data)
	}
}

func (a *App) loadingResults() {
	a.resultBox.Set("innerHTML", "<b>Loading...</b>")
}

func (a *App) clearResults() {
	a.resultBox.Set("innerHTML", "")
}

// Source: https://drailing.net/2018/01/debounce-function-for-golang/
func debounce(interval time.Duration, input chan string, cb func(arg string)) {
	var item string
	timer := time.NewTimer(interval)
	for {
		select {
		case item = <-input:
			timer.Reset(interval)
		case <-timer.C:
			if item != "" {
				cb(item)
			}
		}
	}
}
