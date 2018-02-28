// +build docker

package docker_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/yext/edward/services"
	"github.com/yext/edward/services/backends/docker"
)

func TestMain(m *testing.M) {
	// Register necessary backends
	services.RegisterDefaultBackend(&docker.Loader{})

	os.Exit(m.Run())
}

func TestStart(t *testing.T) {
	service := &services.ServiceConfig{
		Name: "testservice",
		BackendConfig: &docker.Backend{
			Repository: "hello-world-nginx",
			Tag:        "latest",
			Ports: map[string]string{
				"80": "51432",
			},
		},
	}
	b, err := services.GetRunner(service)
	if err != nil {
		t.Error(err)
		return
	}

	err = b.Start(os.Stdout, os.Stderr)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_, err = b.Stop("", nil)
		if err != nil {
			t.Error(err)
		}
	}()

	resp, err := http.Get("http://localhost:51432/")
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(body), "nginx container") {
		t.Errorf("Response was not as expected:\n%s", string(body))
	}

}
