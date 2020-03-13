#!/bin/bash
dnsserialfile=/usr/local/src/dnsserial-DO_NOT_DELETE_BEFORE_ASKING_CHRISTIAN.txt
zonefile=/var/named/zonefile.db
if [ -f zonefile ] ; then
	echo $[ $(grep serial ${zonefile}  | tr -d "\t"" ""\n"  | cut -d';' -f 1) + 1 ] | tee ${dnsserialfile}
else
	if [ ! -f ${dnsserialfile} ] || [ ! -s ${dnsserialfile} ]; then
		echo $(date +%Y%m%d00) | tee ${dnsserialfile}
	else
		echo $[ $(< ${dnsserialfile}) + 1 ] | tee ${dnsserialfile}
	fi
fi
##
##-30-
