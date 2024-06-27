# Generate Zendesk API Mock
To mock the Zendesk API, we will use the `outofcoffee/imposter`  tool to generate a mock server from the OpenAPI specification. The mock server Dockerfile will be generated in the `mock` directory.
## Generate Dockerfile
First run the following command to generate the mock server Dockerfile:
```shell
docker build -t zendesk-mock .
```
## Run Mock Server locally
Then run the following command to start the mock server locally:
```shell
docker run -it -p 8080:8080 zendesk-mock
```
## Use Mock Server in Tests
Use `testcontainers` to start the mock server in your tests. Here is an example of how to start the mock server in a test:
