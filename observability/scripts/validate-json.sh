#!/bin/bash

# This script checks that each Grafana dashboard contains valid JSON.

validate_json(){
	for DASHBOARD in $DASHBOARDS_DIR*;
	do
		let LINE_NUMBER=($(awk '/json:/{ print NR }' "$DASHBOARD")+1)
		VALIDATE=$(tail --lines=+$LINE_NUMBER $DASHBOARD | python -mjson.tool)

		if [[ ! ${VALIDATE} ]]; then
			echo "JSON validation failed for dashboard `basename $DASHBOARD` due to the above error.
			The line number of the error corresponds to the general.json spec of the dashboard, and doesn't take into account the YAML definition at the top of the file"
			exit 1
		fi
	done

	echo "SUCCESS: dashboard validation check"
}

validate_json