package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	// Imports the Google Cloud Speech API client package.
	speech "cloud.google.com/go/speech/apiv1"
	"github.com/scorredoira/email"
	"github.com/spf13/viper"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	ctx := context.Background()

	//Set Flags
	//toEmailptr := flag.String("toEmail", "to pull from flag/asterisk", "define where email transcription should send")
	filenameptr := flag.String("filename", "voicemail.wav", "load voicemail file location (temp file?)")
	callerIDptr := flag.String("callerID", "", "passed from asterisk ${VM_CALLERID}")
	exentsionptr := flag.String("extension", "", "Passed from asterisk, VM_Mailbox")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		fmt.Println("Can't find config file for email auth")
		fmt.Println(viperErr)
	}
	//call asteriskConfig
	toEmail := asteriskConfig(*exentsionptr)
	println("Line 41..... ", toEmail)

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
	var transcript string
	var confidence float32
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			transcript = alt.Transcript
			confidence = alt.Confidence
		}
	}
	send(*callerIDptr, transcript, confidence, toEmail, *filenameptr)
}

func send(callerIDptr string, transcript string, confidence float32, toEmail string, filenameptr string) {
	// compose the message
	m := email.NewMessage("New Voicemail From -> "+callerIDptr, transcript)
	m.From = mail.Address{Name: "TTS Voicemail", Address: viper.GetString("emailSource")}
	m.To = []string{toEmail}

	// add attachments
	if err := m.Attach(filenameptr); err != nil {
		log.Fatal(err)
	}

	// send it
	auth := smtp.PlainAuth("", viper.GetString("emailSource"),
		viper.GetString("emailSourcePass"), "smtp.gmail.com")
	if err := email.Send("smtp.gmail.com:587", auth, m); err != nil {
		log.Fatal(err)
	}
	log.Print("voicemail sent (attachment --> ", filenameptr, ")")
}

func asteriskConfig(extensionptr string) string {
	//pull voicemail.conf file from asterisk
	file, err := os.Open("/etc/asterisk/voicemail.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), extensionptr+" =>") {
			line = (scanner.Text())
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(line)
	//parse returned extension config for destination email address
	s := strings.Split(line, ",")
	destemail := s[2]
	fmt.Println("Line 124.... ", destemail)
	return destemail
}
