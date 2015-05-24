#!/usr/bin/perl

use strict;

use Image::Magick;

my $domain = shift;
my $size = shift || 160;
my $subfolder = shift || 'large';

print "Resizing pngs/$domain.png to $subfolder/$domain.png at $size x $size\n";

my $image = Image::Magick->new;
my $x = $image->Read("pngs/$domain.png");
$image->Resize(width => $size, height => $size, filter => 'Point');

$x = $image->Write("$subfolder/$domain.png");
