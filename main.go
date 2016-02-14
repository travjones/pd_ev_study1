package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SubjectData struct {
	AParam  float64    `json:"aValue"`
	Delays  [7]float64 `json:"delays"`
	Indiffs [7]float64 `json:"indiffVals"`
	KParam  float64    `json:"kValue"`
}

type StoredData struct {
	ID   int    `db:"id"`
	Data string `db:"data"`
}

var cachedData []byte

func dataPost(w http.ResponseWriter, r *http.Request) {

	// request body as byte slice, need byte slice for unmarshal
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var s SubjectData

	// unmarshal body into s
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(s.KParam)

	// marshal json before storing in db
	sjson, err := json.Marshal(s)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(sjson)

	cachedData = sjson

	// // select all records in subject_data
	// sd := []StoredData{}
	// err = db.Select(&sd, "select * from subject_data")
	// if err != nil {
	// 	log.Fatalln(err)
	// 	return
	// }

	// fmt.Println(sd[0])
}

func dataPost2(w http.ResponseWriter, r *http.Request) {
	// request body as byte slice, need byte slice for unmarshal
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var s SubjectData

	// unmarshal body into s
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln(err)
	}

	// marshal json before storing in db
	sjson, err := json.Marshal(s)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(sjson)

	// connect to db
	db, err := sqlx.Connect("postgres", "user=travisjones dbname=pd sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// insert subject data json into db
	_, err = db.Exec("insert into data (data, ev_data) values ($1, $2)", cachedData, sjson)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/data", dataPost)
	http.HandleFunc("/data2", dataPost2)
	http.ListenAndServe(":3000", nil)
}
