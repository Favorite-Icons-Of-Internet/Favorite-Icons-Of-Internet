WIDTH=1280
NUMICONS=30000
# Quantcast top 1M
DOMAINLIST=Quantcast-Top-Million.txt
# Alexa top 1M
#DOMAINLIST=top-1m.csv

html:
	head -n${NUMICONS} ${DOMAINLIST} | perl geticons.pl --nogenimages --nofetch --width=${WIDTH}

images:
	head -n${NUMICONS} ${DOMAINLIST} | perl geticons.pl --width=${WIDTH}

parallel:
	head -n${NUMICONS} ${DOMAINLIST} | xargs -P10 -n1 perl geticons.pl --width=${WIDTH} --nogenpage >run.log 2>&1

test:
	head -n${NUMICONS} ${DOMAINLIST} | perl geticons.pl --nogenimages --nofetch --width=${WIDTH}
