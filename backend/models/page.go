package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type AnalysedPage struct {
	gorm.Model
	AgentID uint
	Category string `json:"category"`
	URL string `json:"url"`
	Keywords string `json:"keywords"`
	Timestamp int64 `json:"timestamp"`
}

func SerializeKeywords(keywords []string) string {
	var serialized string
	for _, k := range keywords {
		serialized += k + "::"
	}
	return serialized
}

func DeserializeKeywords(serialized string) []string {
	return strings.Split(serialized, "::")
}

func NewAnalysedPage(agentID uint, category string, url string, keywords []string, timestamp int64) error {
	page := AnalysedPage{
		AgentID: agentID,
		Category: category,
		URL: url,
		Keywords: SerializeKeywords(keywords),
		Timestamp: timestamp,
	}
	if result := gormDB.Create(&page); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserPages(userID uint, from int64, to int64, category string) ([]AnalysedPage, error) {
	agents, err := GetUserAgents(userID)
	if err != nil {
		return nil, err
	}

	var pages []AnalysedPage
	query := "agent_id = ?"
	if category != "" {
		query += fmt.Sprintf(" AND category = '%s'", category)
	}
	if from != 0 || to != 0 {
		query += fmt.Sprintf(" AND timestamp BETWEEN %s AND %s", from, to)
	}

	for _, agent := range agents {
		var agentPages []AnalysedPage
		if result := gormDB.Where(query, agent.ID).Find(&agentPages); result.Error != nil {
				return nil, result.Error
		}
		pages = append(pages, agentPages...)
	}
	return pages, nil
}