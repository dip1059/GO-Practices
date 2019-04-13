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



func main() {
	resp,_ := http.Get("http://localhost:2000/json/Dipankar Saha/Male/25")
	Bytes,_ := ioutil.ReadAll(resp.Body)
	var d Data
	_ = json.Unmarshal(Bytes, &d)
	_, _ = fmt.Println(d)
}
