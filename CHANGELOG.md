# CHANGELOG

## V1.3.0

This release will support the direct convertion from a `service.yml` to the `kubernetes.yml` format. This give you the option to run a Cloud 66 service definition on your local Kubernetes cluster (for example minikube). The `kubernetes` pack will be available to be used both through the API and the CLI tool.

## V1.2.0

This release will support the direct convertion from a `docker-compose.yml` to the `service.yml` format. The `docker-compose` pack will be available to be used both through the API and the CLI tool.

## V1.1.0

This release will support API mode. You can use Starter as a microservice to analyse codebases on demand. Also, this version will use Docker Registry to check if a specific version public version exists. We add Meteor support to the Node buildpack. The binaries are compiled using Go 1.7 (compatible with macOS Sierra)

## V1.0.3

Support for nodejs and PHP

## V1.0.2

Includes support for compose.yml generation

## V1.0.1

Bump up version

## V1.0.0 - Initial Alpha

This is the first release of Cloud 66 Starter. This is an alpha release, however running it does no harm as it only writes a Dockerfile and service.yml in the project directory. You can manually modify the files and commit them into your repository once you're happy with them.
