// Once In-Use Image data has been gathered, this package reports the data to Nextlinux
package reporter

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/nextlinux/k8s-inventory/internal/config"
	"github.com/nextlinux/k8s-inventory/internal/log"
	"github.com/nextlinux/k8s-inventory/internal/tracker"
	"github.com/nextlinux/k8s-inventory/pkg/inventory"
)

const ReportAPIPath = "v1/enterprise/kubernetes-inventory"

// This method does the actual Reporting (via HTTP) to Nextlinux
//
//nolint:gosec
func Post(report inventory.Report, nextlinuxDetails config.NextlinuxInfo) error {
	defer tracker.TrackFunctionTime(time.Now(), "Reporting results to Nextlinux for cluster: "+report.ClusterName+"")
	log.Debug("Reporting results to Nextlinux")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: nextlinuxDetails.HTTP.Insecure},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(nextlinuxDetails.HTTP.TimeoutSeconds) * time.Second,
	}

	nextlinuxURL, err := buildURL(nextlinuxDetails)
	if err != nil {
		return fmt.Errorf("failed to build url: %w", err)
	}

	reqBody, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to serialize results as JSON: %w", err)
	}

	req, err := http.NewRequest("POST", nextlinuxURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to build request to report data to Nextlinux: %w", err)
	}
	req.SetBasicAuth(nextlinuxDetails.User, nextlinuxDetails.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-nextlinux-account", nextlinuxDetails.Account)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to report data to Nextlinux: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("failed to report data to Nextlinux: %+v", resp)
	}
	log.Debug("Successfully reported results to Nextlinux")
	return nil
}

func buildURL(nextlinuxDetails config.NextlinuxInfo) (string, error) {
	nextlinuxURL, err := url.Parse(nextlinuxDetails.URL)
	if err != nil {
		return "", err
	}

	nextlinuxURL.Path += ReportAPIPath

	return nextlinuxURL.String(), nil
}
