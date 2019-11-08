package endpoint

import (
	"fmt"
	"net/http"
)

func Example_happyPaths() {
	r := InitRouter()
	var tests = []struct {
		method   string
		url      string
		payload  string
		code     int
		response string
	}{
		{method: "GET", url: "/dls", code: 200, response: `{}`},
		{method: "POST", url: "/dls", payload: `{"A":["a1", "a2"], "B":["b1"]}`, code: 201, response: `[]`},
		{method: "POST", url: "/dls", payload: `{"A":["a1", "a3"], "B":["b2"], "C":["c1"]}`, code: 201, response: `[]`},
		{method: "POST", url: "/dls", payload: `{"a1":["a1a", "a1b"]}`, code: 201, response: `[]`},
		{method: "GET", url: "/dls", code: 200, response: `{}`},
		{method: "GET", url: "/members/A", code: 200, response: `{}`},
		{method: "GET", url: "/parents/a1b", code: 200, response: `{}`},
		{method: "DELETE", url: "/dls/A", code: 200, response: `{}`},
	}

	for _, t := range tests {
		code, body := testRequest(r, t.method, t.url, t.payload)
		fmt.Printf("%s %s %s\n    -> HTTP Status: %s(%d), Body: %s\n", t.url, t.payload, t.method, http.StatusText(code), code, body)
	}

	//Output:
	// 	/dls  GET
	//     -> HTTP Status: OK(200), Body: {}
	// /dls {"A":["a1", "a2"], "B":["b1"]} POST
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2"],"B":["b1"]}
	// /dls {"A":["a1", "a3"], "B":["b2"], "C":["c1"]} POST
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2","a3"],"B":["b1","b2"],"C":["c1"]}
	// /dls {"a1":["a1a", "a1b"]} POST
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2","a3"],"B":["b1","b2"],"C":["c1"],"a1":["a1a","a1b"]}
	// /dls  GET
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2","a3"],"B":["b1","b2"],"C":["c1"],"a1":["a1a","a1b"]}
	// /members/A  GET
	//     -> HTTP Status: OK(200), Body: ["a2","a3","a1a","a1b"]
	// /parents/a1b  GET
	//     -> HTTP Status: OK(200), Body: ["a1","A"]
	// /dls/A  DELETE
	//     -> HTTP Status: OK(200), Body: {"B":["b1","b2"],"C":["c1"],"a1":["a1a","a1b"]}

}

func Example_unhappyPaths() {
	r := InitRouter()
	var tests = []struct {
		method   string
		url      string
		payload  string
		code     int
		response string
	}{
		{method: "POST", url: "/dls", payload: `{"A":["a1", "a2"], "B":["b1"]}`, code: 201, response: `[]`},
		{method: "GET", url: "/members/NonExistent", code: 200, response: `{}`},
		{method: "GET", url: "/parents/NonExistent", code: 200, response: `{}`},
		{method: "DELETE", url: "/dls/NonExistent", code: 200, response: `{}`},
	}

	for _, t := range tests {
		code, body := testRequest(r, t.method, t.url, t.payload)
		fmt.Printf("%s %s %s\n    -> HTTP Status: %s(%d), Body: %s\n", t.url, t.payload, t.method, http.StatusText(code), code, body)
	}

	//Output:
	// 	/dlsss  GET
	//     -> HTTP Status: Not Found(404), Body: 404 page not found
	// /dls {"A":["a1", "a2"], "B":["b1"]} POST
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2"],"B":["b1"]}
	// /dls {"A":["a1", "a3"], "B":["b2"], "C":["c1"]} POST
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2","a3"],"B":["b1","b2"],"C":["c1"]}
	// /dls {"a1":["a1a", "a1b"]} POST
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2","a3"],"B":["b1","b2"],"C":["c1"],"a1":["a1a","a1b"]}
	// /dls  GET
	//     -> HTTP Status: OK(200), Body: {"A":["a1","a2","a3"],"B":["b1","b2"],"C":["c1"],"a1":["a1a","a1b"]}
	// /members/A  GET
	//     -> HTTP Status: OK(200), Body: ["a2","a3","a1a","a1b"]
	// /parents/a1b  GET
	//     -> HTTP Status: OK(200), Body: ["a1","A"]
	// /dls/A  DELETE
	//     -> HTTP Status: OK(200), Body: {"B":["b1","b2"],"C":["c1"],"a1":["a1a","a1b"]}

}
