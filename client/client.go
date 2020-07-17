package client

import (
	"github.ibm.com/coligo/satcon-client/client/actions/clusters"
	"github.ibm.com/coligo/satcon-client/client/actions/groups"
)

type SatCon struct {
	Clusters clusters.ClusterService
	Groups   groups.GroupService
}
