package translation

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/sashabaranov/go-openai"
)

type OpenaiTranslator struct {
	apiKey string
	prompt string
}

// Initialises
func NewOpenaiTranslator(apiKey string, fromLang string, toLang string) *OpenaiTranslator {

	prompt := fmt.Sprintf(`You are a %s to %s translator app. 
	Translate any messages from %s to %s. 
	Make it sound as natural as possible and make sure to keep the newline formatting intact.`, fromLang, toLang, fromLang, toLang)

	return &OpenaiTranslator{
		apiKey,
		prompt,
	}

}

func (t *OpenaiTranslator) Translate(input string) (string, error) {

	// openai client doesn't seem concurrency safe.
	// therefore create one for each translation.

	client := openai.NewClient(t.apiKey)

	// Openai has rate limits.
	// Easy crude way to hande this is with an exponential backoff.
	// Keep retrying with the backoff if status is 429 - rate limit reached.

	b := backoff.NewExponentialBackOff()
	var success_res openai.ChatCompletionResponse

	// Operation retries everytime func returns an error.
	// We only want to retry if the error is due to ratelimiting, thus that is the only time we retry.
	op := func() error {
		res, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:     openai.GPT4,
				MaxTokens: getMaxTokens(input),
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: t.prompt,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: input,
					},
				},
			},
		)
		if isRateLimitErr(err) {
			return err
		} else {
			success_res = res
			return nil
		}

	}

	err := backoff.Retry(op, b)
	if err != nil {
		return "", err
	}

	return success_res.Choices[0].Message.Content, nil
}

func getMaxTokens(text string) int {

	// 1 token is roughly 4 chars in english.
	// Wan't some margin on the output - lets say 30% extra.
	// Always do at least 10.

	return max(int((float64(len(text))/4.0)*1.3), 10)
}

func isRateLimitErr(openaierr error) bool {

	// Handle nil case. Not sure if .As does.
	if openaierr == nil {
		return false
	}

	var openaiError = &openai.APIError{}
	if errors.As(openaierr, &openaiError) {
		switch openaiError.HTTPStatusCode {
		case 429:
			return true
		default:
			return false
		}
	}
	return false
}
