package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type urlDataStruct struct {
	OriginalUrl string
	ShortUrl    string
	CreatedAt   time.Time
}

var myDb = make(map[string]urlDataStruct)

func generateShortUrl(OriginalUrl string) string {
	hashData := md5.New()
	hashData.Write([]byte(OriginalUrl))
	data := hashData.Sum(nil)
	shortString := hex.EncodeToString(data)
	return shortString[:8]
	// return "Himanshu"
}

func createShorData(OriginalUrl string) string {
	shortString := generateShortUrl(OriginalUrl)
	myDb[shortString] = urlDataStruct{
		OriginalUrl: OriginalUrl,
		ShortUrl:    shortString,
		CreatedAt:   time.Now(),
	}
	return shortString
}

func getUrl(id string) (urlDataStruct, error) {
	data, ok := myDb[id]
	if !ok {
		return urlDataStruct{}, errors.New("url Not found")
	}
	return data, nil
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	out := createShorData(data.URL)
	// fmt.Fprintf(w, shortURL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: out}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getUrl(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

func main() {
	// OriginalUrl := "www.google.com"
	// shortUrl := createShorData(OriginalUrl)

	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server Not Starting|Something is Wrong", err)
	}
}
