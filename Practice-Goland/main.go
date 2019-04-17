package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Data struct {
	Name string `json:"name"`
	Sex string
	Age int
}

var mp =make(map[string]*Data)

func main() {
	resp,_ := http.Get("http://localhost:2000/json/Dipankar Saha/Male/25")
	Bytes,_ := ioutil.ReadAll(resp.Body)
	var d []*Data
	json.Unmarshal(Bytes, &d)
	for _, val := range d {
		mp[val.Name] = val
		fmt.Println(*mp[val.Name])
	}
	//mp["Dipanka Saha"].Name = "DK"
}


