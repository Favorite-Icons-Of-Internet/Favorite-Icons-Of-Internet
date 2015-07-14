#!/bin/bash

cd /home/ec2-user/user-repo/crawlerd
make

cd /home/ec2-user/queue-cli
git pull

sudo yum install -y pngcrush
