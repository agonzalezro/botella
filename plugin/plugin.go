package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"

	"github.com/agonzalezro/botella/utils"
)

type Plugin struct {
	Image string

	client    *docker.Client
	container *docker.Container

	environment map[string]string

	RunOnlyOnChannels       bool
	RunOnlyOnDirectMessages bool
	RunOnlyOnMentions       bool
}

type Input struct {
	Version  int    `json:"version,omitempty"`
	Emitter  string `json:"emitter,omitempty"`
	Receiver string `json:"receiver,omitempty"`
	Body     string `json:"body"`
}

func NewInput(emitter, receiver, body string) Input {
	return Input{
		Version:  -1, // TODO: not used yet
		Emitter:  emitter,
		Receiver: receiver,
		Body:     body,
	}
}

func (i Input) JSON() string {
	b, err := json.Marshal(i)
	if err != nil {
		log.Warning(err)
		return ""
	}
	return string(b)
}

func environmentAsArrayOfString(image string, environment map[string]string) []string {
	var (
		arrayOfEnvs []string
		err         error
	)
	for k, v := range environment {
		// We want to override it with a value from the environment
		if v == "" {
			v, err = utils.GetFromEnvOrFromMap(image, nil, k)
			log.Warning(err)
		}
		arrayOfEnvs = append(arrayOfEnvs, fmt.Sprintf("%s=%s", k, v))
	}
	return arrayOfEnvs
}

func New(image string, environment map[string]string) (*Plugin, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	// TODO: don't always pull images, check imagePullPolicy from yaml
	if err := client.PullImage(
		docker.PullImageOptions{Repository: image},
		docker.AuthConfiguration{},
	); err != nil {
		return nil, err
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        image,
			Env:          environmentAsArrayOfString(image, environment),
			AttachStdin:  true, // TODO: not sure what of these are needed
			AttachStdout: true,
			OpenStdin:    true,
			StdinOnce:    true,
		},
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("Plugin/Container (%s) created: %+v", image, container)
	return &Plugin{
		Image:       image,
		client:      client,
		container:   container,
		environment: environment,
	}, nil
}

func (p *Plugin) Stop() error {
	return p.client.RemoveContainer(
		docker.RemoveContainerOptions{ID: p.container.ID, Force: true})
}

func (p *Plugin) Run(input Input) (string, error) {
	// TODO: not sure if we should do this or keep an ongoing container running
	if err := p.client.StartContainer(p.container.ID, nil); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := p.client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    p.container.ID,
		Stdin:        true,
		Stdout:       true,
		InputStream:  strings.NewReader(input.JSON()),
		OutputStream: &buf,
		Stream:       true,
	}); err != nil {
		return "", err
	}

	if _, err := p.client.WaitContainer(p.container.ID); err != nil {
		return "", err
	}

	return buf.String(), nil
}
