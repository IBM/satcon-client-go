## 0.3.0 16 June 2022

### **BREAKING CHANGES**
- Remove HTTPClient arg from client.New(). (#43) [61152f7]

  95% of the time people were just passing nil to this. If you want to pass in a custom HTTPClient, use the new `NewWithCustomHTTPClient()` function.
- Add `SubscriptionIdsForCluster()` method. (#42) [9d4f411]

  This would only break something if you have your own implementation of `subscriptions.SubscriptionService`.
### New Features
- Add `SubscriptionIdsForCluster()` method. (#42) [9d4f411]
- Add `Children` field to `SearchableData` struct. (#44) [bef4328]
### Dependency major version upgrades
- Update IBM/go-sdk-core dependency to v5 (#48) [4f54f37]
- Update golang-jwt to v4 (#47) [8b15ccc]

## 0.2.0 15 July 2021

### New Features
- Functions that query for groups now return more detailed information about the clusters in the groups [95a8bcb]
- Allow passing alternate url for IAM endpoint (#31) [ece05a6]

### Fixes
- All characters properly handled when converting to JSON (#39) [fb58ff1]

### Documentation
- Update README to address Satellite Config vs. Razee terminology confusion (#32) [56442d3]

### Stand-alone Razee features (EXPERIMENTAL)
- Support for authentication using an apikey of an OAuth service (#29) [4633884]
- Support for local user/password authentication (#24) [49a390b]

### For contributors
- More flexibility in configuration(s) for running integration tests [b8fa1c7, ece05a6]
- Add guards against accidentally pushing a real api key (for integration tests) (#34) [0caf65b]

## 0.1.1

- Initial Public Release