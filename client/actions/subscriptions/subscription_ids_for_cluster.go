package subscriptions

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	//QuerySubscriptionIdsForCluster specifies the query
	QuerySubscriptionIdsForCluster = "subscriptionsForCluster"
	//SubscriptionIdsForClusterVarTemplate is the template used to create the graphql query
	SubscriptionIdsForClusterVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"clusterId":{{json .ClusterID}}{{end}}`
)

//SubscriptionIdsForClusterVariables are the variables used for the subscription query
type SubscriptionIdsForClusterVariables struct {
	actions.GraphQLQuery
	OrgID     string
	ClusterID string
}

//NewSubscriptionIdsForClusterVariables generates variables used for query
func NewSubscriptionIdsForClusterVariables(orgID string, clusterID string) SubscriptionIdsForClusterVariables {
	vars := SubscriptionIdsForClusterVariables{
		OrgID: orgID,
		ClusterID: clusterID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QuerySubscriptionIdsForCluster
	vars.Args = map[string]string{
		"orgId": "String!",
		"clusterId": "String!",
	}
	vars.Returns = []string{
		"uuid",
	}

	return vars
}

//SubscriptionIdsForClusterResponse response from query
type SubscriptionIdsForClusterResponse struct {
	Data *SubscriptionIdsForClusterResponseData `json:"data,omitempty"`
}

// SubscriptionIdsForClusterResponseData data response from query
type SubscriptionIdsForClusterResponseData struct {
	SubscriptionIdsForCluster []types.UuidOnly `json:"subscriptionsForCluster,omitempty"`
}

func (c *Client) SubscriptionIdsForCluster(orgID string, clusterID string) ([]string, error) {
	var response SubscriptionIdsForClusterResponse

	vars := NewSubscriptionIdsForClusterVariables(orgID, clusterID)

	err := c.DoQuery(SubscriptionIdsForClusterVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		uuidResponses := response.Data.SubscriptionIdsForCluster
		uuids := make([]string, len(uuidResponses))
		for i := 0; i < len(uuidResponses); i++ {
			uuids[i] = uuidResponses[i].UUID
		}
		return uuids, nil
	}

	return nil, err
}
