#!/bin/bash

RUNS=$1
shift
SCRIPT=$@

for i in `seq 1 $RUNS`;
do
	$SCRIPT	
done   
