package main

import (
	ac "autocomplete"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "Roses are red\nYou're my treasure\nI love Patricia\nBe with me forever") // send data to client side

	//http.ServeFile(w, r, "test.html")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

//autocomplete service
func (autocomplete_service *autocomplete_wrapper) handler(w http.ResponseWriter, r *http.Request) {
	s := autocomplete_service
	defer timeTrack(time.Now(), "autocomplete handler")
	//parse the url and return stuff
	r.ParseForm()
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	s.channel_base <- r.URL.Path[1:]
	words := <-s.channel_result
	//?MUST FIND A WAY TO FORMAT THE OUTPUTTUTU
	for _, word := range words {
		//fmt.Println(word)
		fmt.Fprintln(w, word)

	}

}

type autocomplete_wrapper struct {
	channel_base   chan string
	channel_result chan []string
	service        *ac.Trie
}

func setup_autocomplete_service() {
	//setup wordbank
	s := new(ac.Trie)
	s.Init()
	ac.LoadDictionary(s, "./src/autocomplete/words.txt")

	//setup handler
	ac_service := new(autocomplete_wrapper)
	ac_service.service = s
	ac_service.channel_base = make(chan string)
	ac_service.channel_result = make(chan []string)
	
	//run the service
	go func() {
		for {
			ac_service.channel_result <- ac_service.service.Autocomplete(<-ac_service.channel_base)
		}
	}()

	//set the route
	http.HandleFunc("/", ac_service.handler)
}

func main() {
	//spin up the autocomplete service
	setup_autocomplete_service()

	//http.HandleFunc("/favicon.ico", ac.ac_handler)
	http.HandleFunc("/favicon.ico", faviconHandler)

	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
