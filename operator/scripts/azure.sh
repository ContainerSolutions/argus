#! /bin/bash

set -euo pipefail
az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET --tenant $AZURE_TENANT_ID >> /dev/null
state=$(az policy state list --resource-group argus --query "[?policyDefinitionId=='/providers/microsoft.authorization/policydefinitions/e8eef0a8-67cf-4eb4-9386-14b0e78733d4']" | jq -r '.[].complianceState')

if [[ $state == "NonCompliant" ]]; then
    echo "Policy state is currently non compliant"
    exit 1
else
    echo "Policy state is compliant"
    exit 0
fi