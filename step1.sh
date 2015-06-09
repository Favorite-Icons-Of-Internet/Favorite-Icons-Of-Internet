#!/bin/bash

curl -s http://s3.amazonaws.com/alexa-static/top-1m.csv.zip | funzip > top-1m.csv
php step1.php < top-1m.csv
