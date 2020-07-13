package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	Details DeleteClustersResponseDataDetails `json:"deleteClusterByClusterID,omitempty"`
}

type DeleteClustersResponseDataDetails struct {
	DeletedClusterCount  int `json:"deletedClusterCount,omitempty"`
	DeletedResourceCount int `json:"deletedResourceCount,omitempty"`
}

func (d *DeleteClustersResponseData) String() string {
	var response string
	if d.Details.DeletedClusterCount > 0 {
		response = fmt.Sprintf("Deleted Clusters: %d\n", d.Details.DeletedClusterCount)
	}

	if d.Details.DeletedResourceCount > 0 {
		response += fmt.Sprintf("Deleted Resources: %d\n", d.Details.DeletedResourceCount)
	}

	return response
}

func (c *Client) DeleteClusterByClusterID(orgID, clusterID, token string) (*DeleteClustersResponseData, error) {
	var response DeleteClustersResponse

	vars := NewDeleteClusterByClusterIDVariables(orgID, clusterID)

	payload, err := actions.BuildRequestBody(DeleteClusterByClusterIDVarTemplate, vars, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to build request body: %s", err)
	}

	req, _ := http.NewRequest(http.MethodPost, c.Endpoint, payload)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

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
		return response.Data, err
	}

	return nil, err
}
