package cluster

import (
	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryClustersByOrgID = "clustersByOrgID"
)

type ClustersByOrgIDVariables struct {
	actions.GraphQLQuery
	OrgID string
}

func NewClustersByOrgIDVariables(orgID string) ClustersByOrgIDVariables {
	vars := ClustersByOrgIDVariables{
		OrgID: orgID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryClustersByOrgID
	vars.Args = map[string]string{
		"org_id": "String!",
	}
	vars.Returns = []string{
		"_id",
		"org_id",
		"cluster_id",
		"metadata",
	}

	return vars
}

func (c *Client) ClustersByOrgID(orgID, token string) error {
	return nil
}
