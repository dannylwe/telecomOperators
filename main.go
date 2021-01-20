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
	Phone   string `json:"phone"`
	Country string `json:"country"`
}

func getProviders() map[string][]string {
	providers := make(map[string][]string)
	providers["uganda"] = []string{
		"MTN Uganda", "Airtel Uganda", "Uganda Telecom",
		"Africell Uganda", "Smile Telecom", "Vodafone Uganda",
		"Lycamobile Uganda"}
	providers["kenya"] = []string{"Safaricom", "Airtel", "Telkom Kenya"}
	providers["tanzania"] = []string{"Vodacom", "tiGo", "Airtel", "Viettel"}
	return providers
}

func getPrefixes(country string) map[string]string {
	operator := make(map[string]string)
	if country == "" || strings.ToLower(country) == "uganda" {
		operator["mtn"] = "077, 078, 039"
		operator["airtel"] = "075, 070"
		operator["africell"] = "079"
		operator["utl"] = "071"
	}

	if strings.ToLower(country) == "kenya" {
		operator["safaricom"] = "701, 702, 703, 704, 705, 706, 707, 708, 709, 710, 711, 712, 713, 714, 715, 716, 717, 718, 719, 720, 721, 722, 723, 724, 725"
		operator["airtel"] = "730, 731, 732, 733, 734, 735, 736, 737, 738, 739"
		operator["Telkom Kenya"] = "770, 771, 772, 773, 774, 775, 776, 777, 778, 779"
	}
	return operator
}

func getLine(number, country string) string {
	prefix := getPrefixes(country)

	MTN := "mtn"
	AIRTEL := "airtel"
	UTL := "utl"
	AFRICEL := "africell"
	SAFARICOM := "safaricom"
	TELKOM := "Telkom Kenya"
	providerUnknown := "Unknown provider"
	notSupported := "Number not supported"
	insufficient := "Insufficient digits"
	Unknown := "Unknown format"

	if len(number) < 10 {
		log.Println("Insufficient digits " + number)
		return insufficient
	}

	if len(number) > 13 {
		log.Println("Unknown format " + number)
		return Unknown
	}

	if len(number) == 10 {
		if country == "" || strings.ToLower(country) == "uganda" {

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
		}

		if strings.ToLower(country) == "kenya" {
			num := number[1:4]
			safariVal := prefix["safaricom"]
			airtelVal := prefix["airtel"]
			telkomVal := prefix["Telkom Kenya"]

			safaricom := strings.Contains(safariVal, num)
			airtel := strings.Contains(airtelVal, num)
			telekom := strings.Contains(telkomVal, num)

			if safaricom {
				return SAFARICOM
			}

			if airtel {
				return AIRTEL
			}

			if telekom {
				return TELKOM
			}
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
	prefixes := getPrefixes("")
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
	result := getLine(PhoneNumber.Phone, PhoneNumber.Country)
	log.Println("\nphone number: " + PhoneNumber.Phone + "\nCarrier: " + result + "\nCountry: " + strings.ToLower(PhoneNumber.Country))
	temp := make(map[string]string)
	temp["carrier"] = result

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temp)
}
