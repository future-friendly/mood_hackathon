package data

import (
	"github.com/future-friednly/mood/backend/models"
	"github.com/future-friednly/mood/backend/util"
	"net/http"
)

type NewPageRequest struct {
	AgentToken string `json:"agent_token"`
	Category string `json:"category"`
	URL string `json:"url"`
	Keywords []string `json:"keywords"`
	Timestamp int64 `json:"timestamp"`
}

func HandleNewPage(w http.ResponseWriter, r *http.Request) {
	var req NewPageRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	agent, err := models.GetAgent(nil, &req.AgentToken)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}
	if err := models.NewAnalysedPage(agent.ID, req.Category, req.URL, req.Keywords, req.Timestamp); err != nil {
		util.WriteError(w, 500, err)
		return
	}
	util.WriteSuccess(w)
}
