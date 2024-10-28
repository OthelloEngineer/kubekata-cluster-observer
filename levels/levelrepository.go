package levels

import (
	"errors"
	"strings"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

type Level interface {
	GetName() string
	GetDesiredCluster() client.Cluster
	GetClusterStatus(cluster client.Cluster, msg string) string
	SetFinished()
}

type LevelRepository struct {
	levels       []Level
	currentLevel string
}

func NewLevelRepository() *LevelRepository {
	levels := []Level{
		new(WhatIsKubeKata),
		new(ComponentsOfKubeKata),
		new(DeployingTheApp),
		new(CurlingTheApp),
		new(dns_and_services),
	}
	repo := &LevelRepository{
		levels:       levels,
		currentLevel: "what is KubeKata",
	}
	return repo
}

func (s *LevelRepository) GetAllLevels() []Level {
	return s.levels
}

func (s *LevelRepository) GetLevelByName(name string) (Level, error) {
	if s.currentLevel == "" {
		return nil, errors.New("no current level set")
	}
	for _, level := range s.levels {
		if strings.EqualFold(level.GetName(), name) {
			return level, nil
		}
	}
	return nil, errors.New("level not found")
}

func (s *LevelRepository) GetCurrentLevel() (Level, error) {
	if s.currentLevel == "" {
		return nil, errors.New("no current level set")
	}
	return s.GetLevelByName(s.currentLevel)
}

func (s *LevelRepository) SetCurrentLevel(name string) (Level, error) {
	s.currentLevel = name
	return s.GetLevelByName(name)
}
