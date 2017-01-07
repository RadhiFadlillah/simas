#!/bin/bash
# Script to generate password for config file

password=`date +%s | sha256sum | base64 | head -c 64`

{
	echo 'package main'
	echo
	printf '%s"%s"\n' 'const configPassword = ' "$password"
}>config-pass.go

