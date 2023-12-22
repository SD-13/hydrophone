/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"flag"
	"fmt"
	"github.com/dims/hydrophone/pkg/client"
	"k8s.io/client-go/rest"
	"log"
	"os"
)

// ArgConfig stores the argument passed when running the program
type ArgConfig struct {
	// Focus set the E2E_FOCUS env var to run a specific test
	// e.g. - sig-auth, sig-apps
	Focus string

	// Skip set the E2E_SKIP env var to skip specified tests
	Skip string

	// ConformanceImage let's people use the conformance container image of their own choice
	// Get the list of images from https://console.cloud.google.com/gcr/images/k8s-artifacts-prod/us/conformance
	// default registry.k8s.io/conformance:v1.28.0
	ConformanceImage string

	// BusyboxImage lets folks use an appropriate busybox image from their own registry
	BusyboxImage string

	// Kubeconfig is the path to the kubeconfig file
	Kubeconfig string

	// Parallel sets the E2E_PARALLEL env var for tests
	Parallel int

	// Verbosity sets the E2E_VERBOSITY env var for tests
	Verbosity int

	// OutputDir is where the e2e.log and junit_01.xml is saved
	OutputDir string
}

func InitArgs() (*ArgConfig, error) {
	var cfg ArgConfig

	outputDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.Focus, "focus", "", "focus runs a specific e2e test. e.g. - sig-auth. allows regular expressions.")
	flag.StringVar(&cfg.Skip, "skip", "", "skip specific tests. allows regular expressions.")
	flag.StringVar(&cfg.ConformanceImage, "conformance-image", containerImage,
		"image let's you select your conformance container image of your choice.")
	flag.StringVar(&cfg.BusyboxImage, "busybox-image", busyboxImage,
		"image let's you select an alternate busybox container image.")
	flag.StringVar(&cfg.Kubeconfig, "kubeconfig", "", "path to the kubeconfig file.")
	flag.IntVar(&cfg.Parallel, "parallel", 1, "number of parallel threads in test framework.")
	flag.IntVar(&cfg.Verbosity, "verbosity", 4, "verbosity of test framework.")
	flag.StringVar(&cfg.OutputDir, "output-dir", outputDir, "directory for logs.")

	flag.Parse()

	if cfg.Focus == "" {
		return nil, fmt.Errorf("missing --focus argument (use '[Conformance]' to run all conformance tests)")
	}

	return &cfg, nil
}

func ValidateArgs(err error, client *client.Client, config *rest.Config, cfg *ArgConfig) {
	serverVersion, err := client.ClientSet.ServerVersion()
	if err != nil {
		log.Fatal("Error fetching server version: ", err)
	}
	log.Printf("API endpoint : %s", config.Host)
	log.Printf("Server version : %#v", *serverVersion)
	log.Printf("Running tests : '%s'", cfg.Focus)
	if cfg.Skip != "" {
		log.Printf("Skipping tests : '%s'", cfg.Skip)
	}
	log.Printf("Using conformance image : '%s'", cfg.ConformanceImage)
	log.Printf("Using busybox image : '%s'", cfg.BusyboxImage)
	log.Printf("Test framework will start '%d' threads and use verbosity '%d'",
		cfg.Parallel, cfg.Verbosity)

	if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			log.Fatalf("Error creating output directory [%s] : %v", cfg.OutputDir, err)
		}
	}
}
