#!/usr/bin/env bash
WORKING_DIR=$(dirname $1)
CONFIG_FILE=$(basename $1)
elixir copy_deps.exs ${WORKING_DIR}
COPY_RESULT=$?
cd ${WORKING_DIR}
if [ ${COPY_RESULT} -eq 0 ];then
    gcloud app deploy ${CONFIG_FILE} --project $2
    rm -rf tmp_deps
    mv mix.exs.backup mix.exs
fi
