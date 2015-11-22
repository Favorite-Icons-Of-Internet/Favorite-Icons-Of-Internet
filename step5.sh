#!/bin/bash

DELAY=10m
ICONSTORAGE=/www/data/favoriteiconsofinternet.com/icons_archive/

# Waiting for the queue to empty up
if ! dequeue -n 0 FaviconPipelineDomains; then
	echo "Processing queue is still not empty, sleeping for $DELAY"
	sleep $DELAY
	exit
fi

UUID=`uuidgen`
RESULTSFOLDER="/tmp/.favicon_results_$UUID/"

# replace with "mv" command
aws s3 mv --recursive s3://favoriteiconsofinternet.com/results/ $RESULTSFOLDER 2>&1 >/dev/null

if [ ! -d $RESULTSFOLDER ]; then
	echo "No results downloaded, sleeping for $DELAY";
	sleep $DELAY
	exit
fi

TOTALFILES=`ls -1 $RESULTSFOLDER/ | wc -l`

if [ $TOTALFILES -eq 0 ]; then
	echo "No files to process after sync, sleeping for $DELAY"
	sleep $DELAY
	exit
fi

rm -rf /tmp/.favicon_result
mkdir -p /tmp/.favicon_result
for TARBALL in $RESULTSFOLDER/*; do
	tar -C /tmp/.favicon_result -xzf $TARBALL

	if [ -f /tmp/.favicon_result/favicon.manifest ]; then
		cat /tmp/.favicon_result/favicon.manifest | php step5ingest.php /tmp/.favicon_result $ICONSTORAGE
	fi

	rm /tmp/.favicon_result/*
	rm $TARBALL
done

rm -rf /tmp/.favicon_result

# Clean up temporary folder
rm -rf $RESULTSFOLDER

# Generate tile metadata and HTML file jobs for all domains that have icons and html file with links to all tiles
rm -rf /tmp/.step5_tiles
mkdir -p /tmp/.step5_tiles
php step5tiles.php /tmp/.step5_tiles /tmp/.step5.html

for TILE in $( ls /tmp/.step5_tiles/); do
	enqueue FaviconPipelineTiles </tmp/.step5_tiles/$TILE
done
