# satcon-client

## Testing

### Running the integration tests

The integration tests are run as a suite of [Ginkgo](https://github.com/onsi/ginkgo) tests.  You will first need to install `ginkgo` (this has probably already been done by `go mod` for you).

Next, navigate to the `test/integration` directory and update the `integration.json` config file.

- Set `apiKey` to an IAM API key with sufficient permissions.  _Be sure not to push any commits that contain actual credentials._
- Set `iamEndpoint` to the IAM token endpoint you want to use.  This is prepopulated with the staging IAM token endpoint, but you can change it to the production one if you want.  (_Note: this will not work pre-SatCon beta._)
- Set `satconEndpoint` to the Satellite Config API endpoint you want to use.  This is also prepopulated with the staging SatCon endpoint.
- Set `orgId` to the Satellite organization ID you will use for running the tests.  At this time, the tests only support using a single orgId value for all of the tests.

Then, from that same directory, you can just run `ginkgo .` to execute the integration suite.  This is a much simpler suite than e.g. the [CF Acceptance Tests](https://github.com/cloudfoundry/cf-acceptance-tests), and there is not currently a way to execute only a specific set of tests other than to use the ginkgo-specific focus/pending prefixes within the test code files themselves. 
