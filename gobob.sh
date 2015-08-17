#!/bin/bash

# this project params
id_params="--name=starter --project=starter"
path_to_gobob="../gobob/gobob"

command="$1"
version="$2"
branch="$3"
subdir="$4"

if [[ $version == "" ]] ; then
  version="dev"
fi
if [[ $branch == "" ]] ; then
  branch="master"
fi
if [[ $command == "" ]] ; then
  command=latest
fi

$path_to_gobob $command $id_params --version=$version --branch=$branch --subdir=$subdir
