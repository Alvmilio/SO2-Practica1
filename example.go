package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
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
		//file := r.URL.Query()["path"][0]
		w.WriteHeader(http.StatusOK)
		fileCnt := readFile("/proc/practica1")
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
		killProcess(r.URL.Query()["pid"][0])
		w.Write([]byte(`{"message": "Ok"}`))
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

func killProcess(pid string) {
	//realPid, _ := strconv.Atoi(pid)
	out, _ := exec.Command("kill", "-9", pid).Output()
	fmt.Println("Kill output: " + strings.Trim(string(out), "\n"))

}

func getRetObject(fileString string) ApiResponse {
	var retObj = ApiResponse{}
	var procArray []Proc
	var procs = strings.Split(fileString, "\n")

	for i := 0; i < len(procs)-1; i++ {
		var proc = procs[i]
		pieces := strings.Split(proc, ":")
		fmt.Println(i)
		if i <= 1 {
			fmt.Println("entro 1")
			if strings.Contains(pieces[0], "total") {
				retObj.TotalRam = MemInfo{"Total Ram", pieces[1]}
			} else {
				retObj.FreeRam = MemInfo{"Free Ram", pieces[1]}
			}
		}
		if i >= 3 {
			fmt.Println("entro 2")

			//fmt.Print(index)
			//fmt.Print(pieces)
			id, _ := strconv.Atoi(pieces[0])
			rid, _ := strconv.Atoi(pieces[2])
			state, _ := strconv.Atoi(pieces[3])
			procObj := Proc{Id: id, Name: pieces[1], RID: rid, State: state}
			out, _ := exec.Command("ps", "-o", "user=", "-p", pieces[0]).Output()
			procObj.User = strings.Trim(string(out), "\n")
			out2, _ := exec.Command("ps", "-p", pieces[0], "-o", "%mem").Output()
			procObj.Mem = strings.TrimLeft(strings.Split(string(out2), "\n")[1], " ")
			fmt.Print(procObj)
			procArray = append(procArray, procObj)
		}

		retObj.ProcArray = procArray
	}
	return retObj
}
