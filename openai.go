package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/cenkalti/backoff"
	"github.com/sashabaranov/go-openai"
)

const PROMPT = `You are a Swedish to English translator app. 
Translate any messages from swedish to english. 
Make it sound as natural as possible and make sure to keep the newline formatting intact.`

func Translate(input string, responseBuf []string, index int, wg *sync.WaitGroup) {

	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(apiKey)

	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	var err error
	var res openai.ChatCompletionResponse

	// Ticks will continue to arrive when the previous operation is still running,
	// so operations that take a while to fail could run in quick succession.
	for range ticker.C {

		res, err = client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:     openai.GPT4,
				MaxTokens: GetMaxTokens(input),
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: PROMPT,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: input,
					},
				},
			},
		)

		if err != nil {
			// Rate limit is also error here
			log.Printf("error: %s, backing off.", err)
			continue
		}

		responseBuf[index] = res.Choices[0].Message.Content

		ticker.Stop()
		wg.Done()
		break
	}

}

func GetMaxTokens(text string) int {
	// 1 token is roughly 4 chars in english.
	// Wan't some margin on the output - lets say 30% extra.
	// Always do at least 10.

	return max(int((float64(len(text))/4.0)*1.3), 10)
}
