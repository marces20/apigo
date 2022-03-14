package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
)

type Satelite struct {
	Name     string   `json:"Name"`
	Distance float32  `json:"Distance"`
	Message  []string `json:"Message"`
}

type Satellites struct {
	Satellites []Satelite `json:"Satellites"`
}

func topsecret(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)

	var post Satellites

	json.Unmarshal(reqBody, &post)

	//json.NewEncoder(w).Encode(&post)

	//data, err := json.Marshal(post)

	//fmt.Println(string(data))
	//respondWithJSON(w, http.StatusOK, post)

	mx, my := GetLocation(post.Satellites[0].Distance, post.Satellites[1].Distance, post.Satellites[2].Distance)
	resp := make(map[string]interface{})

	p := make(map[string]float32)
	p["x"] = mx
	p["y"] = my

	resp["position"] = p
	resp["message"] = "Status Created"

	//jsonResp, _ := json.Marshal(resp)

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//w.Write(jsonResp)

	//fmt.Println(mx)
	///fmt.Println(my)

	//data, _ := json.Marshal(res)

	//if err != nil {
	//	respondWithError(w, http.StatusBadRequest, err.Error())
	//	return
	//}
	respondWithJSON(w, http.StatusOK, resp)

}

func GetMessage(messages ...[5]string) (msg string) {

	var mess [5]string

	for _, mensaje := range messages {

		for x, dato := range mensaje {
			fmt.Println(dato)
			if mess[x] == "" {
				mess[x] = dato
			}
		}
	}
	fmt.Println(mess)
	return "-------------"
}

func GetLocation(distances ...float32) (rx, ry float32) {
	fmt.Println(xySatelites)
	var dist [3]float32

	for i, param := range distances {
		fmt.Println(param)
		dist[i] = param
	}
	//UBICACIONES DE LOS SATELITES Kenobi: [-500, -200] Skywalker: [100, -100] Sato: [500, 100]

	var x1 = float64(xySatelites["Kenobi"][0])
	var y1 = float64(xySatelites["Kenobi"][1])
	var r1 = float64(dist[0])

	var x2 = float64(xySatelites["Skywalker"][0])
	var y2 = float64(xySatelites["Skywalker"][1])
	var r2 = float64(dist[1])

	var x3 = float64(xySatelites["Sato"][0])
	var y3 = float64(xySatelites["Sato"][1])
	var r3 = float64(dist[2])

	var A float64 = 2*x2 - 2*x1
	var B float64 = 2*y2 - 2*y1

	var C float64 = math.Pow(float64(r1), 2) - math.Pow(float64(r2), 2) - math.Pow(float64(x1), 2) + math.Pow(float64(x2), 2) - math.Pow(float64(y1), 2) + math.Pow(float64(y2), 2)
	var D float64 = 2*x3 - 2*x2
	var E float64 = 2*y3 - 2*y2

	var F float64 = math.Pow(float64(r2), 2) - math.Pow(float64(r3), 2) - math.Pow(float64(x2), 2) + math.Pow(float64(x3), 2) - math.Pow(float64(y2), 2) + math.Pow(float64(y3), 2)

	var x float64 = (C*E - F*B) / (E*A - B*D)
	var y float64 = (C*D - A*F) / (B*D - A*E)

	return float32(x), float32(y)

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

var xySatelites = make(map[string][2]float32)

func main() {

	myarray := [5]string{"", "este", "es", "un", "mensaje"}
	myarray2 := [5]string{"", "este", "", "un", "mensaje"}

	rx := GetMessage(myarray, myarray2)
	fmt.Println(rx)

	//UBICACIONES DE LOS SATELITES Kenobi: [-500, -200] Skywalker: [100, -100] Sato: [500, 100]
	xySatelites["Kenobi"] = [2]float32{-500, -200}
	xySatelites["Skywalker"] = [2]float32{100, -100}
	xySatelites["Sato"] = [2]float32{500, 100}

	r := mux.NewRouter()

	r.HandleFunc("/topsecret", topsecret).Methods("POST")
	r.HandleFunc("/", topsecret)

	log.Printf("Listening...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
