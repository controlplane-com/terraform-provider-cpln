#!/bin/bash

echo -n "" > ./terraform_reference.mdx


PROVIDER_TEXT=$(<"../docs/index.md")
PROVIDER=$(echo "$PROVIDER_TEXT" | awk 'NR>9' )
CLEANPROVIDER=$(echo "$PROVIDER" | sed 's/~> //g;' | sed 's/## Authentication/## Authentication {} --exclude/g;' | sed 's/(Resource)//g;' | sed 's/## Declaration/## Declaration {} --exclude/g;' | sed 's/### Required/### Required {} --exclude/g;' | sed 's/### Optional/### Optional {} --exclude/g;' | sed 's/## Outputs/## Outputs {} --exclude/g;' | sed 's/## Example Usage/## Example Usage {} --exclude/g;'  | sed 's/<a id=.*//g;')


cat << EOF >> ./terraform_reference.mdx
# Terraform Provider {#terraform-provider}

$CLEANPROVIDER

EOF

for FILE in ../docs/resources/*.md; do 


CURRENT_TEXT=$(<$FILE)
TEXT=$(echo "$CURRENT_TEXT" | awk 'NR>7' )
NORMALIZETEXT=$(echo "$TEXT" | sed -E 's/([a-z])--([a-z])/\1.\2/g;' )
CLEANTEXT=$(echo "$NORMALIZETEXT" | sed 's/~> //g;' | sed 's/(Resource)//g;' | sed 's/## Declaration/## Declaration {} --exclude/g;' | sed 's/### Required/### Required {} --exclude/g;' | sed 's/### Optional/### Optional {} --exclude/g;' | sed 's/## Outputs/## Outputs {} --exclude/g;' | sed 's/## Example Usage/## Example Usage {} --exclude/g;'  | sed 's/<a id=.*//g;')
LINKTEXT=$(echo "$CLEANTEXT" | sed -E 's/^# cpln_(.*).$/# cpln_\1 {#cpln_\1}/g;' | sed -E 's/^ ### `(.*)`$/ ### \1 {#nestedblock.\1} --exclude/g;')

cat << EOF >> ./terraform_reference.mdx
$LINKTEXT

EOF


done