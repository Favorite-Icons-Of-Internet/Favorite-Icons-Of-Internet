import os
import sys
import argparse
from PIL import Image


def main():
    parser = argparse.ArgumentParser(
      description='Convert given image to png format.')

    parser.add_argument('imagepath', metavar='image', type=str,
                        help='path to image')
    parser.add_argument('--width', default=16, type=int,
                        help='resize to specified width')
    parser.add_argument('--height', default=16, type=int,
                        help='resize to specified height')
    parser.add_argument('--prefix', default='p', type=str,
                        help='prefix for processed image')
    args = parser.parse_args()

    if not os.path.isfile(args.imagepath):
        sys.exit('Error: file "%s" does not exist' % args.imagepath)

    img_size = (args.width, args.height)

    im = Image.open(args.imagepath)
    if im.size != img_size:
        im = im.resize(img_size)

    # return file's basename without file extension
    fname, _ = os.path.splitext(os.path.basename(args.imagepath))
    dname = os.path.dirname(args.imagepath)

    # new name for processed file
    nname = args.prefix + '-' + fname + '.png'

    output = os.path.join(dname, nname)

    im.save(output, 'png')
    im.close()

    print(output)


if __name__ == '__main__':
    main()
