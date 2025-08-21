#!/bin/bash

FUNCTION="whatismyip"
ENTRYPOINT="WhatIsMyIP"
SECRETS="BASIC_AUTH=whatismyip-auth:latest"
GCP_PROJECT="gcp-isovalentmarket-nprd-33910"

gcloud functions deploy "${FUNCTION}" \
  --project="${GCP_PROJECT}" \
  --gen2 --runtime=go122 --region=europe-west1 --source=. \
  --entry-point="$ENTRYPOINT" --trigger-http --allow-unauthenticated \
  --set-env-vars GCP_PROJECT="${GCP_PROJECT}" \
  --set-secrets="$SECRETS"
