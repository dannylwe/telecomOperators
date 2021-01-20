package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	// providers := getProviders()

	// m := getLine("0772504991")
	// fmt.Println(m)

	const PORT string = ":9000"
	log.Println("starting application on " + PORT)
	router := mux.NewRouter()
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/prefix", listPrefixes).Methods("GET")
	router.HandleFunc("/providers", listProviders).Methods("GET")
	router.HandleFunc("/carrier", getCarrier).Methods("POST")
	log.Fatal(http.ListenAndServe(PORT, router))
}

type phoneNumber struct {
	Phone string `json:"phone"`
}

func getProviders() map[string][]string {
	providers := make(map[string][]string)
	providers["uganda"] = []string{
		"MTN Uganda", "Airtel Uganda", "Uganda Telecom",
		"Africell Uganda", "Smile Telecom", "Vodafone Uganda",
		"Lycamobile Uganda"}
	return providers
}

func getPrefixes() map[string]string {
	m := make(map[string]string)
	m["mtn"] = "077, 078, 039"
	m["airtel"] = "075, 070"
	m["africell"] = "079"
	m["utl"] = "071"
	return m
}

func getLine(number string) string {
	prefix := getPrefixes()

	MTN := "mtn"
	AIRTEL := "airtel"
	UTL := "utl"
	AFRICEL := "africell"
	providerUnknown := "Unknown provider"
	notSupported := "Number not supported"

	if len(number) == 10 {
		num := number[0:3]
		mtnVal := prefix["mtn"]
		airtelVal := prefix["airtel"]
		africellVal := prefix["africell"]
		utlVal := prefix["utl"]

		mtn := strings.Contains(mtnVal, num)
		airtel := strings.Contains(airtelVal, num)
		africell := strings.Contains(africellVal, num)
		utl := strings.Contains(utlVal, num)

		if mtn {
			return MTN
		}

		if airtel {
			return AIRTEL
		}

		if africell {
			return AFRICEL
		}

		if utl {
			return UTL
		}
		log.Println("unknown provider for: " + number)
		return providerUnknown
	}
	log.Println("unsupported number: " + number)
	return notSupported
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tel Provider is healthy")
}

func listPrefixes(w http.ResponseWriter, r *http.Request) {
	prefixes := getPrefixes()
	payload, _ := json.Marshal(prefixes)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

func listProviders(w http.ResponseWriter, r *http.Request) {
	providers := getProviders()
	payload, _ := json.Marshal(providers)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

func getCarrier(w http.ResponseWriter, r *http.Request) {
	var PhoneNumber phoneNumber
	_ = json.NewDecoder(r.Body).Decode(&PhoneNumber)
	result := getLine(PhoneNumber.Phone)
	log.Println("\nphone number: " + PhoneNumber.Phone + "\nCarrier: " + result)
	temp := make(map[string]string)
	temp["carrier"] = result

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temp)
}
