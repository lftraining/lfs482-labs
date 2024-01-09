package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Vessel struct {
	gorm.Model
	Name string
}

var db *gorm.DB
var err error

func main() {
	start := time.Now()
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=8085 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	db.AutoMigrate(&Vessel{})

	db.Delete(&Vessel{})
	db.Create(&Vessel{Name: "Maritime Mover"})
	db.Create(&Vessel{Name: "Tidal Transporter"})
	db.Create(&Vessel{Name: "Cargo Clipper"})

	router := mux.NewRouter()

	router.Path("/").Handler(http.FileServer(http.Dir(".")))

	router.HandleFunc("/api/vessels", ListVessels).Methods("GET")
	router.HandleFunc("/api/vessels/{id:[0-9]+}", GetVessel).Methods("GET")
	router.HandleFunc("/api/vessels", AddVessel).Methods("POST")
	router.HandleFunc("/api/vessels/{id:[0-9]+}", EditVessel).Methods("PUT")
	router.HandleFunc("/api/vessels/{id:[0-9]+}", DeleteVessel).Methods("DELETE")

	log.Println("Vessel Manifest Server starting at:")
	log.Println(time.Since(start))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func ListVessels(w http.ResponseWriter, r *http.Request) {
	var vessel []Vessel
	db.Find(&vessel)
	if err := json.NewEncoder(w).Encode(vessel); err != nil {
		logError(err)
	}
}

func GetVessel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vesselId := vars["id"]

	if !checkIfVesselExists(vesselId) {
		if err := json.NewEncoder(w).Encode("Vessel Manifest Not Found!"); err != nil {
			logError(err)
		}
		return
	}

	var vessel Vessel
	db.Where("id = ?", vesselId).First(&vessel)
	if err := json.NewEncoder(w).Encode(vessel); err != nil {
		logError(err)
	}
}

func AddVessel(w http.ResponseWriter, r *http.Request) {
	var vessel Vessel
	if err := json.NewDecoder(r.Body).Decode(&vessel); err != nil {
		logError(err)
	}

	db.Create(&vessel)
	w.WriteHeader(204)
}

func EditVessel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vesselId := vars["id"]

	var vessel Vessel
	db.Where("id = ?", vesselId).First(&vessel)
	if err := json.NewDecoder(r.Body).Decode(&vessel); err != nil {
		logError(err)
	}

	db.Save(&vessel)
}

func DeleteVessel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vesselId := vars["id"]
	db.Where("id = ?", vesselId).Delete(&Vessel{})
}

func checkIfVesselExists(vesselId string) bool {
	var vessel Vessel
	db.First(&vessel, vesselId)
	return vessel.Name != ""
}

func logError(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
}
