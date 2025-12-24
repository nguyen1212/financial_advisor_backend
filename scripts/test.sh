#!/bin/bash

RED=$(tput setaf 1)
GRN=$(tput setaf 2)
PUR=$(tput setaf 13)
RESET=$(tput sgr0)
TEST_RESULT_DIR="${TEST_RESULTS:-./test-results}"
mkdir -p ${TEST_RESULT_DIR}
echo "${GRN}Listing all packages${RESET}"
PKG_LIST+="$(go list ./... |
  grep -v /vendor/ |
  grep -v /migration |
  grep -v /mock | grep -v /cmd |
  grep -v /app/config |
  grep -v /app/registry |
  grep -v /app/delivery/rest/docs |
  grep -v /app/delivery/rest/handler |
  grep -v /app/external/db/gorm/mysql) "
echo "----------------"
echo "${GRN}Test:${RESET}"

go test -v -covermode=count ${PKG_LIST} -coverprofile ${TEST_RESULT_DIR}/.testCoverage.txt | tee ${TEST_RESULT_DIR}/test.log
echo ${PIPESTATUS[0]} >${TEST_RESULT_DIR}/test.out

cat ${TEST_RESULT_DIR}/test.log | go-junit-report >${TEST_RESULT_DIR}/report.xml

echo "----------------"
echo "${GRN}Result:${RESET}"
go tool cover -func ${TEST_RESULT_DIR}/.testCoverage.txt
exit $(cat ${TEST_RESULT_DIR}/test.out)
