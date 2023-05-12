package reporter

import (
	"testing"

	"github.com/nextlinux/k8s-inventory/internal/config"
)

func TestBuildUrl(t *testing.T) {
	nextlinuxDetails := config.NextlinuxInfo{
		URL:      "https://ancho.re",
		User:     "admin",
		Password: "foobar",
	}

	expectedURL := "https://ancho.re/v1/enterprise/kubernetes-inventory"
	actualURL, err := buildURL(nextlinuxDetails)
	if err != nil || expectedURL != actualURL {
		t.Errorf("Failed to build URL:\nexpected=%s\nactual=%s", expectedURL, actualURL)
	}
}
