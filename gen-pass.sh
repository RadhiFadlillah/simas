#!/bin/bash
# Script to generate password for config file

if [ -s config-pass.go ]
then
	exit 0
fi

password=`date +%s | sha256sum | base64 | head -c 64`

{
	echo 'package main'
	echo
	printf '%s"%s"\n' 'const configPassword = ' "$password"
}>config-pass.go

