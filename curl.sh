#!/bin/bash


curl http://localhost:4355/api/v1/getparams.execute -H "Authorization: Bearer supersecret" -d \
'{
  "applicationSetName": "fb-matrix",
  "input": {
    "parameters": {
      "labels": "[advertiser no-preview]",
      "excludeLabel": "no-preview",
      "blackHole": "/foo/bar",
      "path": "/foo/real",
      "number": "123"

    }
  }
}'