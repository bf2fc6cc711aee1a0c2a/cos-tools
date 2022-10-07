#!/bin/bash

# This script checks that the rules in the following files are syntactically valid:
#
# connectors-slo-rules.yaml

CONTAINER_NAME=prometheus-rules-check-rhoc-observability

docker run -t -p 9090:9090 --name "${CONTAINER_NAME}" \
	-v "${CONNECTORS_SLO_RULES}":/prometheus/connectors-slo-rules.yaml:z \
	-v "${CAMELK_OPERATOR_RULES}":/prometheus/camel-k-operator-rules.yaml:z \
	--entrypoint=/bin/sh \
	"${IMAGE}" -c '(tail -n +8 connectors-slo-rules.yaml ) > connectors-slo-rules-check.yaml \
  && (tail -n +8 camel-k-operator-rules.yaml) > camel-k-operator-rules-check.yaml \
	&& promtool check rules *-rules-check.yaml' \
	docker stop "${CONTAINER_NAME}"

if  docker logs -f ${CONTAINER_NAME} | grep -q FAILED:; then
	docker logs -f ${CONTAINER_NAME}
	docker rm "${CONTAINER_NAME}"
	exit 1
else
	echo "SUCCESS: rules check"
	docker rm "${CONTAINER_NAME}"
fi
