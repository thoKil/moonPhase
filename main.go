/*****************************************************************
* Moon Phase Calender
*
* Return moon phase as json over API.
* Show the moon phase on HTML website.
*
* !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
* IMPORTANT: To run this application, call:
* export PORT=8081 (or another port)
* export API_KEY=33902b95ee204b5eb31165429251106
* go run main.go
* !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
*
* Images from: https://www.flaticon.com/free-icons/moon-phase
*
* author: thokil
* last update: 12.06.2025
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
	Phase        string `json:"Phase"`
	Illumination int    `json:"Illumination"`
	Image        string `json:"Image"`
}

func main() {
	http.HandleFunc("/api/moonphase", apiHandler)
	http.HandleFunc("/", htmlHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// export PORT=8082 (or another port)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("The moon-phase-App runs on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func apiHandler(writer http.ResponseWriter, request *http.Request) {
	moon, err := fetchMoonData(request)
	if err != nil {
		http.Error(writer, "Error to get moon phase", http.StatusInternalServerError)
		return
	}
	//return the moon phase as json
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(moon)
}

// show HTML page
func htmlHandler(writer http.ResponseWriter, request *http.Request) {
	moon, err := fetchMoonData(request)
	if err != nil {
		http.Error(writer, "Error to get moon phase", http.StatusInternalServerError)
		return
	}
	template, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(writer, "Error to load template", http.StatusInternalServerError)
		return
	}
	if err := template.Execute(writer, moon); err != nil {
		log.Println("Template execute Error: ", err)
	}
}

func fetchMoonData(request *http.Request) (MoonPhase, error) {
	date := request.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	moon, err := getMoonPhase(date)
	if err != nil {
		return moon, err
	}
	moon.Phase = translatePhase(moon.Phase)
	return moon, nil
}

// return moon phase
func getMoonPhase(date string) (MoonPhase, error) {

	// export API_KEY=33902b95ee204b5eb31165429251106
	// go run main.go

	//apiKey := os.Getenv("API_KEY")
	apiKey := "33902b95ee204b5eb31165429251106"
	if apiKey == "" {
		return MoonPhase{}, fmt.Errorf("Api Key not set")
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/astronomy.json?key=%s&q=Solingen&dt=%s", apiKey, date)

	response, err := http.Get(url)
	if err != nil {
		return MoonPhase{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return MoonPhase{}, fmt.Errorf("API Error %s", response.Status)
	}

	var result struct {
		Astronomy struct {
			Astro struct {
				MoonPhase        string `json:"moon_phase"`
				MoonIllumination int    `json:"moon_illumination"`
			} `json:"astro"`
		} `json:"astronomy"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return MoonPhase{}, err
	}

	return MoonPhase{
		Phase:        result.Astronomy.Astro.MoonPhase,
		Illumination: result.Astronomy.Astro.MoonIllumination,
		Image:        getPhaseImage(result.Astronomy.Astro.MoonPhase),
	}, nil
}

func translatePhase(phase string) string {
	switch phase {
	case "New Moon":
		return "Neumond"
	case "Waxing Crescent":
		return "Zunehmender Sichelmond"
	case "First Quarter":
		return "Zunehmender Halbmond"
	case "Waxing Gibbous":
		return "Zunehmender Dreiviertelmond"
	case "Full Moon":
		return "Vollmond"
	case "Waning Gibbous":
		return "Abnehmender Dreiviertelmond"
	case "Last Quarter":
		return "Abnehmender Halbmond"
	case "Waning Crescent":
		return "Abnehmendert Sichelmond"
	default:
		return phase
	}
}

func getPhaseImage(phase string) string {
	switch phase {
	case "New Moon":
		return "/static/new_moon.png"
	case "Waxing Crescent":
		return "/static/waxing_crescent.png"
	case "First Quarter":
		return "/static/first_quarter.png"
	case "Waxing Gibbous":
		return "/static/waxing_gibbous.png"
	case "Full Moon":
		return "/static/full_moon.png"
	case "Waning Gibbous":
		return "/static/waning_gibbous.png"
	case "Last Quarter":
		return "/static/last_quarter.png"
	case "Waning Crescent":
		return "/static/waning_crescent.png"
	default:
		return ""
	}
}
