package endpoint

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
)

var r = InitRouter()

func TestConcurrency(t *testing.T) {
	for i := 0; i < 500; i++ {
		for _, t := range tests {
			go testRequest(r, t.method, t.url, t.payload)
		}
	}
}

func BenchmarkPerformance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := tests[rand.Intn(len(tests))]
		go testRequest(r, t.method, t.url, t.payload)
	}

}

var tests = []struct {
	desc     string
	method   string
	url      string
	payload  string
	response int
	body     string
}{
	{desc: "List the dls", method: "GET", url: "/dls", response: 200, body: `{}`},
	{desc: "Add  dls", method: "POST", url: "/dls", payload: `{"A":["a1", "a2"], "B":["b1"]}`, response: 201, body: `[]`},
	{desc: "Add  more dls", method: "POST", url: "/dls", payload: `{"A":["a1", "a3"], "B":["b2"], "C":["c1"]}`, response: 201, body: `[]`},
	{desc: "Add  more dls", method: "POST", url: "/dls", payload: `{"a1":["a1a", "a1b"]}`, response: 201, body: `[]`},
	{desc: "List the dls", method: "GET", url: "/dls", response: 200, body: `{}`},
	{desc: "List members of A", method: "GET", url: "/members/A", response: 200, body: `{}`},
	{desc: "List parents of a1b", method: "GET", url: "/parents/a1b", response: 200, body: `{}`},
	{desc: "Delete A", method: "DELETE", url: "/dls/A", response: 200, body: `{}`},
	{method: "POST", url: "/dls", payload: `{"A":["a1", "a1", "a2"], "B":["b1"]}`, response: 201, body: `[]`},
	{method: "GET", url: "/members/NonExistent", response: 200, body: `{}`},
	{method: "GET", url: "/parents/NonExistent", response: 200, body: `{}`},
	{method: "DELETE", url: "/dls/NonExistent", response: 200, body: `{}`},
}

func Example() {

	for _, t := range tests {
		code, body := testRequest(r, t.method, t.url, t.payload)
		fmt.Printf("%s: %s %s %s ----> HTTP: %s(%d), Body: %s\n", t.desc, t.method, t.url, t.payload, http.StatusText(code), code, body)
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
