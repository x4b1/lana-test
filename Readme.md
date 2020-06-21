
# Lana challenge

## Requirements

- Docker
- Docker compose

## Available commands

- `make test` : Run all tests unit and functional inside docker container
- `make run` : Runs the api (service) and the testing client

## Testing

### Local

You can use the provided client but also the docker container exposes an endpoint to test it:
`http://localhost:3000`

### Production

I have deployed the service to heroku (https://xmartinez-lana-api.herokuapp.com/) so you can also test it there using the following collection

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/f12d127122248f1bc16b)

## Comments

- When I initialize the DB I add the products and discounts as it is not a requirement the manage of them.
- I have added 2 middlewares to the API, one to add extra data of the request to the context and a second one to log all the requests to the api.
- When there is a new push to master the app automatically runs the linter and test, if both success it deploys to heroku.

## Architecture

Code is divided in the following structure

- `cmd/`: Main applications for this project.
  - `api/`: Entry point of the API
  - `cli/`: Entry point of the testing client
- `internal/`: Where application resides.
  - `/` Domain of the aplication
    - `creating/`: application service. Creates a new basket in system
    - `adding/`: application service. Given a basket and product checks if both exists in the system adds a new item to basket
    - `calculating/`: application service. Given a basket checks if exists in the system calculates the total of the basket applying discounts
    - `deletinf/`: application service. Given a basket checks if exists in the system deletes it.
    - `server/`: Here is defined the API handler and responses. Is responsible to get the request and process and return a response.
    - `storage/`: Here is defined the multiple repositories for the different storage types. In this case only memory.
- `pkg/`: Packages that could be in other repository to be used by other services
  - `errors/` errors types of the application, I only have defined 2 type of error because we dont have more.
  - `log/`: Log services, here we can define multiple log services, they should implement logger interface to allow change between them by dependency injection
  - `money/`: Implementation of the money patern to handle al money transactions.
  - `uuid/`: unique ID generator
- `./github`: Here is defined the flows that will follow github runners.
