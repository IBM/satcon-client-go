package types

// type Channel struct {
// 	UUID string `json:"uuid,omitempty"`
// 	OrgID string `json:"orgId,omitempty"`
// 	Name string `json:"name,omitempty"`
// 	Created string `json:"created,omitempty"`
// 	Versions ChannelVersionList `json:"versions,omitempty"`
// 	Subscriptions
// }

type Cluster struct {
	ID        string `json:"id,omitempty"`
	OrgID     string `json:"orgId,omitempty"`
	ClusterID string `json:"clusterId,omitempty"`
	// Metadata          []byte         `json:"metadata,omitempty"`
	Metadata          interface{}    `json:"metadata,omitempty"`
	Comments          []Comment      `json:"comments,omitempty"`
	Registration      Registration   `json:"registration,omitempty"`
	RegistrationState string         `json:"regState,omitempty"`
	Groups            []ClusterGroup `json:"groups,omitempty"`
	Created           string         `json:"created,omitempty"`
	Updated           string         `json:"updated,omitempty"`
	Dirty             bool           `json:"dirty,omitempty"`
}

type ClusterList []Cluster

type ClusterGroup struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type ClusterGroupList []ClusterGroup

type Comment struct {
	UserId  string `json:"user_id"`
	Content string `json:"content"`
	Created string `json:"created"`
}

// Registration is the encapsulation of the JSON registration body, which at this
// point is used primarily to specify the name of the cluster to register.
type Registration struct {
	Name string `json:"name"`
}
