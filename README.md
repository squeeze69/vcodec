# vcodec - Squeeze69

## quick and dirty program to test for particular vcodec in an AVI (or any format supported by RIFF library)

## program written in [GO](https://golang.org)

Scan a RIFF file for the codec (vids in the strh and the structure in the following strf section), then exit with an error level 1 if the video codec is listed on the command line. The "heavy lift" is made by the image/riff library.

If it's used only with the file name, it prints the codec and exits with errorlevel 0.

It's partially based on one of the examples in the image/riff section (the read chunk,etc...).

i.e.:

vcodec file.avi DIV3 DX50 ... (case insensitive)

or:
vcodec -c codeclist file.avi

where codeclist is a "text" file with a codec per line

It's quite useful when you want to re-encode with ffmpeg or others a bunch of files only if they use certain codecs

Reference (Even if I do need to study them a little more)

- [Microsoft AVI riff file refernce](https://docs.microsoft.com/it-it/windows/win32/directshow/avi-riff-file-reference)
- [Microsoft bitmapinfoheader info](https://docs.microsoft.com/it-it/windows/win32/api/wingdi/ns-wingdi-bitmapinfoheader)
- [Alberta's University AVI format description](https://sites.ualberta.ca/dept/chemeng/AIX-43/share/man/info/C/a_doc_lib/ultimdia/ultiprgd/AVI.htm)
