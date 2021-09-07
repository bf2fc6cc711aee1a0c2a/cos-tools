#!/usr/bin/env bash

for cmd in kustomize kubectl curl; do
  if ! command -v $cmd &> /dev/null
  then
      echo "$cmd could not be found"
      exit
  fi
done

kustomize build overlays/staging | kubectl apply -f -