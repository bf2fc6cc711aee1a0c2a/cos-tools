UNIT_TEST_DIR ?= $(shell pwd)/resources/prometheus/unit_tests/
export PROMETHEUS_RULES_DIR ?= $(shell pwd)/resources/prometheus/
export CONNECTORS_SLO_RULES ?= $(shell pwd)/resources/prometheus/connectors-slo-rules.yaml
export CAMELK_OPERATOR_RULES ?= $(shell pwd)/resources/prometheus/camel-k-operator-rules.yaml
export FLEETSHARD_CAMEL_OPERATOR_RULES ?= $(shell pwd)/resources/prometheus/cos-fleetshard-operator-camel-rules.yaml
export FLEETSHARD_DEBEZIUM_OPERATOR_RULES ?= $(shell pwd)/resources/prometheus/cos-fleetshard-operator-debezium-rules.yaml
export FLEETSHARD_SYNC_RULES ?= $(shell pwd)/resources/prometheus/cos-fleetshard-sync-rules.yaml
export STRIMZI_OPERATOR_RULES ?= $(shell pwd)/resources/prometheus/strimzi-operator-rules.yaml
export IMAGE ?= quay.io/prometheus/prometheus
export DASHBOARDS_DIR ?= $(shell pwd)/resources/grafana/
export INDEX_FILE_PATH ?= $(shell pwd)/resources/index.json
export UNIT_TEST_FILES ?= $(shell pwd)/resources/prometheus/unit_tests/
export CRITICAL_SEVERITY="critical"
export WARNING_SEVERITY="warning"
export RHOC_SOPS_REPO_ORG ?= bf2fc6cc711aee1a0c2a
export FILES_TO_UNIT_TEST ?= *

# Run unit unit_tests files located in the unit_tests directory.
# By default, all unit unit_tests are executed. If you want to
# test a single unit test, export the name of the test file:
# e.g export FILES_TO_UNIT_TEST=OfflinePartitions.yaml
.PHONY: unit/test
unit/test:
	docker run --rm -t \
    -v $(CONNECTORS_SLO_RULES):/prometheus/connectors-slo-rules.yaml:z \
    -v $(CAMELK_OPERATOR_RULES):/prometheus/camel-k-operator-rules.yaml:z \
    -v $(FLEETSHARD_CAMEL_OPERATOR_RULES):/prometheus/cos-fleetshard-operator-camel-rules.yaml:z \
    -v $(FLEETSHARD_DEBEZIUM_OPERATOR_RULES):/prometheus/cos-fleetshard-operator-debezium-rules.yaml:z \
    -v $(FLEETSHARD_SYNC_RULES):/prometheus/cos-fleetshard-sync-rules.yaml:z \
    -v $(STRIMZI_OPERATOR_RULES):/prometheus/strimzi-operator-rules.yaml:z \
    -v $(UNIT_TEST_DIR):/prometheus/unit_tests:z --entrypoint=/bin/sh \
$(IMAGE) -c '(tail -n +8 connectors-slo-rules.yaml; tail -n +9 camel-k-operator-rules.yaml; \
			  tail -n +9 cos-fleetshard-operator-debezium-rules.yaml; tail -n +9 cos-fleetshard-operator-camel-rules.yaml; \
			  tail -n +9 cos-fleetshard-sync-rules.yaml; tail -n +9 strimzi-operator-rules.yaml) \
              > rules.yaml && cd unit_tests && promtool test rules ${FILES_TO_UNIT_TEST}'

# Checks the prometheus rules in the given rules files
.PHONY: check/rules
check/rules:
	./scripts/rules-check.sh

# Check that each dashboard is valid JSON
.PHONY: validate/dashboards
validate/dashboards:$(shell pwd)
	./scripts/validate-json.sh

# Check that the index file is valid JSON
.PHONY: validate/index
validate/index:$(shell pwd)
	./scripts/validate-index.sh

# Check each alert has a valid unit test
.PHONY: check/unit-tests
check/unit-tests:$(shell pwd)
	./scripts/unit-test-check.sh

# Check each alert has a SOP
.PHONY: alerts/sop_url_exists
alerts/sop_url_exists:$(shell pwd)
	./scripts/validate-sop-url-exists.sh

.PHONY: validate/sop_url_links
validate/sop_url_links:$(shell pwd)
	./scripts/validate-sop-urls.sh

# Run all test targets
.PHONY: run/tests
run/tests: unit/test alerts/sop_url_exists validate/sop_url_links validate/dashboards validate/index check/rules check/unit-tests
