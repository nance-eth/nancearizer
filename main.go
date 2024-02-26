package main

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/daulet/tokenizers"
	"github.com/joho/godotenv"
)

var (
	s    *discordgo.Session
	tk   *tokenizers.Tokenizer
	port string
)

func init() {
	_, err := os.Stat(".env")
	if !os.IsNotExist(err) {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v\n", err)
		}
	}

	for _, v := range []string{"DISCORD_TOKEN", "OPENAI_API_KEY", "OPENAI_API_URL"} {
		if os.Getenv(v) == "" {
			log.Fatalf("No %s found in env\n", v)
		}
	}

	apiUrl = os.Getenv("OPENAI_API_URL")
	apiKey = os.Getenv("OPENAI_API_KEY")

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating Discord session: %v\n", err)
	}

	tk, err = tokenizers.FromFile("tokenizer.json")
	if err != nil {
		log.Fatalf("Error loading tokenizer: %v\n", err)
	}
}

func main() {
	defer tk.Close()

	http.HandleFunc("GET /proposal/{space}/{id}", summarizeProposal)
	http.HandleFunc("POST /proposal", summarizeProposal)
	http.HandleFunc("GET /thread/{space}/{id}", summarizeThread)
	http.HandleFunc("POST /thread", summarizeThread)

	log.Println("Listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func summarizeProposal(w http.ResponseWriter, req *http.Request) {
	p, err := fetchProposal(w, req)

	if err != nil {
		return
	}

	var b bytes.Buffer

	proposalUserTmpl.Execute(&b, p.Data)
	userPrompt := b.String()

	out := make(chan InferenceResult)

	go inference(InferenceRequest{
		proposalSystemPrompt,
		userPrompt,
	}, out)

	inferenceRes := <-out
	if inferenceRes.err != nil {
		http.Error(w, inferenceRes.err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(inferenceRes.result))
}

func summarizeThread(w http.ResponseWriter, req *http.Request) {
	p, err := fetchProposal(w, req)

	if err != nil {
		return
	}

	userPrompts, err := threadPrompts(p.Data.DiscussionThreadURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make(chan InferenceResult)
	go inference(InferenceRequest{
		threadSystemPrompt,
		userPrompts[0],
	}, out)

	inferenceRes := <-out
	if inferenceRes.err != nil {
		http.Error(w, inferenceRes.err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(inferenceRes.result))
}

func fetchProposal(w http.ResponseWriter, req *http.Request) (*ProposalResponse, error) {
	space := req.PathValue("space")
	id := req.PathValue("id")

	var p *ProposalResponse
	var err error

	if req.Method == "POST" {
		p, err = processProposal(req.Body)
	} else if req.Method == "GET" {
		if space == "" || id == "" {
			http.Error(w, "Invalid path (missing space or id)", http.StatusBadRequest)
			return nil, err
		}
		p, err = proposal(space, id)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return p, nil
}
