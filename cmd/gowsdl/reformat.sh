#!/bin/bash

echo $1


gsed -i  's/type ChangePassword ChangePasswordRequest/type ChangePasswordType ChangePasswordRequest/' netsuite/netsuite.go
gsed -i  's/type ChangeEmail ChangeEmailRequest/type ChangeEmailType ChangeEmailRequest/' netsuite/netsuite.go
gsed -i  's/ChangePassword *ChangePassword /ChangePassword *ChangePasswordType /' netsuite/netsuite.go
gsed -i  's/ChangeEmail *ChangeEmail /ChangeEmail *ChangeEmailType /' netsuite/netsuite.go
gsed -i  's/[^e]Record \*RecordRef/RelatedRecord \*RecordRef/' netsuite/netsuite.go

gsed -i  's/type CampaignResponse string/type CampaignResponseType string/' netsuite/netsuite.go
gsed -i  's/CampaignResponse =/CampaignResponseType =/' netsuite/netsuite.go

# For some reason the template adds *Record to the Address struct but this breaks SOAP so we remove it
gsed -i  '/type Address struct /,/Address/ s/*Record//' netsuite/netsuite.go
gsed -i  '/type CustomFieldList struct /,/CustomFieldList/ s/*CustomFieldRef/any/' netsuite/netsuite.go

cat << EOF >> netsuite/netsuite.go
type TimeBillTimeType string 
const (
    TimeBillTimeType_actualTime TimeBillTimeType = "_actualTime"
    TimeBillTimeType_allocatedTime TimeBillTimeType = "_allocatedTime"
    TimeBillTimeType_plannedTime TimeBillTimeType = "_plannedTime"
)
EOF

gsed -i  's/urn:support_'$1'.lists.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:website_'$1'.lists.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:core_'$1'.platform.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:employees_'$1'.lists.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:faults_'$1'.platform.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:marketing_'$1'.lists.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:common_'$1'.platform.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:accounting_'$1'.lists.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:messages_'$1'.platform.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:bank_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:supplychain_'$1'.lists.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:sales_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:customization_'$1'.setup.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:relationships_'$1'.lists.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:general_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:scheduling_'$1'.activities.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:filecabinet_'$1'.documents.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:communication_'$1'.general.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:financial_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:inventory_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:purchases_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:employees_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:customers_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go
gsed -i  's/urn:demandplanning_'$1'.transactions.webservices.netsuite.com //'  netsuite/netsuite.go