#!/bin/bash

ORG=jenkinsxio
APP_NAME=status-badge

jx step create pr chart --name $ORG/$APP_NAME --version ${VERSION} --repo https://github.com/jenkins-x-apps/jx-app-statusbadge.git
