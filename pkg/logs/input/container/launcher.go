// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package container

import (
	"time"

	"github.com/DataDog/datadog-agent/pkg/logs/auditor"
	"github.com/DataDog/datadog-agent/pkg/logs/config"
	"github.com/DataDog/datadog-agent/pkg/logs/input/docker"
	"github.com/DataDog/datadog-agent/pkg/logs/input/kubernetes"
	"github.com/DataDog/datadog-agent/pkg/logs/pipeline"
	"github.com/DataDog/datadog-agent/pkg/logs/restart"
	"github.com/DataDog/datadog-agent/pkg/logs/service"

	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// NewLauncher returns a new container launcher depending on the environment.
// By default returns a docker launcher if the docker socket is mounted and fallback to
// a kubernetes launcher if '/var/log/pods' is mounted ; this behaviour is reversed when
// kubernetesCollectFromFiles is enabled.
// If dockerCollectFromFiles is enabled the docker launcher will first attempt to tail
// containers from file instead of the docker socket if '/var/lib/docker/containers'
// is mounted.
// If none of those volumes are mounted, returns a lazy docker launcher with a retrier to handle the cases
// where docker is started after the agent.
// dockerReadTimeout is a configurable read timeout for the docker client.
func NewLauncher(collectAll bool, kubernetesCollectFromFiles bool, dockerCollectFromFiles bool, dockerReadTimeout time.Duration, sources *config.LogSources, services *service.Services, pipelineProvider pipeline.Provider, registry auditor.Registry) restart.Restartable {
	var (
		launcher restart.Restartable
		err      error
	)

	if kubernetesCollectFromFiles {
		launcher, err = kubernetes.NewLauncher(sources, services, collectAll)
		if err == nil {
			log.Info("Kubernetes launcher initialized")
			return launcher
		}
		log.Infof("Could not setup the kubernetes launcher: %v", err)

		launcher, err = docker.NewLauncher(dockerReadTimeout, sources, services, pipelineProvider, registry, false, dockerCollectFromFiles)
		if err == nil {
			log.Info("Docker launcher initialized")
			return launcher
		}
		log.Infof("Could not setup the docker launcher: %v", err)
	} else {
		launcher, err = docker.NewLauncher(dockerReadTimeout, sources, services, pipelineProvider, registry, false, dockerCollectFromFiles)
		if err == nil {
			log.Info("Docker launcher initialized")
			return launcher
		}
		log.Infof("Could not setup the docker launcher: %v", err)

		launcher, err = kubernetes.NewLauncher(sources, services, collectAll)
		if err == nil {
			log.Info("Kubernetes launcher initialized")
			return launcher
		}
		log.Infof("Could not setup the kubernetes launcher: %v", err)
	}

	launcher, err = docker.NewLauncher(dockerReadTimeout, sources, services, pipelineProvider, registry, true, dockerCollectFromFiles)
	if err != nil {
		log.Warnf("Could not setup the docker launcher: %v. Will not be able to collect container logs", err)
		return NewNoopLauncher()
	}

	log.Infof("Container logs won't be collected unless a docker daemon is eventually started")

	return launcher
}
