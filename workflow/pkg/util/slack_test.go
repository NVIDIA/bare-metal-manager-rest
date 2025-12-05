// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlackClient_sendHttpRequest(t *testing.T) {
	type fields struct {
		WebHookUrl string
		TimeOut    time.Duration
	}
	type args struct {
		slackRequest SlackMessage
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test sending http request with Slack message",
			fields: fields{
				WebHookUrl: testServer.URL,
				TimeOut:    5 * time.Second,
			},
			args: args{
				slackRequest: SlackMessage{
					IconEmoji: ":ghost:",
					Text:      "This is a test message",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := SlackClient{
				webHookUrl: tt.fields.WebHookUrl,
				timeOut:    tt.fields.TimeOut,
			}

			err := sc.sendHttpRequest(tt.args.slackRequest)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
