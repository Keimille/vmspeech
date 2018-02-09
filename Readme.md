# vmspeech

## current issue > properly invokes.. but extension/CID variables not passing properly in mailmcd

### A package for integrating wav file transcription and emailing attachments
#### current use case is for asterisk PBX implementations

## current issue > properly invokes.. but extension/CID variables not passing properly in mailmcd
1. Create Google Cloud Project, and active Speech API
2. Create service account, and download json config file, placing it in same directory as executable
3. run -> export GOOGLE_APPLICATION_CREDENTIALS=googleconfig.json
4. or export GOOGLE_APPLICATION_CREDENTIALS=/opt/vmspeech/googleconfig.json

_TODOs_
1. Integrate asterisk/FreePBX mailmcd to call built app, w/ cli Flags
2. locate logic for asterisk's temp WAV storage and pass to cli call
3. delete WAV file after transmission

_build on mac for linux_
env GOOS=linux GOARCH=amd64 go build

_sample mailmcd trigger_

/opt/vmspeech/./main-linux --callerID=${VM_CALLERID} --extension=${VM_MAILBOX}
