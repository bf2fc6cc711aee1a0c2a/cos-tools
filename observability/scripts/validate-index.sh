#!/bin/bash

# This script checks that the index.json file is valid JSON.

validate_index(){
	VALIDATE=$(python -mjson.tool $INDEX_FILE_PATH)

	if [[ ! ${VALIDATE} ]]; then
		echo "JSON validation failed for the index.json file due to the above error."
		exit 1
	fi

	echo "SUCCESS: index.json validation check"
}

validate_index