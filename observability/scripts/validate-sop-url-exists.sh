#!/bin/bash

# This script checks that all of the alerts in the files specified by
# the following environment variables have a corresponding
# SOP_URL. The output is a list of alerts without a SOP_URL:
#
# PROMETHEUS_RULE_FILE
# PROMETHEUS_DEVELOPER_RULE_FILE

# The deadmanssnitch and rule-evaluation groups are omitted
# as the alert within each group does not require a SOP

failed=0

function check() {
    local prometheus_rules_cr=${1}

    readarray ALL_ALERTS < <(docker run --rm -v "${PWD}":/workdir:z mikefarah/yq e '.spec.groups[]
      | select(.name != "deadmanssnitch")
      | select(.name != "rule-evaluation")
      | .rules[].annotations
      | select(length!=0)
      | parent
      | .alert' "${prometheus_rules_cr}")

    readarray ALERTS_WITH_SOPS < <(docker run --rm -v "${PWD}":/workdir:z mikefarah/yq e '.spec.groups[]
      | select(.name != "deadmanssnitch")
      | select(.name != "rule-evaluation")
      | .rules[].annotations
      | select(length!=0)
      | select(.sop_url)
      | parent
      | .alert' "${prometheus_rules_cr}")

    NUM_ALERTS_WITHOUT_SOPS=$((${#ALL_ALERTS[@]}-${#ALERTS_WITH_SOPS[@]}))

    if  [ $NUM_ALERTS_WITHOUT_SOPS -gt 0 ]; then
        echo -e "The following $NUM_ALERTS_WITHOUT_SOPS alert(s) in the $(basename "${prometheus_rules_cr}") file do not have a corresponding SOP:\n"
        echo "${ALL_ALERTS[@]}" "${ALERTS_WITH_SOPS[@]}" | tr ' ' '\n' | sort | uniq -u
        echo
        failed=1;
    fi
}

check "$(realpath --relative-to . "${CONNECTORS_SLO_RULES}")"
check "$(realpath --relative-to . "${CAMELK_OPERATOR_RULES}")"
check "$(realpath --relative-to . "${FLEETSHARD_CAMEL_OPERATOR_RULES}")"
check "$(realpath --relative-to . "${FLEETSHARD_DEBEZIUM_OPERATOR_RULES}")"
check "$(realpath --relative-to . "${FLEETSHARD_SYNC_RULES}")"
check "$(realpath --relative-to . "${STRIMZI_OPERATOR_RULES}")"

if [ ${failed} -ne 0 ]; then
    exit 1;
else
    echo "SUCCESS: SOP URL check"
fi