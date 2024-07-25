package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var filename = "dataset.xml"

type Person struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
}

type Data struct {
	Name xml.Name `xml:"root"`
	Rows []Person `xml:"row"`
}

func sortUsers(orderField string, orderBy int, users []User) bool {
	var less func(i, j int) bool
	switch {
	case orderField == "ID":
		if orderBy == OrderByAsc {
			less = func(i, j int) bool {
				return users[i].ID < users[j].ID
			}
		} else if orderBy == OrderByDesc {
			less = func(i, j int) bool {
				return users[i].ID > users[j].ID
			}
		}
	case orderField == "Name" || orderField == "":
		if orderBy == OrderByAsc {
			less = func(i, j int) bool {
				return users[i].Name < users[j].Name
			}
		} else if orderBy == OrderByDesc {
			less = func(i, j int) bool {
				return users[i].Name > users[j].Name
			}
		}
	case orderField == "Age":
		if orderBy == OrderByAsc {
			less = func(i, j int) bool {
				return users[i].Age < users[j].Age
			}
		} else if orderBy == OrderByDesc {
			less = func(i, j int) bool {
				return users[i].Age > users[j].Age
			}
		}
	default:
		return false
	}

	if less != nil {
		sort.Slice(users, less)
	}

	return true
}

func writeError(w http.ResponseWriter, textError string) {
	w.WriteHeader(http.StatusBadRequest)
	response := SearchErrorResponse{}
	response.Error = textError

	result, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("AccessToken")

	if accessToken != "CorrectToken" {
		http.Error(w, "Uncorrect AccessToken", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderBy, err := strconv.Atoi(q.Get("order_by"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if orderBy != OrderByAsc && orderBy != OrderByAsIs && orderBy != OrderByDesc {
		writeError(w, "OrderBy invalid")
		return
	}

	orderField := q.Get("order_field")
	query := q.Get("query")

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := new(Data)
	err = xml.Unmarshal(input, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var builder strings.Builder
	var users []User
	var user User

	for _, person := range data.Rows {
		builder.WriteString(person.FirstName)
		builder.WriteString(" ")
		builder.WriteString(person.LastName)

		if strings.Contains(builder.String(), query) || strings.Contains(person.About, query) {
			user = User{
				ID:     person.ID,
				Name:   builder.String(),
				Age:    person.Age,
				About:  person.About,
				Gender: person.Gender,
			}

			users = append(users, user)
		}

		builder.Reset()
	}

	isSorted := sortUsers(orderField, orderBy, users)
	if !isSorted {
		writeError(w, ErrorBadOrderField)
		return
	}

	end := offset + limit
	if end > len(users) {
		users = users[offset:]
	} else {
		users = users[offset:end]
	}

	var result []byte
	result, err = json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
