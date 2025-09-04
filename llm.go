package main

import (
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	// "github.com/openai/openai-go/v2/shared"
)

var llmClient openai.Client

func init() {
	llmClient = openai.NewClient(
		option.WithAPIKey(""),
		option.WithBaseURL(os.Getenv("OPENAI_API_URL")),
	)
}
