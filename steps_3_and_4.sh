#!/bin/bash
#
# Each time this script is executed, it will grab a list of domains from the queue and crawl them
#

HOME=/home/ec2-user/user-repo

UUID=`uuidgen`
TEMP_NAME3=favicon_step3_$UUID
TEMP_FOLDER3=/tmp/$TEMP_NAME3

mkdir $TEMP_FOLDER3

TEMP_NAME4=favicon_step4_$UUID
TEMP_FOLDER4=/tmp/$TEMP_NAME4

mkdir $TEMP_FOLDER4

$HOME/crawlerd/bin/crawlerd --output $TEMP_FOLDER3 2>/dev/null >$TEMP_FOLDER3/favicon.manifest

perl $HOME/step4.pl --folder $TEMP_FOLDER4 <$TEMP_FOLDER3/favicon.manifest >$TEMP_FOLDER4/favicon.manifest 2>/dev/null
rm -rf $TEMP_FOLDER3

$HOME/smu.sh $TEMP_FOLDER4/* >/dev/null 2>&1

tar -C /tmp/$TEMP_NAME4 -c . |gzip >/tmp/$TEMP_NAME4.tar.gz
rm -rf $TEMP_FOLDER4

aws s3 cp /tmp/$TEMP_NAME4.tar.gz s3://favoriteiconsofinternet.com/results/$TEMP_NAME4.tar.gz

rm /tmp/$TEMP_NAME4.tar.gz
