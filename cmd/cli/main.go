package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mikeifomin/swarmupd/server"
)

func main() {
	serviceId := flag.String("service-id", "", "")
	newImage := flag.String("new-image", "", "")
	token := flag.String("token", "", "")
	url := flag.String("url", "", "")
	flag.Parse()

	params := server.Params{
		ServiceId: *serviceId,
		NewImage:  *newImage,
		Token:     *token,
	}
	b, _ := json.Marshal(params)
	resp, err := http.Post(*url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	body := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		log.Fatal(body)
	}
	fmt.Println(body)
}
