#!/usr/bin/env bash

CHAINCODE_NAME=faceid
CHAINCODE_VERSION=1.9
CHANNEL_NAME=studychain

#set -x

KUBECTL_VERSION=v2
if [ "$1" == "v1" ]; then
    KUBECTL_VERSION=v1
fi

CHAINCDE_OPTION=init
if [ "$2" == "upgrade" ]; then
    CHAINCDE_OPTION=upgrade
fi


function applyBy(){
    cat ./.kubernetes/$1 \
    | sed "s/__CHAINCODE_NAME__/${CHAINCODE_NAME}/g" \
    | sed "s/__CHAINCODE_VERSION__/${CHAINCODE_VERSION}/g" \
    | sed "s/__CHANNEL_NAME__/${CHANNEL_NAME}/g" \
    | kubectl apply --record -f -
}

function getPodsStatus() {
    podStatus=Error
    if [ "${KUBECTL_VERSION}" == "v2" ]; then
        podStatus=$(kubectl get pods -n servicechain-fabric | grep "$1" | awk '{print $3}')
    else
        podStatus=$(kubectl get pods -n servicechain-fabric | grep "$1" | awk '{print $4}')
    fi
    if [ "${podStatus}" == "Error" ]; then
        echo "There is an error in $1 job. Please check logs."
        exit 1
    fi
}

function getJobStatus() {
    if [ "${KUBECTL_VERSION}" == "v2" ]; then
        jobStatus=$(kubectl get jobs -n servicechain-fabric |grep "$1" | awk '{print $2}')
        if [ ${jobStatus} == "1/1" ]; then
            return 0
        fi
        return 1
    else
        jobStatus=$(kubectl get jobs -n servicechain-fabric |grep "$1" | awk '{print $3}')
        if [ ${jobStatus} == "1" ]; then
            return 0
        fi
        return 1
    fi
}

function checkJob() {
    job=$1
    getJobStatus ${job}
    while [ $? != "0" ]; do
        _pod=$(kubectl get pods -n servicechain-fabric --selector=job-name=$1 |sed '1d')
        echo "Waiting for ${job} job to be completed at -> ${_pod}"
        sleep 1;
        getPodsStatus ${job}
        getJobStatus ${job}
    done
    echo "Job ${job} Completed Successfully"
}

kubectl get jobs -n servicechain-fabric | sed '1d' | awk '{ print $1 }' | while read line; do
    kubectl delete jobs ${line} -n servicechain-fabric
done


applyBy 0copy-chaincode.yaml

pod=$(kubectl get pods -n servicechain-fabric --selector=job-name=copychaincode --output=jsonpath={.items..metadata.name})
podSTATUS=$(kubectl get pods -n servicechain-fabric --selector=job-name=copychaincode --output=jsonpath={.items..phase})
while [ "${podSTATUS}" != "Running" ]; do
  echo "Wating for container of delete pod to run. Current status of ${pod} is ${podSTATUS}"
  sleep 5;
  if [ "${podSTATUS}" == "Error" ]; then
    echo "There is an error in copychaincode job. Please check logs."
    exit 1
  fi
  pod=$(kubectl get pods -n servicechain-fabric --selector=job-name=copychaincode --output=jsonpath={.items..metadata.name})
  podSTATUS=$(kubectl get pods -n servicechain-fabric --selector=job-name=copychaincode --output=jsonpath={.items..phase})
done

echo "start copy ${CHAINCODE_NAME} to ${pod}"
kubectl cp -n servicechain-fabric ./chaincode/${CHAINCODE_NAME} $pod:/shared/artifacts/chaincode/
echo "copy ${CHAINCODE_NAME} end"

echo "Waiting for 10 more seconds for copying chaincode ${CHAINCODE_NAME} to avoid any network delay"
sleep 10


checkJob copychaincode
echo "Copy chaincode job completed"

echo -e "\nCreating installchaincode job"
applyBy 1install-chaincode.yaml
checkJob chaincodeinstall
echo "Chaincode Install Completed Successfully"

if [ ${CHAINCDE_OPTION} == "init" ]; then
    echo -e "\nCreating chaincodeinstantiate job"
    applyBy 2instantiate-chaincode.yaml
    checkJob chaincodeinstantiate
    echo "Chaincode Instantiation Completed Successfully"
else
    echo -e "\nCreating chaincodeupgrade job"
    applyBy 3upgrade-chaincode.yaml
    checkJob chaincodeupgrade
    echo "Chaincode Upgrade Completed Successfully"
fi
