package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/julienschmidt/httprouter"
	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/cors"
)

func typeListRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	list := make([]string, 0, len(typeMap))
	for t := range typeMap {
		list = append(list, t)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func typeSpecRoute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	typeName := ps.ByName("typeName")
	typeSpec, ok := typeMap[typeName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("party type not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(typeSpec)
}

// NewParty ...
type NewParty struct {
	PartyType string          `json:"party_type"`
	Data      json.RawMessage `json:"data"`
}

func newPartyRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	var p NewParty
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var party interface{}
	var err error
	switch p.PartyType {
	case "MovieParty":
		var mp MovieParty
		err = json.Unmarshal(p.Data, &mp)
		party = mp
	case "PoolParty":
		var pp PoolParty
		err = json.Unmarshal(p.Data, &pp)
		party = pp
	case "DinnerParty":
		var dp DinnerParty
		err = json.Unmarshal(p.Data, &dp)
		party = dp
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown Party Type"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error unmarshalling party: " + err.Error()))
		return
	}
	if err := validate.Struct(party); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error validating party: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(shortuuid.New()))
}

func failSometimes(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if rand.Intn(100) < 20 {
			defer r.Body.Close()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Random prod failure: " + petname.Generate(3, " ")))
			return
		}
		h(w, r, p)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	router := httprouter.New()
	router.GET("/", aboutRoute)
	router.GET("/partytypes", typeListRoute)
	router.GET("/partytype/:typeName", typeSpecRoute)
	router.POST("/bookparty", newPartyRoute)
	router.POST("/bookpartyprod", failSometimes(newPartyRoute))
	corsHandler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), corsHandler))
}
