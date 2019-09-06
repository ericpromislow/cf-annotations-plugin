#!/usr/bin/env bash


script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
buildPath="$( cd "${script_dir}/.." && pwd )"

pushd $buildPath
    cf uninstall-plugin RenamerPlugin; go build && cf install-plugin renamer-plugin -f
popd