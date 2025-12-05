// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vault

import (
	"context"
	"encoding/pem"
	"net/http"

	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

// PKICACertificateHandlerAPIVersion defines the version
const PKICACertificateHandlerAPIVersion = "v1"

type pkiCACertificateHandler struct {
	certificateIssuer CertificateIssuer
}

func (h *pkiCACertificateHandler) reply(ctx context.Context, cert string, err Error, w http.ResponseWriter) {
	var resp string
	log := core.GetLogger(ctx)
	if err == ErrorNone {
		resp = cert
	} else {
		resp = err.Error()
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(err.Code())
	_, errWrite := w.Write([]byte(resp))
	if errWrite != nil {
		log.Error(errWrite)
		http.Error(w, errWrite.Error(), http.StatusInternalServerError)
	}
}

// ServeHTTP implements /v1/pki/ca/*
func (h *pkiCACertificateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core.GetLogger(ctx)

	var cert string
	var err error

	switch r.URL.Path {
	case "/v1/pki/ca":
		cert, err = h.certificateIssuer.GetCACertificate(ctx)
		if err != nil {
			log.WithField("err", ErrorRequestCACertificate.String()).Errorf("failed to request PKI CA certificate: %s", err.Error())
			h.reply(ctx, "", ErrorRequestCACertificate, w)
			return
		}

		block, _ := pem.Decode([]byte(cert))
		if block == nil || block.Type != "CERTIFICATE" {
			log.WithField("err", ErrorDecodeCACertificate.String()).Errorf("failed to decode PKI CA certificate")
			h.reply(ctx, "", ErrorDecodeCACertificate, w)
			return
		}
		cert = string(block.Bytes)
	case "/v1/pki/ca/pem":
		cert, err = h.certificateIssuer.GetCACertificate(ctx)
		if err != nil {
			log.WithField("err", ErrorRequestCACertificate.String()).Errorf("failed to request PKI CA certificate: %s", err.Error())
			h.reply(ctx, "", ErrorRequestCACertificate, w)
			return
		}
	default:
		log.WithField("err", ErrorBadPKIRequest.String()).Errorf("invalid path")
		h.reply(ctx, "", ErrorBadPKIRequest, w)
		return
	}

	h.reply(ctx, cert, ErrorNone, w)
}
