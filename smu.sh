#!/bin/bash

for FILE in "$@"
do
	# optipng (skip -o7 to run faster)
	optipng -o7 $FILE

	# pngcrush (skip -brute to run faster)
	pngcrush -rem alla -brute -reduce $FILE $FILE.temp
	mv $FILE.temp $FILE

	# pngout - closed source, non-windows binaries here
	# (add parameter -s2 to run faster)
	pngout $FILE

	# advpng (use -z2 to run faster)
	advpng -z4 $FILE

	# deflopt - windows only
	#$ deflopt my.png
done
