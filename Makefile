WIDTH=1280

html:
	head -n30000 Quantcast-Top-Million.txt | perl geticons.pl --nogenimages --nofetch --width=${WIDTH}

images:
	head -n30000 Quantcast-Top-Million.txt | perl geticons.pl --width=${WIDTH}

parallel:
	head -n30000 Quantcast-Top-Million.txt | xargs -P10 -n1 perl geticons.pl --width=${WIDTH} --nogenpage >run.log 2>&1

test:
	head -n3000 Quantcast-Top-Million.txt | perl geticons.pl --nogenimages --nofetch --width=${WIDTH}
