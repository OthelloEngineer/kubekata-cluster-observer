package levels

import (
	"errors"

	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
)

type Level interface {
	GetID() int
	GetName() string
	GetDesiredCluster() client.Cluster
	GetClusterStatus(cluster client.Cluster, msg string) string
}

type LevelRepository struct {
	levels       []Level
	currentLevel int
}

func NewLevelRepository() *LevelRepository {
	levels := []Level{
		new(Level1),
		new(Level2),
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
	if s.currentLevel > len(s.levels) {
		return nil, errors.New("current level out of bounds")
	}
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

func (s *LevelRepository) SetCurrentLevel(id int) (Level, error) {
	s.currentLevel = id
	return s.GetLevelByID(id)
}
