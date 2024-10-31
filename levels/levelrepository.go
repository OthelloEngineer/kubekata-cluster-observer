package levels

import (
	"errors"
	"fmt"
	"strings"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

type Level interface {
	GetName() string
	GetDesiredCluster(client client.Client) client.Cluster
	GetClusterStatus(cluster client.Cluster, msg string) string
	SetFinished()
	GetIsFinished() bool
}

type LevelRepository struct {
	levels       []Level
	currentLevel string
}

func NewLevelRepository() *LevelRepository {
	levels := []Level{
		new(WhatIsKubeKata),
		new(ComponentsOfKubeKata),
		new(WhatIsKubectl),
		new(DeployingTheApp),
		new(CurlingTheApp),
		new(dns_and_services),
		new(ScalingTheApp),
		new(ExposingToTheWorld),
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
		println("no current level set, level name: ", name)
		return nil, errors.New("no current level set")
	}
	for _, level := range s.levels {
		if strings.EqualFold(level.GetName(), name) {
			fmt.Println("EQUAL: repo level name: ", level.GetName(), "searched for: ", name)
			return level, nil
		}
		fmt.Println("NOT EQUAL: repo level name: ", level.GetName(), "searched for: ", name)
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
