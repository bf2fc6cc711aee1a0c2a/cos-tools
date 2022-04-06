#!/bin/bash

if [ ! "$#" -eq 1 ]; then
    echo "ID required"
    exit 1
fi

export INSTANCE_ID="${1}"
export INSTANCE_REGION="eu-de-1"

export OCP_FLAVOR="bx2.8x32"
export OCP_VERSION="4.9.21_openshift"

ibmcloud is vpc-create "${INSTANCE_ID}-vpc"
VPC_ID=$(ibmcloud is vpc "${INSTANCE_ID}-vpc" --output JSON | jq -r '.id') 

ibmcloud is public-gateway-create "${INSTANCE_ID}-vpc-gw" "${VPC_ID}" "${INSTANCE_REGION}"
GW_ID=$(ibmcloud is public-gateway "${INSTANCE_ID}-vpc-gw"--output JSON | jq -r '.id')

ibmcloud is subnet-create "${INSTANCE_ID}-snet-1" $VPC_ID --zone $INSTANCE_REGION --ipv4-address-count 256 --pgw $GW_ID
SNET_ID=$(ibmcloud is subnet "${INSTANCE_ID}-snet-1" --output JSON | jq -r '.id')

ibmcloud resource service-instance-create "${INSTANCE_ID}-cos" cloud-object-storage standard global -g Default
COS_ID=$(ibmcloud resource service-instance "${INSTANCE_ID}-cos" --output JSON | jq -r '.id')

ibmcloud oc cluster create vpc-gen2 \
    --name "${INSTANCE_ID}" \
    --zone ${INSTANCE_REGION} \
    --version "${OCP_VERSION}" \
    --flavor "${OCP_FLAVOR}" \
    --workers 2 \
    --vpc-id "${VPC_ID}" \
    --subnet-id "${SNET_ID}" \
    --cos-instance "${COS_ID}"
