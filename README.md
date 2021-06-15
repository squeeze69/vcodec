# vcodec - Squeeze69

## quick and dirty program to test for particular vcodec in an AVI (or any format supported by RIFF library)

## program written in [GO](https://golang.org)

Scan a RIFF file for the "vids" section, then exit with an error level 1 if the video codec is listed on the command line. The "heavy lift" is made by the image/riff library.

i.e.:

vcodec file.avi DIV3 DX50 ...

It's quite useful when you want to re-encode with ffmpeg or others a bunch of files only if they use certain codecs
