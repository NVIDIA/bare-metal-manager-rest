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
	"encoding/json"
	"net/http"

	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

type pkiCloudCertificateHandler struct {
	certificateIssuer CertificateIssuer
}

func (h *pkiCloudCertificateHandler) reply(ctx context.Context, cert, key string, err Error, w http.ResponseWriter) {
	log := core.GetLogger(ctx)

	resp := &CertificateResponse{}
	if err == ErrorNone {
		resp.Key = key
		resp.Certificate = cert
	}

	respBytes, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		log.WithField("err", ErrorMarshalJSON.String()).Errorf("Failed to json.Marshal ClientCertificateResponse %+v, err: %s", resp, marshalErr.Error())
		http.Error(w, ErrorMarshalJSON.Error(), ErrorMarshalJSON.Code())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err.Code())
	_, errWrite := w.Write(respBytes)
	if errWrite != nil {
		log.Error(errWrite)
		http.Error(w, errWrite.Error(), http.StatusInternalServerError)
		return
	}
}

// ServeHTTP implements /v1/pki/cloud-cert
func (h *pkiCloudCertificateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core.GetLogger(ctx)

	req := &CertificateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.WithField("err", ErrorParseRequest.String()).Errorf("failed to parse request body as CertificateRequest: %s", err.Error())
		h.reply(ctx, "", "", ErrorParseRequest, w)
		return
	}

	cert, privKey, err := h.certificateIssuer.NewCertificate(ctx, req)

	if err != nil {
		log.WithField("err", ErrorGetCertificate.String()).Errorf("failed certificateIssuer.NewCertificate, err: %s", err.Error())
		h.reply(ctx, "", "", ErrorGetCertificate, w)
		return
	}

	h.reply(ctx, cert, privKey, ErrorNone, w)
}
