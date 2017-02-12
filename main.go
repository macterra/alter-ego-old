package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"bufio"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"net/http"

	"goji.io"
	"goji.io/pat"
	
	"github.com/google/uuid"
)

var datapath string
var sessionID uuid.UUID
var sessionPath string

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

func dispatch(reqPath string) (respPath string, err error) {
	f, _ := os.Open(reqPath)
	r := bufio.NewReader(f)
	request, err := http.ReadRequest(r)
	
	//fmt.Println("dispatch", request, err)
	
	txPath := filepath.Dir(reqPath)
	respPath = filepath.Join(txPath, "response")
	
	response, _ := os.OpenFile(respPath, os.O_CREATE, 0644)
	request.Write(response)
	response.Close()
	
	fmt.Println("dispatch wrote:", respPath)
	
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	txID := uuid.New()
	txPath := filepath.Join(sessionPath, "txs", txID.String())
	os.MkdirAll(txPath, os.ModePerm)
	reqPath := filepath.Join(txPath, "request")
	f, _ := os.Create(reqPath)
	r.Write(f)
	f.Close()

	logPath := filepath.Join(sessionPath, "txs", "log")
	writeLog(logPath, txID)	
	
	respPath, err := dispatch(reqPath)
	
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	} else {
		response, _ := ioutil.ReadFile(respPath)
		w.Write(response)
	}
}

func writeLog(logPath string, id uuid.UUID) {
	log, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	fmt.Fprintf(log, "%v|%v\n", time.Now().UTC(), id)
	log.Close()
}

func main() {

	datapath = filepath.Join(os.Getenv("GOPATH"), "data") 
	err := os.MkdirAll(datapath, os.ModePerm)
	fmt.Println(datapath, err)
	
	configPath := filepath.Join(datapath, "config.json")
	config := Config{8000}
	config.Read(configPath)
	
	sessionID = uuid.New()
	sessionPath = filepath.Join(datapath, "sessions", sessionID.String())
	err = os.MkdirAll(sessionPath, os.ModePerm)
	fmt.Println("session:", sessionID, sessionPath, err)

	logPath := filepath.Join(datapath, "sessions", "log")
	writeLog(logPath, sessionID)	
	
	addr := "localhost:" + strconv.Itoa(config.Port)
	
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/*"), index)

	fmt.Println("Starting alter-ego on ", addr)
	http.ListenAndServe(addr, mux)
}
