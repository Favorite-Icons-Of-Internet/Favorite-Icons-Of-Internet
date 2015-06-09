#!/bin/bash

rm -rf /tmp/.step2/
mkdir -p /tmp/.step2/
php step2.php | split -a 5 -l 1000 - /tmp/.step2/job_

for i in $( ls /tmp/.step2/); do
	enqueue FaviconPipelineDomains </tmp/.step2/$i
done

rm -rf /tmp/.step2/
