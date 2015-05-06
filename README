# Cloud 66 Starter
==================

Cloud 66 starter is an open source command line tool to generate Dockerfile and `service.yml` (Cloud 66 service definition file for docker containers) from source code.

It works in the same way as BuildPacks but only to generate the above mentioned files so the compile step happens on the [BuildGrid](http://help.cloud66.com/building-your-stack/introduction-to-docker-deployments). It also does not have any third party requirements or frameworks to work (it's compiled as a Go executable).

## Get Started
==============
To get started download the executable and run it on your development machine.

```
$ cd /my/project
$ starter 
```

This will analyse the project in the current folder and generates two files: `Dockerfile` and `services.yml` in the same folder. 

```
Cloud 66 Starter - (c) 2015 Cloud 66 
 Detecting framework for the project at /Users/awesome/work/boom
 Found ruby application
 ----> Found non Webrick application server (%s) 
 ----> Found Mysql 
 ----> Found Redis 
 Writing Dockerfile...
 Parsing Procfile
 ----> Found Procfile item faye
 ----> Found Procfile item scheduler
 ----> Found Procfile item worker_high
 ----> Found Procfile item worker_low
 ----> Found Procfile item worker_background
 ----> Found Procfile item worker_mem
 Writing service.yml...
 
 Done 
```

starter supports Procfiles and generates a service in `service.yml` for each item in the `Procfile`.

To use starter on a different folder, you can use the `p` option:

```
$ starter -p /my/project 
```

By default starter does not overwrite the existing `Dockerfile` or `service.yml` files. To force it to do so you can use the `o` option:

```
$ starter -p /my/project -o
```

### Supported Languages / Frameworks
====================================

- Ruby, Rack (Rails, Sinatra, Padrino)

## Adding support for new frameworks and languages
==================================================

We will be adding support for new languages and frameworks over the time. However, if you find yourself interested in adding one, it's fairly easy to do:

- Get the source code and compile it (it's written in Go!)
- Add a file under the `packs` folder with the name of the language you are adding, for example `java.go`.
- The `Java` struct should implement `packs.Pack` interface methods. The main ones are `Detect` and `Compile`
- `Detect` tells starter if a folder contains a project in the applicable language (in this example Java)
- `Compile` returns a `common.ParseContext` which is used to render a `Dockerfile` and `services.yml`
- Create a template with the name of the language under the `templates` folder. For example `java.dockerfile.template`
- Use Golang template syntax to build the template for `Dockerfile`
