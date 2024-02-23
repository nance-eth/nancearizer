package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/daulet/tokenizers"
	"github.com/joho/godotenv"
)

const NANCE_API = "https://api.nance.app/"

var (
	s *discordgo.Session
	t *tokenizers.Tokenizer
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	if os.Getenv("DISCORD_TOKEN") == "" {
		log.Fatal("No DISCORD_TOKEN found in .env")
	}

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
	http.ListenAndServe(":8080", nil)
}
