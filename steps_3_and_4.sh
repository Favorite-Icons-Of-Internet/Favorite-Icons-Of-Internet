#!/bin/bash
#
# Each time this script is executed, it will grab a list of domains from the queue and crawl them
#

UUID=`uuidgen`
TEMP_NAME3=favicon_step3_$UUID
TEMP_FOLDER3=/tmp/$TEMP_NAME3

mkdir $TEMP_FOLDER3

TEMP_NAME4=favicon_step4_$UUID
TEMP_FOLDER4=/tmp/$TEMP_NAME4

mkdir $TEMP_FOLDER4

dequeue FaviconPipelineDomains ./crawlerd/bin/crawlerd --output $TEMP_FOLDER3 2>/dev/null >$TEMP_FOLDER3/favicon.manifest

perl step4.pl --folder $TEMP_FOLDER4 <$TEMP_FOLDER3/favicon.manifest >$TEMP_FOLDER4/favicon.manifest
rm -rf $TEMP_FOLDER3

./smu.sh $TEMP_FOLDER4/* 2>&1 >/dev/null

tar -C /tmp/$TEMP_NAME4 -c . |gzip >/tmp/$TEMP_NAME4.tar.gz
rm -rf $TEMP_FOLDER4

aws s3 cp /tmp/$TEMP_NAME4.tar.gz s3://favoriteiconsofinternet.com/results/$TEMP_NAME4.tar.gz

rm /tmp/$TEMP_NAME4.tar.gz
