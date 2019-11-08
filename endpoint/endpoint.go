package endpoint

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

type Email string

type DLs map[Email][]Email

var dls DLs = make(map[Email][]Email)

func index(emails []Email, email Email) int {
	for i, e := range emails {
		if e == email {
			return i
		}
	}
	return -1
}

func testRequest(r *mux.Router, verb string, url string, payload string) (int, string /*, models.Workshop, []models.Workshop*/) {
	p := strings.NewReader(payload)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(verb, url, p)
	r.ServeHTTP(rr, req)
	return rr.Code, strings.TrimSpace(rr.Body.String())
}

func extractWorkshopFromPayload(r *http.Request) (DLs, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newDLs DLs
	e := json.Unmarshal(reqBody, &newDLs)
	if e != nil {
		return nil, errors.New("Bad Json")
	}
	return newDLs, nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func removeDuplicates(emails []Email) []Email {
	alreadyAdded := make(map[Email]bool)
	uniqueEmails := []Email{}
	for _, e := range emails {
		if !alreadyAdded[e] {
			uniqueEmails = append(uniqueEmails, e)
			alreadyAdded[e] = true
		}
	}
	return uniqueEmails
}

func DeleteDL(w http.ResponseWriter, r *http.Request) {
	email := Email(mux.Vars(r)["email"])
	delete(dls, email)
	for k, v := range dls {
		if i := index(v, email); i > -1 {
			dls[k] = append(dls[k][:i], dls[k][i:]...)
		}
	}
	respondWithJSON(w, http.StatusOK, dls)
}

func UpdateDLs(w http.ResponseWriter, r *http.Request) {
	newDLs, e := extractWorkshopFromPayload(r)
	if e != nil {
		respondWithError(w, http.StatusUnprocessableEntity, e.Error())
		return
	}
	for k, v := range newDLs {
		dls[k] = append(dls[k], v...)
		dls[k] = removeDuplicates(dls[k])
	}
	respondWithJSON(w, http.StatusOK, dls)
}

func GetDLs(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, dls)
}

func getMembers(emails []Email) []Email {
	members := emails
	for _, e := range emails {
		members = append(members, getMembers(dls[e])...)
	}
	return members
}

func GetMembers(w http.ResponseWriter, r *http.Request) {
	email := Email(mux.Vars(r)["email"])
	a := []Email{email}
	members := getMembers(a)

	leaves := []Email{}
	for _, m := range members {
		_, b := dls[m]
		if !b && email != m {
			leaves = append(leaves, m)
		}
	}

	respondWithJSON(w, http.StatusOK, leaves)
}

func getParents(email Email) []Email {
	parents := []Email{}
	for k, v := range dls {
		if index(v, email) > -1 {
			parents = append(parents, k)
			parents = append(parents, getParents(k)...)
		}
	}
	return parents
}

func GetParents(w http.ResponseWriter, r *http.Request) {
	email := Email(mux.Vars(r)["email"])
	parents := getParents(email)
	respondWithJSON(w, http.StatusOK, parents)
}

func InitRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/members/{email}", GetMembers).Methods("GET")
	router.HandleFunc("/parents/{email}", GetParents).Methods("GET")
	router.HandleFunc("/dls/{email}", DeleteDL).Methods("DELETE")
	router.HandleFunc("/dls", UpdateDLs).Methods("POST")
	router.HandleFunc("/dls", GetDLs).Methods("GET")
	return router
}
