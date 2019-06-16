package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"net/http"
	"os"
)

type Params struct {
	ServiceName string `json:"serviceName"`
	NewTag      string `json:"newTag"`
	Token       string `json:"token"`
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
		serviceID := params.ServiceName
		serv, _, err := cli.ServiceInspectWithRaw(ctx, serviceID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		image := serv.Spec.TaskTemplate.ContainerSpec.Image
		version := serv.Meta.Version
		w.Write([]byte(image))
		ctx = r.Context()

		newSpec := serv.Spec
		//newSpec.
		opts := types.ServiceUpdateOptions{}
		updResp, err := cli.ServiceUpdate(ctx, serviceID, version, newSpec, opts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", updResp.Warnings)))
		return
	})

	fmt.Println("will listen ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
