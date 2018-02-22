#!/bin/bash
# Asterisk voicemail attachment conversion script, including voice recognition
# Use Voice Recognition Engine provided by Google API
#
#based on original by -> http://bernaerts.dyndns.org/linux/179-asterisk-voicemail-mp3

# set language for voice recognition (en-US, en-GB, fr-FR, ...)
LANGUAGE="en-US"

# set PATH
PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

# save the current directory
pushd .

# create a temporary directory and cd to it
TMPDIR=$(mktemp -d)
cd $TMPDIR

# dump the stream to a temporary file
cat >> stream.org

# get the boundary
BOUNDARY=$(grep "boundary=" stream.org | cut -d'"' -f 2)

# cut the file into parts
# stream.part - header before the boundary
# stream.part1 - header after the bounday
# stream.part2 - body of the message
# stream.part3 - attachment in base64 (WAV file)
# stream.part4 - footer of the message
awk '/'$BOUNDARY'/{i++}{print > "stream.part"i}' stream.org

# if mail is having no audio attachment (plain text)
PLAINTEXT=$(cat stream.part1 | grep 'plain')
if [ "$PLAINTEXT" != "" ]
then

  # prepare to send the original stream
  cat stream.org > stream.new

# else, if mail is having audio attachment
else

  # cut the attachment into parts
  # stream.part3.head - header of attachment
  # stream.part3.wav.base64 - wav file of attachment (encoded base64)
  sed '7,$d' stream.part3 > stream.part3.wav.head
  sed '1,6d' stream.part3 > stream.part3.wav.base64

  # convert the base64 file to a wav file
  dos2unix -o stream.part3.wav.base64
  base64 -di stream.part3.wav.base64 > stream.part3.wav

  # change Type: audio/x-wav or audio/x-WAV to Type: audio/mpeg
  # change name="msg----.wav" or name="msg----.WAV"
  sed 's/x-[wW][aA][vV]/mpeg/g' stream.part3.wav.head

# Send WAV to Watson Speech to Text API. Must use "Narrowband" (aka 8k) model since WAV is 8k sample.
curl -s -X POST \
    --limit-rate 40000 \
    --header "Content-Type: audio/wav" \
    --data-binary @stream.part3.wav \
    "https://speech.googleapis.com/v1beta1/speech:syncrecognize?key=${API_KEY}" 1>audio.txt

# Extract transcript results from JSON response
TRANSCRIPT=`cat audio.txt | grep transcript | sed 's#^.*"transcript": "##g' | sed 's# "$##g'`


  # generate first part of mail body, converting it to LF only
  mv stream.part stream.new
  cat stream.part1 >> stream.new
  sed '$d' < stream.part2 >> stream.new

  # beginning of transcription section
  echo "---" >> stream.new

  # end of message body
  tail -1 stream.part2 >> stream.new

  # append mp3 header
  cat stream.part3.wav.head >> stream.new
  dos2unix -o stream.new

  # append base64 mp3 to mail body, keeping CRLF
  unix2dos -o stream.part3.wav.base64
  cat stream.part3.wav.base64 >> stream.new

  # append end of mail body, converting it to LF only
  echo "" >> stream.tmp
  echo "" >> stream.tmp
  cat stream.part4 >> stream.tmp
  dos2unix -o stream.tmp
  cat stream.tmp >> stream.new

# send the mail thru sendmail
cat stream.new | sendmail -t

# go back to original directory
popd

# remove all temporary files and temporary directory
rm -Rf $TMPDIR
