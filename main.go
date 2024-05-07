package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Group struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}
type Relationsstr struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", Artists)
	http.HandleFunc("/relations", RelationsHandler)
	fmt.Println("Port Çalışıyor...")
	http.ListenAndServe(":8080", nil)
}
func RelationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // webden gelen değer post değilse hata yazsın.
		http.Error(w, "Sadece POST metodu kabul edilmektedir", http.StatusMethodNotAllowed)
		return
	}
	artistID := r.FormValue("id") // webden gelen değer post olduğunda ismi id olan veriden gelen değri alır.
	if artistID == "" {
		http.Error(w, "hata", http.StatusBadRequest)
		return
	}
	var relationInfo Relationsstr                                                               // strcut yapısnı kullanmak için değişken atandı
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + artistID) //webden gelen değeri relationdaki değeri doğru okumak için relations ana linkine ekledik.
	if err != nil {
		fmt.Printf("HTTP isteği başarısız: %s\n", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		jsonData, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("HTTP cevabı okunamadı: %s\n", err)
		}
		err = json.Unmarshal(jsonData, &relationInfo)
		if err != nil {
			fmt.Println("error :", err)
		}
	}
	tmpl, err := template.ParseFiles("locations.html") // html i alır şablona dönüşrütürür.
	if err != nil {
		log.Println("Error parsing HTML template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, &relationInfo) // html şablonunun içine iletilmesini , yazılmasını sağlar.
	if err != nil {
		log.Println("Writing Error...", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func Artists(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(url) // sunucudan veri alınması gerektiği için get kullandık.
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close() // .body => isteğin yanıt olarak aldığı verileri içerir. .close ile bunu kapatırız

	var groups []Group                        //strcut yapsını kullanmak için
	if response.StatusCode == http.StatusOK { // statuscode kodların işeleme ve hata durumlarını içerir.
		body, err := io.ReadAll(response.Body) // okumak için kullanılır.
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, &groups) // json verilerini struct yapısına dönüştürür.
		if err != nil {
			log.Fatal(err)
		}
	}
	tmpl, err := template.ParseFiles("index.html") // html i alır şablona dönüşrütürür.
	if err != nil {
		log.Println("Error parsing HTML template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, groups) // html şablonunun içine iletilmesini , yazılmasını sağlar.
	if err != nil {
		log.Println("Writing Error...", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
