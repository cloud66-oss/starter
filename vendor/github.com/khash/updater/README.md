# Updater
 A simple to use Go package for self updating binaries. It supports HTTP download, Semantic versioning, channels and remote forced updates.

 ## Install

```bash
$ go get github.com/khash/updater
```

## Usage

To use updater, you need to push your binaries somewhere they can be downloaded using HTTP (like an S3 bucket). You also need to construct a single JSON file to include details of your versions. By default this file is called `versions.json` and can look like this:

```json
{
  "versions": [
    {
      "version":  "1.0.0-pre",
      "channel": "dev"
    },
    {
      "version": "1.0.0",
      "channel": "stable"
    },
    {
      "version": "1.0.0-pre-1-57dh54",
      "channel": "nightly",
      "force": true
    }
  ]
}
```

Updater supports multiple OS and architectures.

**NOTE**:

Updater strictly requires SemVer compatible versions.

This is an example code:


```go
func update() {
	worker, err := updater.NewUpdater(utils.Version, &updater.Options{
		RemoteURL: "https://s3.amazonaws.com/acme/myapp/",
		Channel:   "dev",
		Silent:    false,
	})
	if err != nil {
		fmt.Println(err)
	}

	err = worker.Run(false)
	if err != nil {
	    fmt.Println(err)
	}
}
```

The code above updates the current executable if the local version is older than the remote version for the `dev` channel.

### Options

`updater.Options` has the following fields:

```go
	RemoteURL            string
	VersionSpecsFilename string
	BinPattern           string
	Channel              string
	Silent               bool
```

**RemoteURL**

The full URL of where the binaries and `versions.json` file can be found. An example could be `https://downloads.acme.org/`. Including the trailing `/` is not mandatory.

**VersionSpecsFilename**

This is the name of the JSON file. If not set, `versions.json`  will be used which will mean the full URL for the JSON file will be `https://downloads.acme.org/versions.json`

**BinPattern**

This is the pattern used to find the relevant binary file to download. If not specified the default is `{{OS}}_{{ARCH}}_{{VERSION}}`. This means the updater will look for version `1.10.20` of the binary compiled for the OSX 64bit architecture at `https://downloads.acme.org/darwin_amd64_1.10.20`

**Channel**

This is the name of the channel to look for. Default is `dev`. Using channels you can choose to have different version tracks for your binaries.

**Silent**

If set to `false` the updater will print out progress of the update to the console (stdout).

### Updater

The Updater itself, can be created using `NewUpdater` and run using the `Run` function. You can force updates by passing the `force` parameter into `Run`. If forced, the binary will be updated even if it's newer than the remote version.

You can also remote force an update. This is useful when you need to rollback to an older version across all clients. To force an update remotely, set the `force` attribute to `true` in `versions.json` as per example above.

## Compiling

Compiling your binaries and construction of the JSON file is up to you. You can use the following as helpers.

### Build

This is an example build bash script:

```bash
#!/bin/bash

version=$(git describe --tags --always)

if [ -z "$1" ]
  then
    echo "No channel supplied"
    exit 1
fi

channel=$1

echo "Building $channel/$version"
echo

rm build/*
curl -s http://s3.amazonaws.com/acme/versions.json | jq '.versions |= map(if (.channel == "'$channel'") then .version = "'$version'" else . end)' > build/versions.json
echo "Current Versions"
cat build/versions.json | jq -r '.versions | map([.channel, .version] | join(": ")) | .[]'
echo

gox -ldflags "-X github.com/acme/myapp/utils.Version=$version -X github.com/acme/myapp/utils.Channel=$channel" -os="darwin linux windows" -arch="amd64" -output "build/{{.OS}}_{{.Arch}}_$version"
```

The above example assumes 2 variables in `utils` package to hold the current version and the current channel for your application and sets them to the current git tag and the user input into the bash script. This means for the updater to work, you'd need to tag your code with a valid SemVer tag.

It also requires [gox](https://github.com/mitchellh/gox) tool to allow cross compiling of the code.

### Publish

To publish your binaries to S3 you can use something like this bash script:


```bash
#!/bin/bash

aws s3 cp build s3://acme/myapp --acl public-read --recursive
```

This script requires a configured AWS cli installed.
