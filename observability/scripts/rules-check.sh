#!/bin/bash

# This script checks that the rules in the following files are syntactically valid:
#
# connectors-slo-rules-old.yaml

CONTAINER_NAME=prometheus-rules-check-rhoc-observability

docker run -t -p 9090:9090 --name "${CONTAINER_NAME}" \
	-v "${CONNECTORS_SLO_RULES}":/prometheus/connectors-slo-rules-old.yaml:z \
	-v "${CAMELK_OPERATOR_RULES}":/prometheus/camel-k-operator-rules.yaml:z \
	-v "${FLEETSHARD_CAMEL_OPERATOR_RULES}":/prometheus/cos-fleetshard-operator-camel-rules.yaml:z \
	-v "${FLEETSHARD_DEBEZIUM_OPERATOR_RULES}":/prometheus/cos-fleetshard-operator-debezium-rules.yaml:z \
	-v "${FLEETSHARD_SYNC_RULES}":/prometheus/cos-fleetshard-sync-rules.yaml:z \
	-v "${STRIMZI_OPERATOR_RULES}":/prometheus/strimzi-operator-rules.yaml:z \
	--entrypoint=/bin/sh \
	"${IMAGE}" -c '(tail -n +8 connectors-slo-rules-old.yaml ) > connectors-slo-rules-check.yaml \
  && (tail -n +8 camel-k-operator-rules.yaml) > camel-k-operator-rules-check.yaml \
  && (tail -n +8 cos-fleetshard-operator-camel-rules.yaml) > cos-fleetshard-operator-camel-rules-check.yaml \
  && (tail -n +8 cos-fleetshard-operator-debezium-rules.yaml) > cos-fleetshard-operator-debezium-rules-check.yaml \
  && (tail -n +8 cos-fleetshard-sync-rules.yaml) > cos-fleetshard-sync-rules-check.yaml \
  && (tail -n +8 strimzi-operator-rules.yaml) > strimzi-operator-rules-check.yaml \
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
