package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os/exec"
	"strconv"
	"strings"
)

type server struct{}

type MemInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Proc struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	RID   int    `json:"rid"`
	State int    `json:"state"`
	Mem   string `json:"mem"`
	User  string `json:"user"`
}

type ApiResponse struct {
	TotalRam  MemInfo `json:"totalRam"`
	FreeRam   MemInfo `json:"freeRam"`
	ProcArray []Proc  `json:"procArray"`
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") ///Para cualquier tipo de peticion seteamos el header contenttype a json
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "GET":
		file := r.URL.Query()["path"][0]
		w.WriteHeader(http.StatusOK)
		fileCnt := readFile(file)
		retObject := getRetObject(fileCnt)
		fmt.Println("Printing retObj")
		fmt.Println(retObject)
		jsonObject, err := json.Marshal(retObject)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonObject)

	case "POST":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "method not implemented"}`))
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "method not implemented"}`))

	}
}

func main() {
	s := &server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	check(err)
	fmt.Print(string(dat))
	return string(dat)

}

func getRetObject(fileString string) ApiResponse {
	var retObj = ApiResponse{}
	var procArray []Proc
	for index, proc := range strings.Split(fileString, "\n") {
		pieces := strings.Split(proc, ":")
		if index <= 1 {

			if strings.Contains(pieces[0], "total") {
				retObj.TotalRam = MemInfo{"Total Ram", pieces[1]}
			} else {
				retObj.FreeRam = MemInfo{"Free Ram", pieces[1]}
			}
		}
		if index >= 3 {

			//fmt.Print(index)
			//fmt.Print(pieces)
			id, _ := strconv.Atoi(pieces[0])
			rid, _ := strconv.Atoi(pieces[2])
			state, _ := strconv.Atoi(pieces[3])
			procObj := Proc{Id: id, Name: pieces[1], RID: rid, State: state}
			//out, _ := exec.Command("id", "-nu", pieces[0]).Output()
			//procObj.User = strings.Trim(string(out), "\n")
			//out2, _ := exec.Command("ps", "-p", pieces[1], "-o", "%mem").Output()
			//procObj.Mem = strings.TrimLeft(strings.Split(string(out2), "\n")[1], " ")
			fmt.Print(procObj)
			procArray = append(procArray, procObj)
		}

		retObj.ProcArray = procArray
	}
	return retObj
}
