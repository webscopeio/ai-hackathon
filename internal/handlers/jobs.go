package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func CreateJob(client *llm.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := decode[models.CreateJobArgs](r)
		if err != nil {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: fmt.Sprintf("Bad request, %v", err),
			})
			return
		}

		if args.Prompt == "" {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: "Prompt is required",
			})
			return
		}

		basePrompt := `Generate comprehensive job listings from user descriptions by extracting or inferring the necessary details to populate a structured JSON object with these fields: "title" (clear, industry-standard job title), "description" (2-3 paragraph role summary including context and value proposition), "requirements" (4-8 specific qualifications as complete sentences), "responsibilities" (5-10 specific duties starting with action verbs), "experience_level" (integer from 0-5, where 0=unknown, 1=junior, 2=mid-level, 3=senior, 4-5=executive/leadership), "skills" (5-15 specific technical and soft skills as short phrases), and "keywords" (8-12 searchable terms including industry terminology and abbreviations). When information is missing, make reasonable inferences based on industry standards and maintain internal consistency across all fields. Return valid JSON that strictly follows this structure without additional commentary unless specifically requested.`
		prompt := fmt.Sprintf("%s %s", basePrompt, args.Prompt)

		tool, toolChoice := llm.GenerateTool[models.CreateJobReturn]("create_job_return", "structure response for generation a job posting based on a prompt")
		res, err := client.GetStructuredCompletion(
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

		var response models.CreateJobReturn
		if err := json.Unmarshal(res, &response); err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Couldn't process response, %v", err),
			})
			return
		}

		if err := response.Validate(); err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Validation fail, %v", err),
			})
			return
		}

		encode(w, http.StatusOK, response)
	}
}
