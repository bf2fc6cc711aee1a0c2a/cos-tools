#!/bin/bash

# Checks to ensure that all warning and critical
# alerts have a corresponding unit test.

function check() {
	SEVERITY=$1

	while IFS= read -r ALERT; do ALL_ALERTS+=($ALERT)
	done < <(docker run --rm -v "${PROMETHEUS_RULES_DIR}":/workdir:z mikefarah/yq -N e '.spec.groups[]
	| select(.name != "deadmanssnitch")
	| .rules[].labels
	| select(length!=0)
	| select(.severity == '$SEVERITY')
	| parent
	| .alert' "${CONNECTORS_SLO_RULES##*/}" | sort -u)

	while IFS= read -r ALERT; do PROMETHEUS_UNIT_TEST_CHECK+=($ALERT)
	done < <(yq -N eval-all '.tests[]
	| .alert_rule_test[]
	| .exp_alerts[]
	| .exp_labels
	| select(.severity == '$SEVERITY')
	| .alertname' $UNIT_TEST_FILES* | sort -u)

	ALERTS_WITHOUT_UNIT_TESTS=(`echo ${ALL_ALERTS[@]} ${PROMETHEUS_UNIT_TEST_CHECK[@]} | tr ' ' '\n' | sort | uniq -u `)

	if  [ ${#ALERTS_WITHOUT_UNIT_TESTS[@]} -gt 0 ]; then
		printf "FAILURE: ${#ALERTS_WITHOUT_UNIT_TESTS[@]} $SEVERITY alert(s) are missing a corresponding unit-test:\n" && echo
		printf '%s\n' "${ALERTS_WITHOUT_UNIT_TESTS[@]}"
		exit 1
	else
		echo -e "SUCCESS: All $SEVERITY alerts have a corresponding unit-test"
	fi
}

check ${CRITICAL_SEVERITY}
# <TO-DO> https://issues.redhat.com/browse/MGDSTRM-6845
# check ${WARNING_SEVERITY}