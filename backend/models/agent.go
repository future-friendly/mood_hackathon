package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

type MonAgentType uint

const (
	Chrome MonAgentType = iota
	Android
	IOS
)

type MonAgent struct {
	gorm.Model
	Name string `json:"name"`
	UserID uint `json:"user_id"`
	Type MonAgentType `json:"agent_type"`
	Token string `json:"agent_token"`
	Confirmed bool `json:"confirmed"`
}

func NewAgent(userID uint, name string, agentType MonAgentType) error {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	agent := &MonAgent{
		Name: name,
		UserID:    userID,
		Type:      agentType,
		Token:     hex.EncodeToString(bytes),
		Confirmed: false,
	}
	if result := gormDB.Create(agent); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAgent(id *uint, token *string) (*MonAgent, error) {
	var agent MonAgent
	if id != nil {
		if result := gormDB.First(&agent, *id); result.Error != nil {
			return nil, result.Error
		}
		return &agent, nil
	}

	if token != nil {
		if result := gormDB.Where(&MonAgent{Token: *token}).First(&agent); result.Error != nil {
			return nil, result.Error
		}
		return &agent, nil
	}

	return nil, errors.New("args not provided")
}

func (a *MonAgent) ConfirmAgent() error {
	a.Confirmed = true
	if result := gormDB.Save(a); result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *MonAgent) Delete() error {
	if result := gormDB.Delete(a); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserAgents(userId uint) ([]MonAgent, error) {
	var agents []MonAgent
	if result := gormDB.Where(&MonAgent{UserID: userId}).Find(&agents); result.Error != nil {
		return nil, result.Error
	}
	return agents, nil
}