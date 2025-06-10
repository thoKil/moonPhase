/*****************************************************************
* Moon Phase Calender
*
* Return moon phase as json over API.
* Show the moon phase on HTML website.
*
* author: thokil
* last update: 10.06.2025
******************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type MoonPhase struct {
	Phase        string  `json:"Phase"`
	Illumination float64 `json:"Illumination"`
}

func main() {
	fmt.Println("Hello Moon")
	//Call funnctions for this urls
	http.HandleFunc("/api/moonphase", apiHandler)
	http.HandleFunc("/", htmlHandler)
	port := os.Getenv("Port")
	if port == "" {
		port = "8081"
	}
	log.Printf("Der Mondphasen-Kalender läuft auf http://localhost:%s", port)
	//Starts webserver on port 8081
	http.ListenAndServe(":"+port, nil)
}

func apiHandler(writer http.ResponseWriter, request *http.Request) {
	//Get the date from the url (...?date=dd-mm-yyyy)
	date := request.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("10-06-2025")
	}
	//Get the Moon Phase
	moon, err := getMoonPhase(date)
	if err != nil {
		http.Error(writer, "Error to get moon phase", http.StatusInternalServerError)
		return
	}
	//return the moon phase as json
	writer.Header().Set("ContentType", "application/json")
	json.NewEncoder(writer).Encode(moon)
}

// show HTML page
func htmlHandler(writer http.ResponseWriter, request *http.Request) {
	date := time.Now().Format("10-06-2025")
	moon, err := getMoonPhase(date)
	if err != nil {
		http.Error(writer, "Error to get moon phase", http.StatusInternalServerError)
		return
	}
	template, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(writer, "Error on template", http.StatusInternalServerError)
		return
	}
	template.Execute(writer, moon)
}

// return moon phase (Fakedata)
func getMoonPhase(date string) (MoonPhase, error) {
	return MoonPhase{
		Phase:        "Full Moon",
		Illumination: 100,
	}, nil
}
