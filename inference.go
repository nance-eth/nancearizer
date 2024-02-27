package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	apiUrl string
	apiKey string
)

const (
	MODEL_NAME     = "gpt-3.5-turbo-0125"
	CONTEXT_LENGTH = 16_384
	TEMPERATURE    = 0.7
)

type InferenceRequest struct {
	systemPrompt string
	userPrompt   string
	// maxTokens    int
}

type InferenceResult struct {
	result string
	err    error
}

// Types for parsing the inference response.
type ResponseBody struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func inference(r InferenceRequest, out chan<- InferenceResult) {
	req, err := http.NewRequest("POST", apiUrl+"/v1/chat/completions", nil)
	if err != nil {
		out <- InferenceResult{"", err}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	requestBody := map[string]interface{}{
		// "max_tokens":  r.maxTokens,
		"model":       MODEL_NAME,
		"temperature": TEMPERATURE,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": r.systemPrompt,
			},
			{
				"role":    "user",
				"content": r.userPrompt,
			},
		},
	}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		out <- InferenceResult{"", err}
		return
	}

	// Must wrap bytes.Buffer to satisfy io.ReadCloser interface
	req.Body = io.NopCloser(bytes.NewBuffer(requestBodyJson))

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		out <- InferenceResult{"", err}
		return
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		out <- InferenceResult{"", err}
		return
	}

	// Check for HTTP errors.
	if resp.StatusCode != http.StatusOK {
		log.Println("Inference request failed with status code " + fmt.Sprint(resp.StatusCode) + ".")
		log.Println("responseBody: " + string(responseBody))
		out <- InferenceResult{"", fmt.Errorf("inference request returned status code %d", resp.StatusCode)}
		return
	}

	// Parse the response body.
	var response ResponseBody
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		out <- InferenceResult{"", err}
		return
	}

	// Extract the result from the response
	if len(response.Choices) == 0 {
		out <- InferenceResult{"", fmt.Errorf("inference request returned no results")}
		return
	}

	if response.Choices[0].FinishReason != "stop" && response.Choices[0].FinishReason != "eos" {
		log.Println("Inference request finished with reason " + response.Choices[0].FinishReason + ".")
	}

	log.Printf(
		"Inference request finished with reason '%s'. Input tokens: %d, output tokens: %d\n",
		response.Choices[0].FinishReason,
		response.Usage.PromptTokens,
		response.Usage.CompletionTokens,
	)

	out <- InferenceResult{result: response.Choices[0].Message.Content}
}
