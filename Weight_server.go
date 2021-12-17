package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type WeightRecord struct {
	Weight float64 `json:"weight"`
	Date   string  `json:"date"`
}

var (
	fileMutex sync.RWMutex
)

type Handler func(*http.Request) error

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
}

func AddWeight(r *http.Request) error {
	vars := r.URL.Query()
	weight, _ := strconv.ParseFloat(vars["weight"][0], 64)
	date := vars["date"][0]
	fmt.Println(weight, date)
	return addWeight(weight, date)
}

func addWeight(weight float64, date string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	wls := []WeightRecord{}
	data, err := ioutil.ReadFile("weight.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &wls)
	if err != nil {
		return err
	}
	wls = append(wls, WeightRecord{
		Weight: weight,
		Date:   date,
	})
	data, err = json.Marshal(wls)
	f, err := os.OpenFile("weight.json", os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func register(pattern string, handler Handler) {
	http.Handle(pattern, handler)
}

func main() {
	register("/add", AddWeight)
	err := http.ListenAndServe(":6783", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
