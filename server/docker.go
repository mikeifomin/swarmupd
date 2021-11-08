package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Docker struct {
	client *client.Client
}

func NewDocker() (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	_, err = cli.Ping(ctxDefault())
	if err != nil {
		return nil, err
	}

	d := Docker{
		client: cli,
	}
	return &d, nil
}

func (d *Docker) PullImage(ctx context.Context, image, user, password string) error {
	opt := types.ImagePullOptions{}
	if user != "" {
		opt.RegistryAuth = authStr(user, password)
	}
	_, err := d.client.ImagePull(ctx, image, opt)
	if err != nil {
		return err
	}
	return nil
}

func (d *Docker) GetImageFromService(serviceID string) (string, error) {
	service, _, err := d.client.ServiceInspectWithRaw(ctxDefault(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return "", err
	}
	return service.Spec.TaskTemplate.ContainerSpec.Image, nil
}

func (d *Docker) UpdateImageService(serviceID string, image string) error {
	service, _, err := d.client.ServiceInspectWithRaw(ctxDefault(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}
	newSpec := service.Spec
	newSpec.TaskTemplate.ContainerSpec.Image = image
	version := service.Meta.Version
	_, err = d.client.ServiceUpdate(ctxDefault(), serviceID, version, newSpec, types.ServiceUpdateOptions{})
	if err != nil {
		return err
	}
	return nil

}

func ctxDefault() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	return ctx
}

func authStr(user, password string) string {
	authConfig := types.AuthConfig{
		Username: user,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(encodedJSON)
}
