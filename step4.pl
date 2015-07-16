#!/usr/bin/perl

use strict;

use Image::Magick;
use Data::Dumper;
use WWW::Favicon;
use JSON;
use Digest::MD5;
use Getopt::Long qw(GetOptions);
use Digest::MD5 qw(md5_hex);
 
my $temp_folder;
GetOptions('folder=s' => \$temp_folder) or die "Must provide folder name for storing icons";

while(<>) {
	chomp;
	s/,$//;
	my $icon = from_json($_);

	my $skip = !exists($$icon{"changed"}) || !$$icon{"changed"} || !exists($$icon{"icon_file"});

	# No need for a temporary file name in results set
	delete $$icon{"icon_file"};

	print to_json($icon) . ",\n";

	next if $skip;

	my $image = Image::Magick->new;
	my $x = $image->Read($$icon{"icon_file"});
	warn $x if $x;

	foreach my $i (@$image)
	{
		if ($i->Get('width') == 16 && $i->Get('height') == 16)
		{
			my $filename = $$icon{"domain"} . '-' . $$icon{"hash"} . ".png"
			my $domain_hash = md5_hex($$icon{"domain"});
			my $folder = $temp_folder . '/' . substr($domain_hash, 0, 2) . '/' . substr($domain_hash, 2, 2);

			mkdir $temp_folder . '/' . substr($domain_hash, 0, 2);
			mkdir $temp_folder . '/' . substr($domain_hash, 0, 2) . '/' . substr($domain_hash, 2, 2);

			$x = $i->Write($temp_folder . '/' .
				substr($domain_hash, 0, 2) . '/' .
				substr($domain_hash, 2, 2) . '/' .
				$filename);

			warn $x if $x;
			last;
		}
	}
}
