#!/usr/bin/env bash

kubectl get gw | awk '{print $1}' | xargs kubectl delete gw
kubectl get dr | awk '{print $1}' | xargs kubectl delete dr
kubectl get vs | awk '{print $1}' | xargs kubectl delete vs
