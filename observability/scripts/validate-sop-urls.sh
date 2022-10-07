#!/bin/bash

# Checks that the each SOP URL link is valid, i.e can be reached via curl
# and returns a 2* status code. When running in GitHub actions,
# the ACCESS_TOKEN_SECRET variable must be created within the repository to allow
# the curl command to access the private sops repository.

# The ACCESS_TOKEN_SECRET is a Personal Access Token associated with your account:
# https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token

declare -A BAD_SOP_LINKS
RULE_FILES="${CONNECTORS_SLO_RULES##*/} ${CAMELK_OPERATOR_RULES##*/} ${FLEETSHARD_CAMEL_OPERATOR_RULES##*/} ${FLEETSHARD_DEBEZIUM_OPERATOR_RULES##*/} ${FLEETSHARD_SYNC_RULES##*/} ${STRIMZI_OPERATOR_RULES##*/}"

function check(){
	echo "Validating the SOP URL for alerts in the following files: "$RULE_FILES" ..."

	readarray ALL_SOPS < <(docker run --rm -v "${PROMETHEUS_RULES_DIR}":/workdir:z mikefarah/yq -N e '.spec.groups[]
	| select(.name != "deadmanssnitch")
	| select(.name != "rule-evaluation")
	| .rules[].annotations
	| select(length!=0)
	| .sop_url ' "${CONNECTORS_SLO_RULES##*/}" "${CAMELK_OPERATOR_RULES##*/}" "${FLEETSHARD_CAMEL_OPERATOR_RULES##*/}" "${FLEETSHARD_DEBEZIUM_OPERATOR_RULES##*/}" "${FLEETSHARD_SYNC_RULES##*/}" "${STRIMZI_OPERATOR_RULES##*/}" | sort -u)

	for SOP in "${ALL_SOPS[@]}"; do
		STATUS_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" https://raw.githubusercontent.com/${RHOC_SOPS_REPO_ORG}/cos-tools/main/observability/sops/alerts/${SOP##*/})
		if [[ "$STATUS_RESPONSE" != "2"* ]]; then
		    echo
			echo " Checking: ${SOP##*/} Status Code: $STATUS_RESPONSE"
			BAD_SOP_LINKS[$SOP]=$STATUS_RESPONSE
		else
			echo " Checking: ${SOP##*/} Status Code: $STATUS_RESPONSE"
			echo
		fi
	done

	if  [ ${#BAD_SOP_LINKS[@]} -gt 0 ]; then
		echo "The following SOP URL(s) are invalid, in the wrong folder, or could not be reached:"
		for k in "${!BAD_SOP_LINKS[@]}"; do
			printf "SOP URL Link: $k\nStatus Code: ${BAD_SOP_LINKS[$k]}\n"
		done
		exit 1
	else
		echo
		echo -e " SUCCESS: All SOP URL links are valid"
	fi
}

check
