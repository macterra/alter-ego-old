package main

import (
	"os"
	"fmt"
	"strconv"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"net/http"

	"goji.io"
	"goji.io/pat"
)

type Config struct {
	Port int
}

func (config *Config) Read(filepath string) {

	fmt.Println("read config:", config)

	data, err := ioutil.ReadFile(filepath)
	
	if err != nil {
		fmt.Println(err)
		
		data, err = json.Marshal(config)
		
		if err != nil {
			fmt.Println(err)
		} else {
			ioutil.WriteFile(filepath, data, 0644)
		}
	} else {
		err = json.Unmarshal(data, &config)
		
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("success", config)
		} 		
	}
	
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index")
}

func main() {

	datapath := filepath.Join(os.Getenv("GOPATH"), "data") 
	err := os.Mkdir(datapath, os.ModePerm)
    fmt.Println(datapath, err)
	
	configpath := filepath.Join(datapath, "config.json")
	config := Config{8000}
	config.Read(configpath)
	
	addr := "localhost:" + strconv.Itoa(config.Port)
	
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/"), index)

	fmt.Println("Starting alter-ego on ", addr)
	http.ListenAndServe(addr, mux)
}
