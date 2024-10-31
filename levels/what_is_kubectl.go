package levels

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels/levelutils"
)

type WhatIsKubectl struct {
	isFinished bool
}

func (l *WhatIsKubectl) GetName() string {
	return "what is kubectl"
}

func (l *WhatIsKubectl) GetDesiredCluster() client.Cluster {

	return levelutils.GetEmptyCluster()
}

func (l *WhatIsKubectl) GetClusterStatus(cluster client.Cluster, msg string, client client.Client) string {
	return "success"
}

func (l *WhatIsKubectl) SetFinished() {
	l.isFinished = true
}

func (l *WhatIsKubectl) GetIsFinished() bool {
	return l.isFinished
}
