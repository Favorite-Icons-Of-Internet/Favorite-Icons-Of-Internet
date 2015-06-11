#!/bin/bash

UUID=`uuidgen`
TEMP_NAME=favicon_step4_$UUID
TEMP_FOLDER=/tmp/$TEMP_NAME

mkdir $TEMP_FOLDER

perl step4.pl --folder $TEMP_FOLDER >$TEMP_FOLDER/favicon.manifest

./smu.sh $TEMP_FOLDER/* >/dev/null 2>&1

tar -C /tmp/$TEMP_NAME -c . |gzip >/tmp/$TEMP_NAME.tar.gz
rm -rf $TEMP_FOLDER

aws s3 cp /tmp/$TEMP_NAME.tar.gz s3://favoriteiconsofinternet.com/results/$TEMP_NAME.tar.gz

rm /tmp/$TEMP_NAME.tar.gz
