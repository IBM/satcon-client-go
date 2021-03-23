# satcon-client

## Usage

Although the purpose of the library is to insulate users to some extent from the details of the IBM Cloud Satellite Config (SatCon) API, an understanding of key concepts and, in some cases, actual API details is required to use it successfully.  Users new to SatCon and/or Razee (the IBM-developed open source tooling that is the core of Satellite Config) may want to visit these links to gain familiarity:

- [Satellite Config API Schema](https://github.com/razee-io/Razeedash-api/tree/master/app/apollo/schema) - can help with understanding entity details and relationships
- [Razee.io docs](https://github.com/razee-io/Razee/blob/master/README.md) - primarily useful for users wanting to understand more of what's under the covers.

**This client does not officially support interacting with a stand-alone open source Razee deployment. It will probably work, and we do even accept PRs only relevant to Razee (as long as they don't break overall SatCon functionality). However, bugs filed that are only relevant to Razee usage will not be fixed.**

### Key objects in Satellite Config

#### NOTE: The GraphQL API upon which this client library is currently based is for all practical purposes the standard Razee API. We therefore primarily use Razee terminology as that is what is used in the API. IBM Cloud Satellite Config uses slightly different terminology in its documentation, UI, and command-line client. We point this out when relevant.

There are six primary entity classes in SatCon (or Razee), at least from a user perspective.  These are:

1. **Cluster** - represents an actual Kubernetes cluster to be managed via SatCon
1. **Cluster group** - clusters to be managed similarly are collected into groups to facilitate organized and simplified distribution of resources
1. **Version** - a version represents one or more Kubernetes objects to be created and/or configured. It usually ultimately takes the form of a YAML file as is typical for defining Kubernetes resources
1. **Channel** - a channel collects a set of versions which can potentially then be distributed to cluster groups. *NOTE*: the term Satellite Config uses for this is **Configuration**
1. **Subscription** - a subscription associates a cluster group with a version and the owning channel/configuration, so that the SatCon agent will then know to pull that version down and apply it to every cluster within that group
1. **Resource** - versions applied to clusters produce resources, which are essentially just that - representations of Kubernetes resources which have been deployed to the cluster(s)

This diagram provides an overview of the relationships between the various entities:

![SatCon Entity Relationships](diagrams/images/CE_Isolation_SatCon_Workflow.png)

### SatCon Workflow

A typical SatCon workflow might proceed roughly as follows:

1. **Register a cluster to SatCon.**  This makes the cluster known to SatCon, and the API call will return the URL for a Kubernetes YAML file which can be applied to the actual cluster.  Applying this YAML via e.g. kubectl will deploy and configure the SatCon (razee) agent, which in turn will connect back to SatCon to begin managing the cluster.
1. **Add the cluster to one or more groups.**  If no group exists within the organization (i.e. the IBM Cloud account), you will first need to create a group.  You can then add your cluster to the group. (note *Organization* is the Razee term for the owner of the clusters; in a Satellite Config setting it is always the IBM Cloud account.)
1. **Create a channel.**  To begin defining resources for distribution, you first need a channel.
1. **Add a channel version.**  If you have the correct YAML to deploy a desired Kubernetes resource, you can create a version to encapsulate that YAML within SatCon.  A version is created within the context of a specific channel.
1. **Subscribe a cluster group to a channel/version tuple.**  Creating a subscription, which associates a group with a version (and the owning channel), triggers the agent to download the configuration and apply it to the cluster.

After following these steps, you can then query both SatCon and the k8s API on the cluster itself to see your newly deployed resources.

## Testing

### Running the integration tests

The integration tests are run as a suite of [Ginkgo](https://github.com/onsi/ginkgo) tests.  You will first need to install `ginkgo` (this has probably already been done by `go mod` for you).

Next, navigate to the `test/integration` directory and _copy_ the `integration-sample.json` config file to `integration.json`. Do _not_ change `intergation-sample.json`. Add your credentials to `integration.json`, which is not tracked by git.

- Set `apiKey` to an IAM API key with sufficient permissions.  _Again, be sure not to push any commits that contain actual credentials._
- Set `satconEndpoint` to the Satellite Config API endpoint you want to use.  This is also pre-populated with the production SatCon endpoint.
- Set `orgId` to the IBM Cloud account ID you will use for running the tests. This is generally a 32-character hexadecimal string. At this time, the tests only support using a single orgId/account value for all of the tests.

Then, from that same directory, you can just run `ginkgo .` to execute the integration suite.  This is a much simpler suite than e.g. the [CF Acceptance Tests](https://github.com/cloudfoundry/cf-acceptance-tests), and there is not currently a way to execute only a specific set of tests other than to use the ginkgo-specific focus/pending prefixes within the test code files themselves.

#### [Contribution Guidelines for this project](docs/CONTRIBUTING.md)
