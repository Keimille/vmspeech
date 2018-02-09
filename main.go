package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"

	// Imports the Google Cloud Speech API client package.
	speech "cloud.google.com/go/speech/apiv1"
	"github.com/spf13/viper"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	ctx := context.Background()

	//Set Flags
	//toEmailptr := flag.String("toEmail", "to pull from flag/asterisk", "define where email transcription should send")
	//toEmail := (*toEmailptr)
	filenameptr := flag.String("filename", "voicemail.wav", "load voicemail file")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		fmt.Println("Can't find config file for email auth")
		fmt.Println(viperErr)
	}

	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(*filenameptr)
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

	// Prints the results
	var transcript = ""
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			ioutil.WriteFile("output.txt", []byte(alt.Transcript), 0644)
			transcript = alt.Transcript
			fmt.Println(transcript)
		}
	}
	send(transcript)
}

func send(transcript string) {
	from := viper.GetString("emailSource")
	pass := viper.GetString("emailSourcePass")
	to := "......"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: New Voicemail\n\n" +
		transcript

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("voicemail transcription sent!")
}
