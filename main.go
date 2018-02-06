package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	// Imports the Google Cloud Speech API client package.
	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name of the audio file to transcribe.
	filenamePtr := flag.String("filename", "voicemail.wav", "voicemails from asterisk")
	flag.Parse()
	fmt.Println(*filenamePtr)

	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(*filenamePtr)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 8000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}

	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	fromEmailptr := flag.String("fromEmail", "to pull from asterisk", "from address asterisk")
	fromEmail := (*fromEmailptr)
	toEmailptr := flag.String("toEmail", "to pull from flag/asterisk", "define where email transcription should send")
	toEmail := (*toEmailptr)
	exec.Command("/usr/sbin/sendmail", "-f", fromEmail, toEmail)

}
