package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryDeleteClusterByClusterID       = "deleteClusterByClusterID"
	DeleteClusterByClusterIDVarTemplate = `{{define "vars"}}"org_id":"{{.OrgID}}","cluster_id":"{{.ClusterID}}"{{end}}`
)

type DeleteClusterByClusterIDVariables struct {
	actions.GraphQLQuery
	OrgID     string
	ClusterID string
}

func NewDeleteClusterByClusterIDVariables(orgID, clusterID string) DeleteClusterByClusterIDVariables {
	vars := DeleteClusterByClusterIDVariables{
		OrgID:     orgID,
		ClusterID: clusterID,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryDeleteClusterByClusterID
	vars.Args = map[string]string{
		"org_id":     "String!",
		"cluster_id": "String!",
	}
	vars.Returns = []string{
		"deletedClusterCount",
		"deletedResourceCount",
	}

	return vars
}

type DeleteClustersResponse struct {
	Data *DeleteClustersResponseData `json:"data,omitempty"`
}

type DeleteClustersResponseData struct {
	Details *DeleteClustersResponseDataDetails `json:"deleteClusterByClusterID,omitempty"`
}

type DeleteClustersResponseDataDetails struct {
	DeletedClusterCount  int `json:"deletedClusterCount,omitempty"`
	DeletedResourceCount int `json:"deletedResourceCount,omitempty"`
}

func (d *DeleteClustersResponseDataDetails) String() string {
	var response string
	if d.DeletedClusterCount > 0 {
		response = fmt.Sprintf("Deleted Clusters: %d\n", d.DeletedClusterCount)
	}

	if d.DeletedResourceCount > 0 {
		response += fmt.Sprintf("Deleted Resources: %d\n", d.DeletedResourceCount)
	}

	return response
}

func (c *Client) DeleteClusterByClusterID(orgID, clusterID, token string) (*DeleteClustersResponseDataDetails, error) {
	var response DeleteClustersResponse

	vars := NewDeleteClusterByClusterIDVariables(orgID, clusterID)

	payload, _ := actions.BuildRequestBody(DeleteClusterByClusterIDVarTemplate, vars, nil)

	req := actions.BuildRequest(payload, c.Endpoint, token)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshal response: %s", err)
		}
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
