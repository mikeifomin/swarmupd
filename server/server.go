package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Params struct {
	Token     string
	NewImage  string
	ServiceId string
}

type Server struct {
	Addr             string
	Tokens           []string
	RegistryUser     string
	RegistryPassword string

	AllowedServiceIdPrefixies []string

	docker *Docker
	server *http.Server
}

func (s *Server) Init() error {
	_, err := NewDocker()
	if err != nil {
		return err
	}
	s.server = &http.Server{
		Addr:    s.Addr,
		Handler: s,
	}

	if len(s.AllowedServiceIdPrefixies) == 0 {
		return fmt.Errorf("you must set allowed serviceIds prefixies")
	}
	return nil
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	params, err := reqToParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if s.IsTokenWrong(params.Token) {
		http.Error(w, "no access", http.StatusForbidden)
		return
	}

	if !s.isAllowedServiceId(params.ServiceId) {
		http.Error(w, "service id not allowed", http.StatusForbidden)
		return
	}

	err = s.UpdateDockerService(r.Context(), params.ServiceId, params.NewImage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (s *Server) IsTokenWrong(token string) bool {
	for _, allowed := range s.Tokens {
		if allowed == token {
			return false
		}
	}
	return true
}

func (s *Server) UpdateDockerService(ctx context.Context, serviceId, imageNext string) error {
	docker, err := NewDocker()
	if err != nil {
		return err
	}

	err = docker.PullImage(ctx, imageNext, s.RegistryUser, s.RegistryPassword)
	if err != nil {
		return err
	}

	imageWas, err := docker.GetImageFromService(serviceId)
	if err != nil {
		return err
	}

	if !imageTagChangedOrNoChange(imageWas, imageNext) {
		return fmt.Errorf("only tag update allowed. imageWas: %s, imageNext: %s", imageWas, imageNext)
	}

	err = docker.UpdateImageService(serviceId, imageNext, s.RegistryUser, s.RegistryPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) isAllowedServiceId(serviceId string) bool {
	for _, prefix := range s.AllowedServiceIdPrefixies {
		if strings.HasPrefix(serviceId, prefix) {
			return true
		}
	}
	return false
}

func reqToParams(r *http.Request) (*Params, error) {
	var params Params
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return nil, err
	}
	return &params, nil
}
