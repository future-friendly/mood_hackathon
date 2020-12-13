package agents

import (
	"encoding/json"
	"errors"
	"github.com/future-friednly/mood/backend/models"
	"github.com/future-friednly/mood/backend/util"
	"net/http"
)

type CreateAgentRequest struct {
	Token string `json:"token"`
	Name string `json:"name"`
	Type models.MonAgentType `json:"agent_type"`
}

func HandleCreateAgent(w http.ResponseWriter, r *http.Request) {
	var req CreateAgentRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	token, err := models.GetToken(req.Token)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}
	user, err := token.GetUserFromToken()
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	if err := models.NewAgent(user.ID, req.Name, req.Type); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	util.WriteSuccess(w)
}

type ConfirmAgentRequest struct {
	AgentToken string `json:"agent_token"`
}

func HandleConfirmAgent(w http.ResponseWriter, r *http.Request) {
	var req ConfirmAgentRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	agent, err := models.GetAgent(nil, &req.AgentToken)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}
	if err := agent.ConfirmAgent(); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	util.WriteSuccess(w)
}

type GetAgentsRequest struct {
	Token string `json:"token"`
}

type GetAgentsResponse struct {
	Agents []models.MonAgent `json:"agents"`
}

func HandleGetAgents(w http.ResponseWriter, r *http.Request) {
	var req GetAgentsRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	token, err := models.GetToken(req.Token)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}
	user, err := token.GetUserFromToken()
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	agents, err := models.GetUserAgents(user.ID)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	resp, err := json.Marshal(GetAgentsResponse{Agents: agents})
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	w.WriteHeader(200)
	w.Write(resp)
}

type DeleteAgentRequest struct {
	Token string `json:"token"`
	AgentID uint `json:"agent_id"`
}

func HandleDeleteAgent(w http.ResponseWriter, r *http.Request) {
	var req DeleteAgentRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	token, err := models.GetToken(req.Token)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}
	user, err := token.GetUserFromToken()
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	agent, err := models.GetAgent(&req.AgentID, nil)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}
	if agent.UserID != user.ID {
		util.WriteError(w, 403, errors.New("permission denied"))
		return
	}
	if err := agent.Delete(); err != nil {
		util.WriteError(w, 500, err)
		return
	}
	util.WriteSuccess(w)
}