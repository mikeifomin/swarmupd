package main

import (
	"errors"
	"log"

	"github.com/mikeifomin/swarmupd/server"

	"github.com/kelseyhightower/envconfig"
)

var ErrTagNotFound = errors.New("tag not found")
var ErrWrongToken = errors.New("wrong token")

type Params struct {
	ServiceName string `json:"ServiceName"`
	NewImage    string `json:"NewImage"`
	Token       string `json:"Token"`
}

type Env struct {
	TOKENS                 []string `required:"true"`
	PORT                   string   `required:"true"`
	REGISTRY_USER          string   `required:"true"`
	REGISTRY_PASSWORD      string   `required:"true"`
	SERVICE_PREFIXIES_ONLY []string `required:"true"`
}

func main() {
	var env Env
	envconfig.MustProcess("", &env)

	srv := server.Server{
		Addr:             ":" + env.PORT,
		Tokens:           env.TOKENS,
		RegistryUser:     env.REGISTRY_USERNAME,
		RegistryPassword: env.REGISTRY_PASSWORD,

		AllowedServiceIdPrefixies: env.SERVICE_PREFIXIES_ONLY,
	}
	err := srv.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
