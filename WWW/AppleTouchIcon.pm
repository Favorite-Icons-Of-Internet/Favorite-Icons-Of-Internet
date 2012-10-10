package WWW::AppleTouchIcon;
use strict;
use warnings;
use base qw/Class::Accessor::Fast Exporter/;

use Carp;
use LWP::UserAgent;
use HTML::TreeBuilder;
use HTML::ResolveLink;

our $VERSION = '0.1';
our @EXPORT_OK = qw/detect_apple_touch_icon_url/;

__PACKAGE__->mk_accessors(qw/ua/);

sub new {
    my $self = shift->SUPER::new(@_);

    $self->{ua} = do {
        my $ua = LWP::UserAgent->new;
        $ua->timeout(10);
        $ua->max_size(1024*1024);
        $ua->env_proxy;
        $ua;
    };

    $self;
}

sub detect_apple_touch_icon_url($) {
    __PACKAGE__->detect(shift);
}

sub detect {
    my ($self, $url) = @_;
    $self = $self->new unless ref $self;

    my $res = $self->ua->get($url);
    croak 'request failed: ' . $res->status_line unless $res->is_success;

    my $resolver = HTML::ResolveLink->new( base => $res->base );
    my $html = $resolver->resolve( $res->content );

    my $tree = HTML::TreeBuilder->new;
    $tree->parse($html);
    $tree->eof;

    my ($apple_touch_icon_url) = grep {$_} map { $_->attr('href') } $tree->look_down(
        _tag => 'link',
        rel  => qr/^apple-touch-icon$/i,
    );
    my ($apple_touch_icon_precomposed_url) = grep {$_} map { $_->attr('href') } $tree->look_down(
        _tag => 'link',
        rel  => qr/^apple-touch-icon-precomposed$/i,
    );

    my $icon_url = $apple_touch_icon_precomposed_url || $apple_touch_icon_url;

    unless ($icon_url) {
        $icon_url = $res->base->clone;
        $icon_url->path('/apple-touch-icon.png');
    }

    $tree->delete;

    "$icon_url";
}

=head1 NAME

WWW::AppleTouchIcon - perl module to detect favicon url

=head1 SYNOPSIS

    use WWW::AppleTouchIcon qw/detect_apple_touch_icon_url/;
    my $apple_touch_icon_url = detect_apple_touch_icon_url('http://example.com/');
    
    # or OO way
    use WWW::AppleTouchIcon;
    my $apple_touch_icon = WWW::AppleTouchIcon->new;
    my $apple_touch_icon_url = $apple_touch_icon->detect('http://example.com/');

=head1 DESCRIPTION

This module provide simple interface to detect apple touch icon url of specified url.

=head1 METHODS

=head2 new

Create new WWW::AppleTouchIcon object.

=head2 detect($url)

Detect apple touch icon url of $url.

=head1 EXPORT FUNCTIONS

=head2 detect_apple_touch_icon_url($url)

Same as $self->detect described above.

=head1 AUTHOR

Sergey Chernyshev <cpan.org@antispam.sergeychernyshev.com>
based on WWW::Favicon by Daisuke Murase <typester@cpan.org>

=head1 COPYRIGHT

This program is free software; you can redistribute
it and/or modify it under the same terms as Perl itself.

The full text of the license can be found in the
LICENSE file included with this module.

=cut

1;
