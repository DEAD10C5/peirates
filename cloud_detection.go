package peirates

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var hc = &http.Client{Timeout: 300 * time.Millisecond}

type CloudProvider struct {
	Name              string
	URL               string
	HTTPMethod        string
	CustomHeader      string
	CustomHeaderValue string
	ResultString      string
}

func populateAndCheckCloudProviders() string {
	providers := []CloudProvider{
		{
			Name:              "AWS",
			URL:               "http://169.254.169.254/latest/",
			HTTPMethod:        "GET",
			CustomHeader:      "",
			CustomHeaderValue: "",
			ResultString:      "Amazon Web Services",
		},
		{
			Name:              "Azure",
			URL:               "http://169.254.169.254/metadata/v1/InstanceInfo",
			HTTPMethod:        "GET",
			CustomHeader:      "",
			CustomHeaderValue: "",
			ResultString:      "Microsoft Azure",
		},
		{
			Name:              "DigitalOcean",
			URL:               "http://169.254.169.254/metadata/v1.json",
			HTTPMethod:        "GET",
			CustomHeader:      "",
			CustomHeaderValue: "",
			ResultString:      "Microsoft Azure",
		},
		{
			Name:              "Google Cloud",
			URL:               "http://metadata.google.internal/computeMetadata/v1/instance/tags",
			HTTPMethod:        "GET",
			CustomHeader:      "Metadata-Flavor",
			CustomHeaderValue: "Google",
			ResultString:      "Google Compute Engine",
		},
		{
			Name:              "SoftLayer",
			URL:               "https://api.service.softlayer.com/rest/v3/SoftLayer_Resource_Metadata/UserMetadata.txt",
			HTTPMethod:        "GET",
			CustomHeader:      "",
			CustomHeaderValue: "",
			ResultString:      "SoftLayer",
		},
		{
			Name:              "Vultr",
			URL:               "http://169.254.169.254/v1.json",
			HTTPMethod:        "GET",
			CustomHeader:      "",
			CustomHeaderValue: "",
			ResultString:      "Vultr",
		},
	}

	for _, provider := range providers {
		fmt.Printf("Checking %s...\n", provider.Name)

		var response string

		var lines []HeaderLine
		if provider.CustomHeader != "" {
			line := HeaderLine{LHS: provider.CustomHeader, RHS: provider.CustomHeaderValue}
			lines = append(lines, line)
			response = GetRequest(provider.URL, lines, true)
		} else {
			response = GetRequest(provider.URL, nil, true)
		}
		if strings.Contains(response, provider.ResultString) {
			println("DEBUG: detected provider ", provider.Name)
			return provider.Name
		}
	}
	return ""
}

func detectContainer() string {
	b, err := ioutil.ReadFile("/proc/self/cgroup")
	if err != nil {
		return ""
	}

	fc := string(b)
	kube := strings.Contains(fc, "kube")
	container := strings.Contains(fc, "containerd")

	if kube {
		return "K8S Container"
	}

	if container {
		return "Container"
	}

	return ""
}

func detectOpenStack() string {
	if runtime.GOOS != "windows" {
		data, err := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor")
		if err != nil {
			return ""
		}
		if strings.Contains(string(data), "OpenStack Foundation") {
			return "OpenStack"
		}
		return ""
	}
	return ""
}
