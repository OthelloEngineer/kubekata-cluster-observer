package levels

import (
	"github.com/OthelloEngineer/kubekata-cluster-observer/client"
	"github.com/OthelloEngineer/kubekata-cluster-observer/levels/levelutils"
)

/*
	The goal of level1 is to create an nginx deployment
*/

type WhatIsKubeKata struct {
	isFinished bool
}

func (l *WhatIsKubeKata) GetName() string {
	return "what is kubekata"
}

func (l *WhatIsKubeKata) GetDesiredCluster(client client.Client) client.Cluster {

	return levelutils.GetEmptyCluster()
}

func (l *WhatIsKubeKata) GetClusterStatus(cluster client.Cluster, msg string) string {
	return "success"
}

func (l *WhatIsKubeKata) SetFinished() {
	l.isFinished = true
}

func (l *WhatIsKubeKata) GetIsFinished() bool {
	return l.isFinished
}
