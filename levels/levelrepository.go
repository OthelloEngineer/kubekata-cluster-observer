package levels

import (
	"errors"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

type Level interface {
	GetID() int
	GetName() string
	GetDesiredCluster() client.Cluster
	GetClusterDiff(client.Cluster) string
}

type LevelRepository struct {
	levels       []Level
	currentLevel int
}

func NewLevelRepository() *LevelRepository {
	levels := []Level{
		new(Level1),
	}
	repo := &LevelRepository{
		levels:       levels,
		currentLevel: 1,
	}
	return repo
}

func (s *LevelRepository) GetAllLevels() []Level {
	return s.levels
}

func (s *LevelRepository) GetLevelByID(id int) (Level, error) {
	for _, level := range s.levels {
		if level.GetID() == id {
			return level, nil
		}
	}
	return nil, errors.New("level not found")
}

func (s *LevelRepository) GetCurrentLevel() (Level, error) {
	if s.currentLevel == 0 {
		return nil, errors.New("no current level set")
	}
	return s.GetLevelByID(s.currentLevel)
}

func (s *LevelRepository) SetCurrentLevel(id int) {
	s.currentLevel = id
}
