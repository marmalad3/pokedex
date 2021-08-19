# Pokedex

Tools required
- make
- docker (tested on version 1.12.1)

## Running the server

```
$ make docker-run
# wait for "listening on :5000"

$ curl curl http://localhost:5000/pokemon/oddish
{"description":"During the day,\nit keeps its face\nburied in the\u000cground. At night,\nit wanders around\nsowing its seeds.","habitat":"grassland","isLegendary":false,"name":"Oddish"}
```

### Running the tests

Run the tests using docker
```
$ make docker-run-tests
```

# Technology used
- Docker to ensure the pokedex server runs in the same way regardless of the base environment
- Swagger with [go-swagger](https://goswagger.io/) to generate server models and clients. The `swagger.yaml` file can be provided to downstream services to for documenation, easier integration and to help ensure the API interface is in sync.


# Approach
- Given that there are two APIs to be called, I created a base API client (./internal/app/sources/baseapi) to abstract common features. Clients are created by specifying default information, such as the default API url, in the `GenerateNewAPIClient(..)` method. 
- Upstream API errors are handled in these clients by returning the `baseapi.ApiError{}` type, which can be used to inform retried requests using the `IsRetryable()` method. The `StatusCode` property can also be accessed by casting types if further detail is desired. 
- API retries are not imlpemented in this example project, but my approach would be to use a back-off retry mechanism and a handler timeout if any upstream APIs are taking too long to respond.
- Some upstream APIs are heavily rate-limited. Production use of this server should have authentication support in the API clients. 
- Tests make heavy use of overriding the `http.Client` transport in order to mock the upstream API's response. I feel that this approach is more appropriate than calling the APIs directly in the tests, since these would fail if either the rate limit is hit, or they have uptime issues. This would result in an unreliable test suite.
- Despite using mocks in tests, there is huge value in running full end-to-end tests against  real upstream APIs, to verify that mocks are correct - there could be some undocumented behaviour that needs to be considered. If these APIs were owned by our tech team, I would implement contract testing between APIs in order to verify the interactions are valid on both sides. 
- Both API clients and request handlers are tested using mocks. Injecting API dependecies through the `GetRouter(..)` method allows us to control the upstream API mock interactions from a high level. 
- Translation logic is tested from the HTTP handler perspective. If the logic was more complicated I would move this to some more specific unit tests with more test data.
- Dependencies have been vendored in the project. The main reason for this is to ensure that the build is still successful is a dependency becomes unavailable, from it's source or otherwise.

## Development
Tools required
- make
- go (at least 1.16)
- goimports (for code formatting)

To run tests:

`make test`
