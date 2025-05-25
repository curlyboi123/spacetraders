package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/spacetraders", &SpaceTradersHandler{})
	mux.Handle("/spacetraders/", &SpaceTradersHandler{})

	// Run the server
	http.ListenAndServe(":8080", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

type SpaceTradersHandler struct{}

var (
	SpaceTradersBaseURL = "https://api.spacetraders.io/v2"
	MyAgentRe           = regexp.MustCompile(`^/spacetraders/myagent`)
	AgentToken          = os.Getenv("AGENT_TOKEN")
)

func (h *SpaceTradersHandler) GetAgentData(w http.ResponseWriter, r *http.Request) {
	log.Println("Making get agent data request")

	// TODO put token in secret storage, AWS parameter store/secrets manager?
	var bearer = "Bearer " + AgentToken

	req, err := http.NewRequest("GET", SpaceTradersBaseURL+"/my/agent", nil)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Println(sb)
}

func (h *SpaceTradersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && MyAgentRe.MatchString(r.URL.Path):
		h.GetAgentData(w, r)
		return
	default:
		return
	}
}
