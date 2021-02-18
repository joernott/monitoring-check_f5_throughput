#!/bin/bash
SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"
RPMBUILD_DIR=$(dirname ${SCRIPT_DIR})
SPECS_DIR=${RPMBUILD_DIR}/SPECS
RPMS_DIR=${RPMBUILD_DIR}/RPMS
SRPMS_DIR=${RPMBUILD_DIR}/SRPMS
BUILD_DIR=${RPMBUILD_DIR}/BUILD
SOURCES_DIR=${RPMBUILD_DIR}/SOURCES
LOGLEVEL=5
APP="check_f5_throughput"
VERSION=$1
RELEASE=$2

function init() {
    for DIR in ${SPECS_DIR} ${RPMS_DIR} ${SRPMS_DIR} ${BUILD_DIR} ${SOURCES_DIR}; do
        if [ ! -d ${DIR} ]; then
            echo "$DIR not found. Creating it."
            mkdir -p ${DIR}
        fi
    done
    if [ -z "${VERSION}" ]; then
        cat <<EOF
Usage: $0 VERSION [RELEASE]

where VERSION is the version of RPM, RELEASE is the release-number for that version
EOF
        exit 2
    fi
    if [ -z "${RELEASE}" ]; then
        RELEASE=1
    fi
}

init
cp "../${APP}.tar.gz" "${SOURCES_DIR}/"
cd ${RPMBUILD_DIR}
rpmbuild --define="_topdir ${RPMBUILD_DIR}" \
         --define "version ${VERSION}" \
         --define "release ${RELEASE}" \
         -ba ${SPECS_DIR}/${APP}.spec

