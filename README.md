<img src="http://cdn2-cloud66-com.s3.amazonaws.com/images/oss-sponsorship.png" width=150/>

# Starter

![Codeship Status for cloud66/starter](https://codeship.com/projects/81c5dde0-e914-0133-c219-4eaa3299b296/status)

Starter is an open source command line tool to generate a `Dockerfile` and a `docker-compose.yml` file from arbitrary source code. It will kickstart the journey towards containerizing your apps.

Starter can also generate the `service.yml` file, which is a Cloud 66 service definition file used to define the service configurations that run Docker in production on any cloud, or on your own serve

![Logo Starter an open source dockerfile generator](http://blog.cloud66.com/content/images/2016/08/Starter-open-source-dockerfile-generator-on-github.png)

- Website: http://www.startwithdocker.com/
- Download Starter: https://github.com/cloud66-oss/starter/releases/
- Articles: http://blog.cloud66.com/tag/starter/


### Key features:
___

- Detects **frameworks** and its version, i.e. Rails v5.0, PHP Laravel and Node.JS express to name a few.
- Determines the **ports** that need setting.
- Detects what **database** you’re using, to setup your databases in docker-compose.
- Compatible with **Procfiles** to generate services for you.
- Examines your application to generate appropriate `Dockerfile`, `docker-compose.yml`, and `service.yml` files.
- Has an **API** to integrate Starter into you own product.


### Why Starter?
___

- **You’re new to Docker, have got the basic 101 content and are now ready to start using Docker:**
  Starter is an ideal tool to support you with your first Docker deployment. It’s the easiest way to learn how to containerize your application, and is a great step to ease you through the Docker learning curve. It automates the process, allowing you to focus on the things that matter.

- **You’re in need of containerizing your existing application:**
  Starter helps you analyze your existing application and detects what framework the application is running and in what version. Additionally, it automatically detects what database and ports you’re using.

  Next it will generate a Dockerfile, DockerCompose or service.yml that is ready to run in containers. This helps you with faster builds and prepares you to run Docker in production.

**Reasons to containerize your applications**

- If you need to achieve multi-tenancy
There are a couple reasons why you would need to run multiple applications on the same stack, i.e. the applications that you’re running share common resources or your applications don’t receive enough traffic to run on separate stacks. This is where Starter can help to Dockerize each of your applications, which you can host on a single stack.

- If you have a specialized application that requires a sandbox environment
With the new security features of Docker, you get true isolation of your process. When using Starter you can isolate parts of your apps in containers and make sure they’re in a sandbox (and can’t do harmful things!).

### Documentations:
___

Comprehensive documentation is available on the Starter website:

http://www.startwithdocker.com/

### Quick Start:
___

Head to the Starter releases (https://github.com/cloud66-oss/starter/releases/latest) and download the latest version for your platform. You can copy the file to `/usr/local/bin` and make sure it is renamed to `starter` and you can run it (`chmod a+x /usr/local/bin/starter`). From this point on you can run `starter update` to update it automatically.

    $ cd /my/project
    $ starter -g dockerfile,service,docker-compose

This will analyze the project in the current folder and generate the three files: `Dockerfile, docker-compose.yml and `service.yml` in the same folder, prompting for information when required.


    Cloud 66 Starter ~ (c) 2019 Cloud 66
    Detecting framework for the project at /Users/awesome/work/boom
    Found ruby application
    Enter ruby version: [latest]
    ----> Found config/database.yml
    Found mysql, confirm? [Y/n]
    Found redis, confirm? [Y/n]
    Found elasticsearch, confirm? [Y/n]
    Add any other databases? [y/N]
    ----> Analyzing dependencies
    ----> Parsing Procfile
    ----> Found Procfile item web
    ----> Found Procfile item worker
    ----> Found unicorn
    This command will be run after each build: '/bin/sh -c "RAILS_ENV=_env:RAILS_ENV bundle exec rake db:schema:load"', confirm? [Y/n]
    This command will be run after each deployment: '/bin/sh -c "RAILS_ENV=_env:RAILS_ENV bundle exec rake db:migrate"', confirm? [Y/n]
    ----> Writing Dockerfile...
    ----> Writing docker-compose.yml...
    ----> Writing service.yml
    Done

Starter supports Procfiles and generates a service in `service.yml` for each item in the Procfile. It is highly advised to use a Procfile to define your own service commands as starter will only detect the web service otherwise.

To use starter on a different folder, you can use the `p` option:


    $ starter -p /my/project

For more options, please see:


    $ starter help


### Building Starter using Habitus
___


If you want to contribute to Starter. You can build Starter using [Habitus](http://www.habitus.io). Habitus is an open source build flow tool for Docker.

Run Habitus in the root directory of this repository. The latest version is generated (after tests) inside the `./artifacts/compiled` directory.

<kbd>habitus --keep-artifacts=true</kbd>

To make sure you a have isolated development environment for contribution. You can use the `docker-compose` for developing, testing and compiling.

<kbd>$ docker-compose run starter</kbd>

Building starter inside a docker container:

<kbd>root@xx:/usr/local/go/src/github.com/cloud66/starter# go build</kbd>

Running the tests:

<kbd>root@xx:/usr/local/go/src/github.com/cloud66/starter# go test</kbd>


And you’re ready to start contributing to Starter.
