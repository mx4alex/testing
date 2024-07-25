package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

var token = "CorrectToken"

func TestFindUsers(t *testing.T) {
	testCases := []struct {
		name          string
		request       SearchRequest
		expectedUsers []User
	}{
		{
			name: "TestOrderIDAsc",
			request: SearchRequest{
				Limit:      2,
				OrderField: "ID",
				OrderBy:    OrderByAsc,
			},
			expectedUsers: []User{
				{
					ID:     0,
					Name:   "Boyd Wolf",
					Age:    22,
					About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
					Gender: "male",
				},
				{
					ID:     1,
					Name:   "Hilda Mayer",
					Age:    21,
					About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
					Gender: "female",
				},
			},
		},
		{
			name: "TestOrderIDDesc",
			request: SearchRequest{
				Limit:      2,
				OrderField: "ID",
				OrderBy:    OrderByDesc,
			},
			expectedUsers: []User{
				{
					ID:     34,
					Name:   "Kane Sharp",
					Age:    34,
					About:  "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n",
					Gender: "male",
				},
				{
					ID:     33,
					Name:   "Twila Snow",
					Age:    36,
					About:  "Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n",
					Gender: "female",
				},
			},
		},
		{
			name: "TestOrderNameAsc",
			request: SearchRequest{
				Limit:      1,
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			expectedUsers: []User{
				{
					ID:     15,
					Name:   "Allison Valdez",
					Age:    21,
					About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
					Gender: "male",
				},
			},
		},
		{
			name: "TestOrderNameDesc",
			request: SearchRequest{
				Limit:      1,
				OrderField: "Name",
				OrderBy:    OrderByDesc,
			},
			expectedUsers: []User{
				{
					ID:     13,
					Name:   "Whitley Davidson",
					Age:    40,
					About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
					Gender: "male",
				},
			},
		},
		{
			name: "TestOrderAgeAsc",
			request: SearchRequest{
				Limit:      1,
				OrderField: "Age",
				OrderBy:    OrderByAsc,
			},
			expectedUsers: []User{
				{
					ID:     1,
					Name:   "Hilda Mayer",
					Age:    21,
					About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
					Gender: "female",
				},
			},
		},
		{
			name: "TestOrderAgeDesc",
			request: SearchRequest{
				Limit:      1,
				OrderField: "Age",
				OrderBy:    OrderByDesc,
			},
			expectedUsers: []User{
				{
					ID:     32,
					Name:   "Christy Knapp",
					Age:    40,
					About:  "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n",
					Gender: "female",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(SearchServer))
			defer ts.Close()

			client := &SearchClient{
				AccessToken: token,
				URL:         ts.URL,
			}

			result, err := client.FindUsers(tc.request)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tc.expectedUsers, result.Users) {
				t.Errorf("Wrong result, expected %v, got %v", tc.expectedUsers, result.Users)
			}
		})
	}
}

func TestRequestsError(t *testing.T) {
	testCases := []struct {
		name          string
		request       SearchRequest
		expectedError error
	}{
		{
			name: "TestInvalidLimitError",
			request: SearchRequest{
				Limit: -1,
			},
			expectedError: fmt.Errorf("limit must be > 0"),
		},
		{
			name: "TestInvalidOffsetError",
			request: SearchRequest{
				Offset: -1,
			},
			expectedError: fmt.Errorf("offset must be > 0"),
		},
		{
			name: "TestInvalidOrderFieldError",
			request: SearchRequest{
				OrderField: "user",
			},
			expectedError: fmt.Errorf("OrderFeld user invalid"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(SearchServer))
			defer ts.Close()

			client := &SearchClient{
				AccessToken: token,
				URL:         ts.URL,
			}

			_, err := client.FindUsers(tc.request)
			if err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error, expected %v, got %v", tc.expectedError, err)
			}
		})
	}
}

func TestLimitOffsetUsers(t *testing.T) {
	testCases := []struct {
		name           string
		limit		   int 
		offset 		   int
		expectedCount  int
	}{
		{
			name:          "TestLimitUsers",
			limit:         50,
			offset:        0,
			expectedCount: 25,
		},
		{
			name:          "TestOffsetUsers",
			limit:         50,
			offset:        30,
			expectedCount: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(SearchServer))
			defer ts.Close()

			client := &SearchClient{
				AccessToken: token,
				URL:         ts.URL,
			}

			request := SearchRequest{
				Limit:  tc.limit,
				Offset: tc.offset,
			}

			result, err := client.FindUsers(request)
			if err != nil {
				t.Errorf("Unexpected error: %v", err.Error())
			}
			if len(result.Users) != tc.expectedCount {
				t.Errorf("Wrong number of users, expected %d, got: %d", tc.expectedCount, len(result.Users))
			}
		})
	}
}

func TestAtoiError(t *testing.T) {
	testCases := []struct {
		name           string
		paramGenerator func() url.Values
	}{
		{
			name: "TestAtoiLimitError",
			paramGenerator: func() url.Values {
				searcherParams := url.Values{}
				searcherParams.Add("limit", "one")
				return searcherParams
			},
		},
		{
			name: "TestAtoiOffsetError",
			paramGenerator: func() url.Values {
				searcherParams := url.Values{}
				searcherParams.Add("limit", "1")
				searcherParams.Add("offset", "one")
				return searcherParams
			},
		},
		{
			name: "TestAtoiOrderByError",
			paramGenerator: func() url.Values {
				searcherParams := url.Values{}
				searcherParams.Add("limit", "1")
				searcherParams.Add("offset", "1")
				searcherParams.Add("order_by", "one")
				return searcherParams
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(SearchServer))
			defer ts.Close()

			searcherParams := tc.paramGenerator()

			searcherReq, _ := http.NewRequest("GET", ts.URL + "?" + searcherParams.Encode(), nil) //nolint:errcheck
			searcherReq.Header.Add("AccessToken", token)
			client := &http.Client{}

			resp, err := client.Do(searcherReq)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("Unexpected status code, expected %v, got %v", http.StatusBadRequest, resp.StatusCode)
			}
		})
	}
}

func TestInvalidOrderByError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	expectedError :=  fmt.Errorf("unknown bad request error: OrderBy invalid")

	client := &SearchClient{
		AccessToken: token,
		URL:         ts.URL,
	}

	request := SearchRequest{
		OrderBy: 2,
	}

	_, err := client.FindUsers(request)
	if err.Error() != expectedError.Error() {
		t.Errorf("Unexpected error, expected %v, got %v", expectedError, err)
	}
}

func TestClientError(t *testing.T) {
	testCases := []struct {
		name    string
		server 	*httptest.Server
		expectedError error
	}{
		{
			name: "TestTimeoutError",
			server: httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(1 * time.Second)
				})),
			expectedError: fmt.Errorf("timeout for"),
		},
		{
			name: "TestCantUnpackError",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not SearchErrorResponse", http.StatusBadRequest)
			})),
			expectedError: fmt.Errorf("cant unpack error json"),
		},
		{
			name: "TestCantUnpackResultError",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not SearchErrorResponse", http.StatusOK)
			})),
			expectedError: fmt.Errorf("cant unpack result json"),
		},
		{
			name: "TestBadRequestError",
			server: httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					response := SearchErrorResponse{}
					response.Error = "badRequest"

					result, err := json.Marshal(response)
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					
					_, err = w.Write(result)
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
				})),
			expectedError: fmt.Errorf("unknown bad request error:"),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := tc.server
			defer ts.Close()

			client := &SearchClient{
				AccessToken: token,
				URL:         ts.URL,
			}

			_, err := client.FindUsers(SearchRequest{})
			if !strings.Contains(err.Error(), tc.expectedError.Error()) {
				t.Errorf("Unexpected error, expected %v, got %v", tc.expectedError, err)
			}
		})
	}
}

func TestInvalidTokenError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	expectedError := fmt.Errorf("bad AccessToken")

	client := &SearchClient{
		AccessToken: "invalidToken",
		URL:         ts.URL,
	}

	_, err := client.FindUsers(SearchRequest{})
	if err.Error() != expectedError.Error() {
		t.Errorf("Unexpected error, expected %v, got %v", expectedError, err)
	}
}

func TestInvalidUrlError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	expectedError := fmt.Errorf("unknown error")

	client := &SearchClient{
		AccessToken: token,
		URL:         "invalidUrl",
	}

	_, err := client.FindUsers(SearchRequest{})
	if !strings.Contains(err.Error(), expectedError.Error()) {
		t.Errorf("Unexpected error, expected %v, got %v", expectedError, err)
	}
}

func TestFileError(t *testing.T) {
	testCases := []struct {
		name 		  string
		filename      string
		expectedError error
	}{
		{
			name: 			"TestNotExistDatasetError",
			filename: 		"not_exist_dataset", 
			expectedError: 	fmt.Errorf("SearchServer fatal error"),
		},
		{
			name: 			"TestInvalidDatasetError",
			filename: 		"invalid_dataset.xml", 
			expectedError: 	fmt.Errorf("SearchServer fatal error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(SearchServer))
			defer ts.Close()

			client := &SearchClient{
				AccessToken: token,
				URL:         ts.URL,
			}

			filename = tc.filename
			_, err := client.FindUsers(SearchRequest{})

			if err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error, expected %v, got %v", tc.expectedError, err)
			}
		})
	}
}