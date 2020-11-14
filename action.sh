#!/bin/sh

METHOD="${1}"
if [ "${METHOD}" = "" ]
then
  METHOD="patch"
fi

git-vertag --fetch "${METHOD}"
echo "::set-output name=vertag::$(git-vertag get)"
