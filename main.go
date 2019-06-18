package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"net/http"
	"os"
	"strings"
)

var ErrTagNotFound = errors.New("tag not found")

type Params struct {
	ServiceName string `json:"serviceName"`
	NewTag      string `json:"newTag"`
	Token       string `json:"token"`
}

func main() {
	_, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("Set env TOKEN")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")

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
		imageFullName := serv.Spec.TaskTemplate.ContainerSpec.Image
		version := serv.Meta.Version
		w.Write([]byte(imageFullName))

		ctx = r.Context()
		newSpec := serv.Spec
		image := NewImage(imageFullName)
		err = image.UpdateTag(params.NewTag, username, password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		newSpec.TaskTemplate.ContainerSpec.Image = image.String()
		w.Write([]byte("\n" + image.String()))
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

type Image struct {
	Registry string
	Name     string
	Tag      string
	Digest   string
}

func (i *Image) String() string {
	return i.Registry + "/" + i.Name + ":" + i.Tag + "@" + i.Digest
}

func NewImage(full string) *Image {
	posFirstSlash := strings.Index(full, "/")
	posFirstColon := strings.Index(full, ":")
	posFirstAt := strings.Index(full, "@")
	i := Image{
		Registry: full[:posFirstSlash],
		Name:     full[posFirstSlash+1 : posFirstColon],
		Tag:      full[posFirstColon+1 : posFirstAt],
		Digest:   full[posFirstAt+1:],
	}
	return &i
}

func (i *Image) UpdateTag(newTag, login, password string) error {
	client := &http.Client{}
	url := "https://" + i.Registry + "/v2/" + i.Name + "/manifests/" + newTag
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if login != "" {
		req.SetBasicAuth(login, password)
	}
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrTagNotFound
	}

	i.Tag = newTag
	i.Digest = resp.Header.Get("Docker-Content-Digest")
	return nil
}
