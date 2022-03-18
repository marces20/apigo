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

var varSatellites Satellites
var xySatelites = make(map[string][2]float32)

func topsecret(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)

	var post Satellites

	error2 := json.Unmarshal(reqBody, &post)
	if error2 != nil {
		respondWithError(w, http.StatusBadRequest, "No se puede establecer la ubicacion")
		return
	}

	mx, my := GetLocation(post.Satellites[0].Distance, post.Satellites[1].Distance, post.Satellites[2].Distance)
	resp := make(map[string]interface{})

	p := make(map[string]float32)
	p["x"] = mx
	p["y"] = my
	if mx == 0.0 || my == 0.0 {
		respondWithError(w, http.StatusBadRequest, "No se puede establecer la ubicacion")
		return
	}
	resp["position"] = p

	resp["message"] = GetMessage(post.Satellites[0].Message, post.Satellites[1].Message, post.Satellites[2].Message)

	respondWithJSON(w, http.StatusOK, resp)

}

func topsecretSplit(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	satellite_name := mux.Vars(r)["satellite_name"]

	var post Satelite

	error2 := json.Unmarshal(reqBody, &post)
	if error2 != nil {
		respondWithError(w, http.StatusBadRequest, "No se puede establecer la ubicacion")
		return
	}

	post.Name = satellite_name

	var tempVarSatellites []Satelite

	for _, sal := range varSatellites.Satellites {
		if sal.Name != satellite_name {
			tempVarSatellites = append(tempVarSatellites, sal)
		}
	}
	varSatellites.Satellites = tempVarSatellites
	varSatellites.Satellites = append(varSatellites.Satellites, post)

	if len(varSatellites.Satellites) == 3 {
		mx, my := GetLocation(varSatellites.Satellites[0].Distance, varSatellites.Satellites[1].Distance, varSatellites.Satellites[2].Distance)
		resp := make(map[string]interface{})

		p := make(map[string]float32)
		p["x"] = mx
		p["y"] = my
		if mx == 0.0 || my == 0.0 {
			respondWithError(w, http.StatusBadRequest, "No se puede establecer la ubicacion")
			return
		}
		resp["position"] = p

		resp["message"] = GetMessage(varSatellites.Satellites[0].Message, varSatellites.Satellites[1].Message, varSatellites.Satellites[2].Message)
		respondWithJSON(w, http.StatusOK, resp)
	} else {
		respondWithError(w, http.StatusBadRequest, "No se puede establecer la ubicacion. No hay suficiente información")
		return
	}

}

func GetMessage(messages ...[]string) (msg string) {

	var mess = make([]string, len(messages[0]))
	for _, mensaje := range messages {
		for x, dato := range mensaje {
			if mess[x] == "" {
				mess[x] = dato
			}
		}
	}

	for _, x := range mess {
		msg += x + " "
	}

	return msg
}

func GetLocation(distances ...float32) (rx, ry float32) {
	//fmt.Println(xySatelites)
	var dist [3]float32

	for i, param := range distances {
		fmt.Println(param)
		dist[i] = param
	}
	//http://ramon-gzz.blogspot.com/2013/05/geolocalizacion.html
	//UBICACIONES DE LOS SATELITES Kenobi: [-500, -200] Skywalker: [100, -100] Sato: [500, 100]
	//c1, c2, c3 = (50,50), (300,430), (590,50)
	//P1 = [-500, -200]  //Kenobi    //# Almacenamos las coordenadas de cada transmidor
	//P2 = [100, -100]	//Skywalker
	//P3 = [500, 100]	//Sato

	var d float64 = float64(xySatelites["Skywalker"][0]) - float64(xySatelites["Kenobi"][0]) // 100 - (-500) //self.a2x - self.a1x
	var i float64 = float64(xySatelites["Sato"][0]) - float64(xySatelites["Kenobi"][0])      //500 - (-500) //self.a3x - self.a1x
	var j float64 = float64(xySatelites["Sato"][1]) - float64(xySatelites["Sato"][1])        // 100 - (-200) // self.a3y - self.a1y

	var r1 float64 = float64(dist[0]) //Kenobi
	var r2 float64 = float64(dist[1]) //Skywalker
	var r3 float64 = float64(dist[2]) //Sato

	var x float64 = (math.Pow(r1, 2) - math.Pow(r2, 2) + math.Pow(d, 2)) / (2 * d)
	var y float64 = (math.Pow(r1, 2) - math.Pow(r3, 2) - math.Pow(x, 2) + math.Pow((x-i), 2) + math.Pow(j, 2)) / (2 * j)
	x += -500 //self.a1x
	y += -200 //self.a1y

	/*var x1 = float64(xySatelites["Kenobi"][0])
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
	var y float64 = (C*D - A*F) / (B*D - A*E)*/

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

func main() {

	//UBICACIONES DE LOS SATELITES Kenobi: [-500, -200] Skywalker: [100, -100] Sato: [500, 100]
	xySatelites["Kenobi"] = [2]float32{-500, -200}
	xySatelites["Skywalker"] = [2]float32{100, -100}
	xySatelites["Sato"] = [2]float32{500, 100}

	r := mux.NewRouter()

	r.HandleFunc("/topsecret", topsecret).Methods("POST")
	r.HandleFunc("/topsecret_split/{satellite_name}", topsecretSplit).Methods("POST", "GET")

	log.Printf("Listening v1...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
