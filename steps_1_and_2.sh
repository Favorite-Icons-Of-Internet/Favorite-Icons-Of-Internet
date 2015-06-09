#!/bin/bash

# step 1 - import latest alexa rankings
curl -s http://s3.amazonaws.com/alexa-static/top-1m.csv.zip | funzip > top-1m.csv
php step1.php < top-1m.csv

# step 2 - get all currently ranked domains and send jobs to the queue
rm -rf /tmp/.step2/
mkdir -p /tmp/.step2/
php step2.php | split -a 5 -l 1000 - /tmp/.step2/job_

for i in $( ls /tmp/.step2/); do
	enqueue FaviconPipelineDomains </tmp/.step2/$i
done

rm -rf /tmp/.step2/
