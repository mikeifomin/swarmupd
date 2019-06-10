package main

import (
	"encoding/json"
	"github.com/docker/docker/client"
	"log"
	"net/http"
	"os"
)

type Params struct {
	Name   string `json:"name"`
	NewTag string `json:"newTag"`
	Token  string `json:"token"`
}

func main() {

	token := os.Getenv("TOKEN")
	if token == "" {
		panic("Set env TOKEN")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var params Params
		err := json.NewDecoder(r.Body).Decode(&params)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if params.Token != token {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Wrong token"))
			return
		}

		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		ctx := r.Context()
		serv, _, err := cli.ServiceInspectWithRaw(ctx, params.Name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		image := serv.Spec.TaskTemplate.ContainerSpec.Image
		w.Write([]byte(image))

	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
