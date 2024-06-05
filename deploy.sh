#!/bin/bash

FUNCTION="whatismyip"
ENTRYPOINT="WhatIsMyIP"
SECRETS=""

gcloud functions deploy "${FUNCTION}" \
  --project=cilium-demo \
  --gen2 --runtime=go122 --region=europe-west1 --source=. \
  --entry-point="$ENTRYPOINT" --trigger-http --allow-unauthenticated \
  --set-env-vars GCP_PROJECT=cilium-demo \
  --set-secrets="$SECRETS"
