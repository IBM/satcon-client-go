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

// ChannelSubscription encapsulates a channel's subscription data
type ChannelSubscription struct {
	UUID        string     `json:"uuid,omitempty"`
	OrgID       string     `json:"orgId,omitempty"`
	Name        string     `json:"name,omitempty"`
	Groups      []string   `json:"groups,omitempty"`
	ChannelUUID string     `json:"channelUuid,omitempty"`
	ChannelName string     `json:"channelName,omitempty"`
	Channel     *Channel   `json:"channel,omitempty"`
	Version     string     `json:"version,omitempty"`
	VersionUUID string     `json:"versionUuid,omitempty"`
	Owner       *BasicUser `json:"owner,omitempty"`
	Created     string     `json:"created,omitempty"`
	Updated     string     `json:"updated,omitempty"`
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

type ClusterInfo struct {
	ClusterID string `json:"clusterId,omitempty"`
	Name      string `json:"name,omitempty"`
}

type ClusterList []Cluster

type ClusterGroup struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type ClusterGroupList []ClusterGroup

type Channel struct {
	UUID          string                `json:"uuid,omitempty"`
	OrgID         string                `json:"orgId,omitempty"`
	Name          string                `json:"name,omitempty"`
	Created       string                `json:"created,omitempty"`
	Versions      []ChannelVersion      `json:"versions,omitempty"`
	Subscriptions []ChannelSubscription `json:"subscriptions,omitempty"`
}

type ChannelList []Channel

type ChannelVersion struct {
	UUID        string `json:"uuid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
	Created     string `json:"created,omitempty"`
}

type ChannelVersionList []ChannelVersion

type Comment struct {
	UserId  string `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
	Created string `json:"created,omitempty"`
}

type DeployableVersion struct {
	OrgID       string `json:"orgId,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	ChannelID   string `json:"channelId,omitempty"`
	ChannelName string `json:"channelName,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	Created     string `json:"created,omitempty"`
}

type RequestError struct {
	Errors []RequestErrorDetails `json:"errors,omitempty"`
}

type RequestErrorDetails struct {
	Message string `json:"message,omitempty"`
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

// Resource encapsulates satellite cluster resources
type Resource struct {
	ID                 string              `json:"id,omitempty"`
	OrgID              string              `json:"orgId,omitempty"`
	ClusterID          string              `json:"clusterId,omitempty"`
	Cluster            ClusterInfo         `json:"cluster,omitempty"`
	HistId             string              `json:"histId,omitempty"`
	SelfLink           string              `json:"selfLink,omitempty"`
	Hash               string              `json:"hash,omitempty"`
	Data               string              `json:"data,omitempty"`
	Deleted            bool                `json:"deleted,omitempty"`
	Created            string              `json:"created,omitempty"`
	Updated            string              `json:"updated,omitempty"`
	LastModified       string              `json:"lastModified,omitempty"`
	SearchableData     SearchableData      `json:"searchableData,omitempty"`
	SearchableDataHash string              `json:"searchableDataHash,omitempty"`
	Subscription       ChannelSubscription `json:"subscription,omitempty"`
}

// ResourceList encapsulates list of resource
type ResourceList struct {
	Count     int        `json:"count,omitempty"`
	Resources []Resource `json:"resources,omitempty"`
}

type ResourceContentObj struct {
	ID      string `json:"id,omitempty"`
	HistID  string `json:"histId,omitempty"`
	Content string `json:"content,omitempty"`
	Updated string `json:"updated,omitempty"`
}

// SearchableData encapsulates cluster resource data
type SearchableData struct {
	Kind                 string                 `json:"kind,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	Namespace            string                 `json:"namespace,omitempty"`
	APIVersion           string                 `json:"apiVersion,omitempty"`
	SearchableExpression string                 `json:"searchableExpression,omitempty"`
	Errors               map[string]interface{} `json:"errors,omitempty"`
}

// Subscription encapsulates satellite subscription data
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

// SubscriptionList list of subscriptions
type SubscriptionList []Subscription
