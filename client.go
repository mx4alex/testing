package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	orderAsc = iota
	orderDesc
)

var (
	errTest = errors.New("testing")
	client  = &http.Client{Timeout: time.Second}
)

type User struct {
	ID     int
	Name   string
	Age    int
	About  string
	Gender string
}

type SearchResponse struct {
	Users    []User
	NextPage bool
}

type SearchErrorResponse struct {
	Error string
}

const (
	OrderByAsc  = 1
	OrderByAsIs = 0
	OrderByDesc = -1

	ErrorBadOrderField = `OrderField invalid`
)

type SearchRequest struct {
	Limit      int
	Offset     int
	Query      string
	OrderField string
	OrderBy    int
}

type SearchClient struct {
	AccessToken string
	URL         string
}

func (srv *SearchClient) FindUsers(req SearchRequest) (*SearchResponse, error) {

	searcherParams := url.Values{}

	if req.Limit < 0 {
		return nil, fmt.Errorf("limit must be > 0")
	}
	if req.Limit > 25 {
		req.Limit = 25
	}
	if req.Offset < 0 {
		return nil, fmt.Errorf("offset must be > 0")
	}

	req.Limit++

	searcherParams.Add("limit", strconv.Itoa(req.Limit))
	searcherParams.Add("offset", strconv.Itoa(req.Offset))
	searcherParams.Add("query", req.Query)
	searcherParams.Add("order_field", req.OrderField)
	searcherParams.Add("order_by", strconv.Itoa(req.OrderBy))

	searcherReq, _ := http.NewRequest("GET", srv.URL+"?"+searcherParams.Encode(), nil)
	searcherReq.Header.Add("AccessToken", srv.AccessToken)

	resp, err := client.Do(searcherReq)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return nil, fmt.Errorf("timeout for %s", searcherParams.Encode())
		}
		return nil, fmt.Errorf("unknown error %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("bad AccessToken")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("SearchServer fatal error")
	case http.StatusBadRequest:
		errResp := SearchErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, fmt.Errorf("cant unpack error json: %s", err)
		}
		if errResp.Error == ErrorBadOrderField {
			return nil, fmt.Errorf("OrderFeld %s invalid", req.OrderField)
		}
		return nil, fmt.Errorf("unknown bad request error: %s", errResp.Error)
	}

	data := []User{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("cant unpack result json: %s", err)
	}

	result := SearchResponse{}
	if len(data) == req.Limit {
		result.NextPage = true
		result.Users = data[0 : len(data)-1]
	} else {
		result.Users = data[0:]
	}

	return &result, err
}
