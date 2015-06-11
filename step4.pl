#!/usr/bin/perl

use strict;

use Image::Magick;
use Data::Dumper;
use WWW::Favicon;
use JSON;
use Digest::MD5;
use Getopt::Long qw(GetOptions);
 
my $temp_folder;
GetOptions('folder=s' => \$temp_folder) or die "Must provide folder name for storing icons";

while(<>) {
	chomp;
	s/,$//;
	my $icon = from_json($_);

	print to_json($icon) . ",\n";

	next if (!exists($$icon{"changed"}) || !$$icon{"changed"} || !exists($$icon{"icon_file"}));

	my $image = Image::Magick->new;
	my $x = $image->Read($$icon{"icon_file"});
	warn $x if $x;

	foreach my $i (@$image)
	{
		if ($i->Get('width') == 16 && $i->Get('height') == 16)
		{
			$x = $i->Write("$temp_folder/" . $$icon{"domain"} . ".png");
			warn $x if $x;
			last;
		}
	}
}
