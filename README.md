# Cloud 66 Starter

![Codeship Status for cloud66/starter](https://codeship.com/projects/81c5dde0-e914-0133-c219-4eaa3299b296/status)

Cloud 66 starter is an open-source command line tool to generate a `Dockerfile` and a `service.yml` file from arbitrary source code. The `service.yml` file is a Cloud 66 service definition file which is used to define the service configurations on a stack.

To find out more about Starter, checkout the [Cloud 66 Starter website](http://www.startwithdocker.com)

#### Building Starter using Habitus

If you want to contribute to Starter. You can build Starter using [Habitus](http://www.habitus.io), run Habitus in the root directory of this repository. The latest version is generated (after tests) inside the `./artifacts/compiled` directory.

<kbd>$ sudo habitus –host $DOCKER\_HOST –certs $DOCKER\_CERT\_PATH</kbd>

To make sure you a have isolated development environment for contribution. You can use the `docker-compose` for developing, testing and compiling. 

<kbd>$ docker-compose run starter</kbd>

Building starter inside a docker container:

<kbd>root@xx:/usr/local/go/src/github.com/cloud66/starter# go build</kbd>

Running the tests:

<kbd>root@xx:/usr/local/go/src/github.com/cloud66/starter# go test</kbd>