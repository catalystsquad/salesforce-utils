package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
)

type SoqlResponse struct {
	Done           bool
	TotalSize      int
	Records        []interface{}
	NextRecordsUrl string
}

func (s *SalesforceUtils) ExecuteSoqlQuery(query string) (SoqlResponse, error) {
	uri := s.getQueryUrl(query, s.getSoqlUrl())
	return s.doSoqlQuery(uri)
}

func (s *SalesforceUtils) ExecuteSoqlQueryAll(query string) (SoqlResponse, error) {
	uri := s.getQueryUrl(query, s.getSoqlQueryAllUrl())
	return s.doSoqlQuery(uri)
}

func (s *SalesforceUtils) doSoqlQuery(uri string) (response SoqlResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodGet)
	body, statusCode, deferredFunc, requestErr := s.sendRequest(req)
	defer deferredFunc()
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusOK {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

func (s *SalesforceUtils) GetNextRecords(nextRecordsUrl string) (response SoqlResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getNextRecordsUrl(nextRecordsUrl)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodGet)
	body, statusCode, deferredFunc, requestErr := s.sendRequest(req)
	defer deferredFunc()
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusOK {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

// getSoqlUrl gets a formatted url to the soql endpoint
func (s *SalesforceUtils) getSoqlUrl() string {
	return fmt.Sprintf("%s/services/data/v%s/query", s.Config.BaseUrl, s.Config.ApiVersion)
}

// getSoqlQueryAllUrl gets a formatted url to the queryall soql endpoint
func (s *SalesforceUtils) getSoqlQueryAllUrl() string {
	return fmt.Sprintf("%s/services/data/v%s/queryAll", s.Config.BaseUrl, s.Config.ApiVersion)
}

// getQueryUrl gets a formatted url to the soql endpoint with the formatted query string included
func (s *SalesforceUtils) getQueryUrl(query string, path string) string {
	// url encode the query
	params := url.Values{}
	params.Add("q", query)
	return fmt.Sprintf("%s?%s", path, params.Encode())
}

// getSoqlUrl gets a formatted url to the soql endpoint
func (s *SalesforceUtils) getNextRecordsUrl(nextRecordsUrl string) string {
	return fmt.Sprintf("%s/%s", s.Config.BaseUrl, nextRecordsUrl)
}
