package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/daulet/tokenizers"
	"github.com/joho/godotenv"
)

var (
	s    *discordgo.Session
	t    *tokenizers.Tokenizer
	port string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	for _, v := range []string{"DISCORD_TOKEN", "OPENAI_API_KEY", "OPENAI_API_URL", "PORT"} {
		if os.Getenv(v) == "" {
			log.Fatalf("No %s found in env\n", v)
		}
	}

	apiUrl = os.Getenv("OPENAI_API_URL")
	apiKey = os.Getenv("OPENAI_API_KEY")
	port = os.Getenv("PORT")

	s, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating Discord session: %v\n", err)
	}

	t, err = tokenizers.FromFile("tokenizer.json")
	if err != nil {
		log.Fatalf("Error loading tokenizer: %v\n", err)
	}
}

func main() {
	http.HandleFunc("GET /proposal/{space}/{id}", summarizeProposal)
	http.HandleFunc("GET /thread/{space}/{id}", summarizeThread)

	http.ListenAndServe(":"+port, nil)
}

func summarizeProposal(w http.ResponseWriter, req *http.Request) {}

func summarizeThread(w http.ResponseWriter, req *http.Request) {}
