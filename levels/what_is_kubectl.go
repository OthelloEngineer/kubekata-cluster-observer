package levels

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels/levelutils"
)

type WhatIsKubbectl struct {
	isFinished bool
}

func (l *WhatIsKubbectl) GetName() string {
	return "what is kubectl"
}

func (l *WhatIsKubbectl) GetDesiredCluster() client.Cluster {

	return levelutils.GetEmptyCluster()
}

func (l *WhatIsKubbectl) GetClusterStatus(cluster client.Cluster, msg string) string {
	return "success"
}

func (l *WhatIsKubbectl) SetFinished() {
	l.isFinished = true
}

func (l *WhatIsKubbectl) GetIsFinished() bool {
	return l.isFinished
}
