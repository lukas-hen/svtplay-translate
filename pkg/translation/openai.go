package translation

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

const PROMPT = `You are a Swedish to English translator app. 
Translate any messages from swedish to english. 
Make it sound as natural as possible and make sure to keep the newline formatting intact.`

type Translator interface {
	Translate(string) (string, error)
}

type OpenaiTranslator struct {
	apiKey string
}

func NewOpenaiTranslator() *OpenaiTranslator {

	apiKey := os.Getenv("OPENAI_API_KEY")

	return &OpenaiTranslator{
		apiKey,
	}

}

func (t *OpenaiTranslator) Translate(input string) (string, error) {

	// openai client doesn't seem concurrency safe.
	// therefore create one for each translation.

	client := openai.NewClient(t.apiKey)

	res, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4,
			MaxTokens: getMaxTokens(input),
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
		return "", err
	}

	return res.Choices[0].Message.Content, nil
}

func getMaxTokens(text string) int {
	// 1 token is roughly 4 chars in english.
	// Wan't some margin on the output - lets say 30% extra.
	// Always do at least 10.

	return max(int((float64(len(text))/4.0)*1.3), 10)
}
