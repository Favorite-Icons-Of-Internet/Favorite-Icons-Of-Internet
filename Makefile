html:
	head -n30000 Quantcast-Top-Million.txt | perl geticons.pl --nogenimages --nofetch --width=1024

images:
	head -n30000 Quantcast-Top-Million.txt | perl geticons.pl --width=1024

parallel:
	head -n30000 Quantcast-Top-Million.txt | xargs -P10 -n1 perl geticons.pl --width=1280 --nogenpage >run.log 2>&1

test:
	head -n3000 Quantcast-Top-Million.txt | perl geticons.pl --nogenimages --nofetch --width=1024
