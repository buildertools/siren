// Copyright 2017 Jeff Nickoloff "jeff@allingeek.com"
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

package clients

import (
	"net/http"
)

type ClientError interface {
	IsRetriable() bool
	Error() error
}

type RetriableError struct {
	E error
}

func (r RetriableError) IsRetriable() bool {
	return true
}
func (r RetriableError) Error() error {
	return r.E
}

type NonRetriableError struct {
	E error
}

func (n NonRetriableError) IsRetriable() bool {
	return false
}
func (n NonRetriableError) Error() error {
	return n.E
}

func WrapHttpResponseError(r *http.Response, err error) (*http.Response, ClientError) {
	if r == nil && err == nil {
		return nil, nil
	} else if err != nil {
		return r, RetriableError{E: err}
	}

	if r.StatusCode == http.StatusBadRequest ||
		r.StatusCode == http.StatusUnauthorized ||
		r.StatusCode == http.StatusPaymentRequired ||
		r.StatusCode == http.StatusForbidden ||
		r.StatusCode == http.StatusNotFound ||
		r.StatusCode == http.StatusMethodNotAllowed ||
		r.StatusCode == http.StatusNotAcceptable ||
		r.StatusCode == http.StatusProxyAuthRequired ||
		r.StatusCode == http.StatusRequestTimeout ||
		r.StatusCode == http.StatusConflict ||
		r.StatusCode == http.StatusGone ||
		r.StatusCode == http.StatusLengthRequired ||
		r.StatusCode == http.StatusPreconditionFailed ||
		r.StatusCode == http.StatusRequestEntityTooLarge ||
		r.StatusCode == http.StatusRequestURITooLong ||
		r.StatusCode == http.StatusUnsupportedMediaType ||
		r.StatusCode == http.StatusRequestedRangeNotSatisfiable ||
		r.StatusCode == http.StatusExpectationFailed ||
		r.StatusCode == http.StatusTeapot ||
		r.StatusCode == http.StatusUnprocessableEntity ||
		r.StatusCode == http.StatusLocked ||
		r.StatusCode == http.StatusFailedDependency ||
		r.StatusCode == http.StatusUpgradeRequired ||
		r.StatusCode == http.StatusPreconditionRequired ||
		r.StatusCode == http.StatusTooManyRequests ||
		r.StatusCode == http.StatusRequestHeaderFieldsTooLarge ||
		r.StatusCode == http.StatusUnavailableForLegalReasons {
		return r, NonRetriableError{E: err}
	} else if r.StatusCode == http.StatusInternalServerError ||
		r.StatusCode == http.StatusNotImplemented ||
		r.StatusCode == http.StatusBadGateway ||
		r.StatusCode == http.StatusServiceUnavailable ||
		r.StatusCode == http.StatusGatewayTimeout ||
		r.StatusCode == http.StatusHTTPVersionNotSupported ||
		r.StatusCode == http.StatusVariantAlsoNegotiates ||
		r.StatusCode == http.StatusInsufficientStorage ||
		r.StatusCode == http.StatusLoopDetected ||
		r.StatusCode == http.StatusNotExtended ||
		r.StatusCode == http.StatusNetworkAuthenticationRequired {
		return r, RetriableError{E: err}
	}
	return r, nil
}
