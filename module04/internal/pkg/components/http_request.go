package components

import (
	"fmt"
	"net/http"
	dlog2 "self/internal/pkg/dlog"
	"time"
)

type PageResponse struct {
	Data        interface{} `json:"data"`
	Count       int64       `json:"count"`
	CurrentPage int64       `json:"current_page"`
}

type Nhttp struct{}

//统一GET请求
func (s Nhttp) DoHttpGet(url string, header map[string]string, timeoutSec int) ([]byte, error) {
	startTime := time.Now()
	request := NewRequest().Timeout(time.Duration(timeoutSec) * time.Second)
	resp, body, errs := request.Get(url).Sets(header).End()
	cost := time.Since(startTime).Seconds() * 1e3

	httpStatusCode := 0

	if len(errs) > 0 {
		msg := request.FormatErrors(errs)
		errmsg := fmt.Sprintf("_http_request_error|url=%s|method=%s|errmsg=%s|func=%s|stage=%s|code=%d|cost=%.3fms",
			url, "get", msg, "DoHttpGet", "request", httpStatusCode, cost)
		dlog2.Warning(errmsg)
		return nil, fmt.Errorf(msg)
	}

	httpStatusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("http code %d", resp.StatusCode)
		errmsg := fmt.Sprintf("_http_request_error|url=%s|method=%s|errmsg=%s|func=%s|stage=%s|code=%d|cost=%.3fms",
			url, "get", msg, "DoHttpGet", "checkstatus", httpStatusCode, cost)
		dlog2.Warning(errmsg)
		if body != "" {
			return nil, fmt.Errorf(body)
		} else {
			return nil, fmt.Errorf(errmsg)
		}
	}

	dlog2.Debugf("_http_request_ok|url=%s|method=%s|code=%d|cost=%.3fms", url, "get", httpStatusCode, cost)
	return []byte(body), nil
}

//统一POST请求
func (s Nhttp) DoHttpPost(url, forceType string, data interface{}, header map[string]string, timeoutSec int) ([]byte, error) {
	startTime := time.Now()
	request := NewRequest().Timeout(time.Duration(timeoutSec) * time.Second)
	resp, body, errs := request.Post(url, forceType).Sets(header).Send(data).End()
	cost := time.Since(startTime).Seconds() * 1e3

	httpStatusCode := 0

	if len(errs) > 0 {
		msg := request.FormatErrors(errs)
		errmsg := fmt.Sprintf("_http_request_error|url=%s|method=%s|errmsg=%s|func=%s|stage=%s|code=%d|cost=%.3fms",
			url, "post", msg, "DoHttpPost", "request", httpStatusCode, cost)
		dlog2.Warning(errmsg)
		return nil, fmt.Errorf(msg)
	}

	httpStatusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("http code %d", resp.StatusCode)
		errmsg := fmt.Sprintf("_http_request_error|url=%s|method=%s|errmsg=%s|func=%s|stage=%s|code=%d|cost=%.3fms",
			url, "post", msg, "DoHttpPost", "checkstatus", httpStatusCode, cost)
		dlog2.Warning(errmsg)
		if body != "" {
			return nil, fmt.Errorf(body)
		} else {
			return nil, fmt.Errorf(errmsg)
		}
	}

	//dlog.Debugf("_http_request_ok|url=%s|method=%s|code=%d|cost=%.3fms", url, "post", httpStatusCode, cost)
	return []byte(body), nil
}

//设置http请求的header
func (s Nhttp) SetHeader(header map[string]string) map[string]string {
	if nil == header {
		header = make(map[string]string)
	}
	header["Accept-Encoding"] = "utf8"
	header["Accept"] = "*/*"
	return header
}
