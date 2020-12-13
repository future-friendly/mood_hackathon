package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/future-friednly/mood/backend/models"
	"github.com/future-friednly/mood/backend/util"
)

type ChartType uint

const (
	InterestMap ChartType = iota
	CategoryMap
	SurfPipeline
)

var ChartToEndpoint map[ChartType]string = map[ChartType]string{
	InterestMap:  "interest_chart",
	CategoryMap:  "keyword_chart",
	SurfPipeline: "surf_chart",
}

type GetChartRequest struct {
	Token    string    `json:"token"`
	Type     ChartType `json:"chart_type"`
	Category string    `json:"category"`
	From     int64     `json:"from"`
	To       int64     `json:"to"`
}

func HandleGetChart(w http.ResponseWriter, r *http.Request) {
	var req GetChartRequest
	if err := util.DecodeRequest(r.Body, &req); err != nil {
		util.WriteError(w, 500, err)
		return
	}

	token, err := models.GetToken(req.Token)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	pages, err := models.GetUserPages(token.UserID, req.From, req.To, req.Category)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	body, err := json.Marshal(map[string][]models.AnalysedPage{
		"data": pages,
	})
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("http://analytics:8080/%s", ChartToEndpoint[req.Type]), "application/json", bytes.NewBuffer(body))
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.WriteError(w, 500, err)
		return
	}

	w.WriteHeader(200)
	w.Write(respBody)
}
