#!/bin/bash
for OUTPUT in $(ls *.zip)
do
	echo $(echo $OUTPUT | sed -e "s/-/_/g" -e 's/\.zip//g'): `sha256sum $OUTPUT | awk '{ print $1 }'`
done