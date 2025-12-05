// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
)

var (
	DefaultTemporalCAPath     = "/var/secrets/temporal/certs/ca"
	DefaultTemporalClientPath = "/var/secrets/temporal/certs/client"
	DefaultTemporalServer     = "site.server.temporal.nvidia.com"
	DefualtTemporalID         = "elektractl-proxy"
)

type TemporalClient struct {
	ID                 string
	Publisher          client.Client
	Subscriber         client.Client
	PublishWflowOpts   *client.StartWorkflowOptions
	SubscribeWflowOpts *client.StartWorkflowOptions
	// Worker started with Publish client
	Worker worker.Worker
}

func createTemporalClient() (*TemporalClient, error) {

	// 1. Check temporal parameters from cli env file
	// 2. Check whether we need to use TLS to connect to temporal
	// 3. Check whether we have already registered as a client
	// 4. Start both Publisher & Subscriber workflows
	InitLog()
	Elog.Info().Msg("checking temporal parameters")
	var tclient *TemporalClient
	tclient = &TemporalClient{
		ID: DefualtTemporalID,
	}
	var subscribens, publishns string
	var publishClientOptions client.ConnectionOptions
	var subscribeClientOptions client.ConnectionOptions

	subscribens = DefaultNamespace
	if ns, isPresent := os.LookupEnv("TEMPORAL_SUBSCRIBE_NAMESPACE"); isPresent {
		subscribens = ns
	}
	publishns = DefaultNamespace
	if ns, isPresent := os.LookupEnv("TEMPORAL_PUBLISH_NAMESPACE"); isPresent {
		publishns = ns
	}

	if mcq := os.Getenv("TEMPORAL_SUBSCRIBE_QUEUE"); mcq != "" {
		TemporalSubscribeQueue = mcq
	}

	if mcq := os.Getenv("TEMPORAL_PUBLISH_QUEUE"); mcq != "" {
		TemporalPublishQueue = mcq
	}
	tclient.SubscribeWflowOpts = &client.StartWorkflowOptions{
		ID:        tclient.ID,
		TaskQueue: TemporalSubscribeQueue,
	}
	tclient.PublishWflowOpts = &client.StartWorkflowOptions{
		ID:        tclient.ID,
		TaskQueue: TemporalPublishQueue,
	}

	var clientInterceptors []interceptor.ClientInterceptor
	var workerInterceptors []interceptor.WorkerInterceptor

	// 2. Check whether we need to use TLS to connect to temporal

	if tlsFlag, isPresent := os.LookupEnv("ENABLE_TLS"); isPresent {
		if strings.ToLower(tlsFlag) == "true" {
			Elog.Info().Msg("TemporalClient: Creating proxy Temporal client with tls")

			// Proceed with defualtpath
			TemporalClientCertPath := DefaultTemporalClientPath
			TemporalCACertPath := DefaultTemporalCAPath
			// Check if we need to load TLS from a custom path
			if _, isPresent := os.LookupEnv("TEMPORAL_CA_PATH"); isPresent {
				// update the paths accordingly
			}
			if _, isPresent := os.LookupEnv("TEMPORAL_CLIENT_PATH"); isPresent {
				// update the paths accordingly
			}

			// Load client cert
			clientcert, err := tls.LoadX509KeyPair(fmt.Sprintf("%v/%v", TemporalClientCertPath, "tls.crt"),
				fmt.Sprintf("%v/%v", TemporalClientCertPath, "tls.key"))
			if err != nil {
				Elog.Panic().Str("path", TemporalClientCertPath).Msg("TemporalClient: Unable to read client certificates, exiting...")
				return tclient, err
			}
			// Load CA cert
			caCert, err := os.ReadFile(fmt.Sprintf("%v/%v", TemporalCACertPath, "ca.crt"))
			if err != nil {
				Elog.Panic().Msg("TemporalClient: Unable to read server certificates, exiting...")
				return tclient, err

			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			// provide tls cert option for subcribing workflows
			subscribeClientOptions = client.ConnectionOptions{
				TLS: &tls.Config{
					Certificates: []tls.Certificate{clientcert},
					ServerName:   DefaultTemporalServer,
					RootCAs:      caCertPool,
				},
				KeepAliveTime:    10 * time.Second,
				KeepAliveTimeout: 20 * time.Second,
			}

			// provide tls cert option for publishing workflows
			publishClientOptions = client.ConnectionOptions{
				TLS: &tls.Config{
					Certificates: []tls.Certificate{clientcert},
					ServerName:   DefaultTemporalServer,
					RootCAs:      caCertPool,
				},
				KeepAliveTime:    10 * time.Second,
				KeepAliveTimeout: 20 * time.Second,
			}

		}
	} else {
		// no TLS - this should be development mode for the cli
		Elog.Info().Msg("proceeding with no tls for temporal connection")
	}

	var err error

	temporalHost, isPresent := os.LookupEnv("TEMPORAL_HOST")
	if !isPresent {
		// We cannot proceed
		Elog.Panic().Msg("TemporalClient: cannot find temporal host to connect, exiting...")
	}
	temporalPort, isPresent := os.LookupEnv("TEMPORAL_PORT")
	if !isPresent {
		// We cannot proceed
		Elog.Panic().Msg("TemporalClient: cannot find temporal host to connect, exiting...")
	}

	// Create Publisher Temporal Client
	Elog.Info().Str("Publish NS", publishns).Msgf("TemporalClient: Connecting to Host %v, Port %v", temporalHost, temporalPort)
	publisherOptions := client.Options{
		HostPort:          fmt.Sprintf("%s:%s", temporalHost, temporalPort),
		Namespace:         publishns,
		ConnectionOptions: publishClientOptions,
		DataConverter: converter.NewCompositeDataConverter(
			converter.NewNilPayloadConverter(),
			converter.NewByteSlicePayloadConverter(),
			converter.NewProtoJSONPayloadConverterWithOptions(converter.ProtoJSONPayloadConverterOptions{
				AllowUnknownFields: true,
			}),
			converter.NewProtoPayloadConverter(),
			converter.NewJSONPayloadConverter(),
		),
		Interceptors: clientInterceptors,
	}

	Publisher, err := client.NewLazyClient(publisherOptions)

	if err != nil {
		Elog.Panic().Msg("TemporalClient: Failed to create Temporal client")
		return tclient, err
	}

	tclient.Publisher = Publisher

	// Create Subscriber Temporal Client
	Elog.Info().Str("Subscribe NS", subscribens).Msgf("TemporalClient: Connecting to Host %v, Port %v", temporalHost, temporalPort)

	subscriberOptions := client.Options{
		HostPort:          fmt.Sprintf("%s:%s", temporalHost, temporalPort),
		Namespace:         subscribens,
		ConnectionOptions: subscribeClientOptions,
		DataConverter: converter.NewCompositeDataConverter(
			converter.NewNilPayloadConverter(),
			converter.NewByteSlicePayloadConverter(),
			converter.NewProtoJSONPayloadConverterWithOptions(converter.ProtoJSONPayloadConverterOptions{
				AllowUnknownFields: true,
			}),
			converter.NewProtoPayloadConverter(),
			converter.NewJSONPayloadConverter(),
		),
		Interceptors: clientInterceptors,
	}

	Subscriber, err := client.NewLazyClient(subscriberOptions)

	if err != nil {
		Elog.Panic().Msg("TemporalClient: Failed to create Temporal client")
		return tclient, err
	}

	tclient.Subscriber = Subscriber

	// Create Temporal Worker
	tclient.Worker = worker.New(tclient.Publisher, TemporalPublishQueue, worker.Options{
		Interceptors:                     workerInterceptors,
		MaxConcurrentActivityTaskPollers: 10,
		MaxConcurrentWorkflowTaskPollers: 10,
	})

	// Start listening to the Task Queue
	Elog.Info().Msg("TemporalClient: Starting Temporal Worker")
	err = tclient.Worker.Start()
	if err != nil {
		Elog.Error().Msg("TemporalClient: Failed to start Temporal Worker")
		return tclient, err
	}
	return tclient, err
}
