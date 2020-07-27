package types

type BasicUser struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type BasicChannelSubscription struct {
	UUID        string   `json:"uuid,omitempty"`
	OrgID       string   `json:"orgId,omitempty"`
	Name        string   `json:"name,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	ChannelUUID string   `json:"channelUuid,omitempty"`
	ChannelName string   `json:"channelName,omitempty"`
	Version     string   `json:"version,omitempty"`
	VersionUUID string   `json:"versionUuid,omitempty"`
	Created     string   `json:"created,omitempty"`
	Updated     string   `json:"updated,omitempty"`
}

type ChannelSubscription struct {
	UUID        string    `json:"uuid,omitempty"`
	OrgID       string    `json:"orgId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Groups      []string  `json:"groups,omitempty"`
	ChannelUUID string    `json:"channelUuid,omitempty"`
	ChannelName string    `json:"channelName,omitempty"`
	Channel     Channel   `json:"channel,omitempty"`
	Version     string    `json:"version,omitempty"`
	VersionUUID string    `json:"versionUuid,omitempty"`
	Owner       BasicUser `json:"owner,omitempty"`
	Created     string    `json:"created,omitempty"`
	Updated     string    `json:"updated,omitempty"`
}

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
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type ClusterGroupList []ClusterGroup

type Channel struct {
	UUID          string                     `json:"uuid,omitempty"`
	OrgID         string                     `json:"orgId,omitempty"`
	Name          string                     `json:"name,omitempty"`
	Created       string                     `json:"created,omitempty"`
	Versions      []ChannelVersion           `json:"versions,omitempty"`
	Subscriptions []BasicChannelSubscription `json:"subscriptions,omitempty"`
}

type ChannelList []Channel

type ChannelVersion struct {
	UUID        string `json:"uuid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
	Created     string `json:"created,omitempty"`
}

type Comment struct {
	UserId  string `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
	Created string `json:"created,omitempty"`
}

type GroupList []Group

type Group struct {
	UUID    string    `json:"uuid,omitempty"`
	OrgID   string    `json:"orgId,omitempty"`
	Name    string    `json:"name,omitempty"`
	Owner   BasicUser `json:"owner,omitempty"`
	Created string    `json:"created,omitempty"`
}

// Registration is the encapsulation of the JSON registration body, which at this
// point is used primarily to specify the name of the cluster to register.
type Registration struct {
	Name string `json:"name,omitempty"`
}

type Subscription struct {
	UUID        string    `json:"uuid,omitempty"`
	OrgID       string    `json:"orgId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Groups      []string  `json:"groups,omitempty"`
	ChannelUUID string    `json:"channelUuid,omitempty"`
	ChannelName string    `json:"channelName,omitempty"`
	Channel     Channel   `json:"channel,omitempty"`
	Version     string    `json:"version,omitempty"`
	VersionUUID string    `json:"versionUuid,omitempty"`
	Owner       BasicUser `json:"owner,omitempty"`
	Created     string    `json:"created,omitempty"`
	Updated     string    `json:"updated,omitempty"`
}

type SubscriptionList []Subscription
