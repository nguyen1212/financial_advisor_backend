#!/bin/bash

RED=$(tput setaf 1)
GRN=$(tput setaf 2)
PUR=$(tput setaf 13)
RESET=$(tput sgr0)
TEST_RESULT_DIR="${TEST_RESULTS:-./test-results}"
mkdir -p ${TEST_RESULT_DIR}
echo "${GRN}Listing all packages${RESET}"
PKG_LIST+="$(go list ./... | grep -v /vendor/ | grep -v /migration/ | grep -v /mock) "
echo "----------------"
echo "${GRN}Test:${RESET}"

# Set environment variables for webhook_trigger lambda tests
export AWS_REGION=${AWS_REGION:-ap-northeast-1}
export AWS_ENDPOINT=${AWS_ENDPOINT:-http://localhost:4566}
export AWS_S3_BUCKET=${AWS_S3_BUCKET:-test-bucket}
export AWS_SQS_QUEUE_NAME_WEBHOOK=${AWS_SQS_QUEUE_NAME_WEBHOOK:-test-queue}
export DD_ENV=${DD_ENV:-test}
export LOG_LEVEL=${LOG_LEVEL:-debug}
export ID_HASHER_MIN_LENGTH=${ID_HASHER_MIN_LENGTH:-16}
export ID_HASHER_SALT=${ID_HASHER_SALT:-test-salt}
export NOTIFIER_ENGINE=${NOTIFIER_ENGINE:-rollbar}
export ROLLBAR_TOKEN=${ROLLBAR_TOKEN:-test-token}

go test -v -covermode=count ${PKG_LIST} -coverprofile ${TEST_RESULT_DIR}/.testCoverage.txt | tee ${TEST_RESULT_DIR}/test.log
echo ${PIPESTATUS[0]} >${TEST_RESULT_DIR}/test.out

cat ${TEST_RESULT_DIR}/test.log | go-junit-report >${TEST_RESULT_DIR}/report.xml

echo "----------------"
echo "${GRN}Result:${RESET}"
go tool cover -func ${TEST_RESULT_DIR}/.testCoverage.txt
exit $(cat ${TEST_RESULT_DIR}/test.out)
