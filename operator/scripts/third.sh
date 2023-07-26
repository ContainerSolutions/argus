#! /bin/bash

set -euo pipefail
az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET --tenant $AZURE_TENANT_ID >> /dev/null
state=$(az policy state list --resource "/subscriptions/2f2d0095-02f7-42ee-bd27-9559b16fcd21/resourcegroups/argus/providers/microsoft.containerregistry/registries/argusthird" --query "[?policyDefinitionId=='/providers/microsoft.authorization/policydefinitions/d0793b48-0edc-4296-a390-4c75d1bdfd71']" | jq -r '.[].complianceState')

if [[ $state == "Compliant" ]]; then
    echo "Policy state is compliant"
    exit 0
else
    echo "Policy state is non compliant"
    exit 1
fi
