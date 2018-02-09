# vmspeech

### A package for integrating wav file transcription and emailing attachments
#### current use case is for asterisk PBX implementations

1. export GOOGLE_APPLICATION_CREDENTIALS=googleconfig.json

_TODOs_
1. Integrate asterisk/FreePBX mailmcd to call built app, w/ cli Flags
2. locate logic for asterisk's temp WAV storage and pass to cli call

_sample mailmcd trigger_

/opt/vmspeech/dist/main-linux --filename="/var/spool/asterisk/default/{getExtension?}/.tmp" \
--toEmail="{emailAddress of mailbox}" \
--callerID="${VM_CALLERID}" 
