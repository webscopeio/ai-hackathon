package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamhectorsosa/ai-hackathon/internal/llm"
	"github.com/iamhectorsosa/ai-hackathon/internal/models"
)

func Ask(client *llm.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := decode[models.AskArgs](r)
		if err != nil {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: fmt.Sprintf("Bad request, %v", err),
			})
			return
		}

		if args.Question == "" {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: "Question is required",
			})
			return
		}

		basePrompt := `Please limit your answer to 1 very concise sentence only.`
		prompt := fmt.Sprintf("%s %s", basePrompt, args.Question)

		// INFO: for a structured response the client requires tools, ref: https://docs.anthropic.com/en/docs/build-with-claude/tool-use/overview
		tool, toolChoice := llm.GenerateTool[models.AskReturn]("get_ask_return", "structured response to the question asked")
		answer, err := client.GetStructuredCompletion(
			r.Context(),
			"",
			prompt,
			tool,
			toolChoice,
		)
		if err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Couldn't process request, %v", err),
			})
			return
		}

		var response models.AskReturn
		if err := json.Unmarshal(answer, &response); err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Couldn't process response, %v", err),
			})
			return
		}

		if response.Answer == "" {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Couldn't process response, %v", err),
			})
			return
		}

		encode(w, http.StatusOK, response)
	}
}
