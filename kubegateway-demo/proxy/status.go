// Copyright 2022 ByteDance and its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	StatusReasonNoReadyEndpoints         = "no_ready_endpoints"
	StatusReasonClusterNotBeingProxied   = "cluster_not_being_proxied"
	StatusReasonInvalidRequestContext    = "invalid_request_context"
	StatusReasonCircuitBreaker           = "circuit_breaker"
	StatusReasonRateLimited              = "rate_limited"
	StatusReasonInvalidEndpoint          = "invalid_endpoint"
	StatusReasonUpgradeAwareHandlerError = "upgrade_aware_handler_error"
	StatusReasonReverseProxyError        = "reverse_proxy_error"
)

func CaptureErrorReason(reason string) bool {
	switch reason {
	case StatusReasonUpgradeAwareHandlerError, StatusReasonReverseProxyError:
		return true
	}
	return false
}

// statusError is an object that can be converted into an metav1.Status
type statusError interface {
	Status() metav1.Status
}

// ErrorToProxyStatus converts an error to an metav1.Status object.
func ErrorToProxyStatus(err error) *metav1.Status {
	switch t := err.(type) {
	case statusError:
		status := t.Status()
		if len(status.Status) == 0 {
			status.Status = metav1.StatusFailure
		}
		switch status.Status {
		case metav1.StatusSuccess:
			if status.Code == 0 {
				status.Code = http.StatusOK
			}
		case metav1.StatusFailure:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		default:
			if status.Code == 0 {
				status.Code = http.StatusInternalServerError
			}
		}
		status.Kind = "Status"
		status.APIVersion = "v1"
		//TODO: check for invalid responses
		return &status
	default:
		return &metav1.Status{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Status",
				APIVersion: "v1",
			},
			Status:  metav1.StatusFailure,
			Code:    int32(http.StatusBadGateway),
			Reason:  "KubeGatewayInternalError",
			Message: err.Error(),
		}
	}
}
