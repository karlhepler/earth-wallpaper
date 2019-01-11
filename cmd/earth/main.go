package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const URL_ROOT = "https://cdn.star.nesdis.noaa.gov/GOES16/ABI/FD/GEOCOLOR"

func main() {
	// Define the allowed sizes
	allowedSizes := []string{"339", "678", "1808", "5424", "10848"}
	sizeChoices := strings.Join(allowedSizes, ", ")

	// Create the size flag
	sizeUsage := "The image is a square.\nThis size specifies the width and height of the image.\nChoices include: %s"
	inputSize := flag.String("size", "1808", fmt.Sprintf(sizeUsage, sizeChoices))

	// Create the filepath flag
	inputFilepath := flag.String("filepath", "/tmp/earth", "Where do you want to store the images?")

	// Parse all input flags
	flag.Parse()

	// Select the URL filename
	var urlFilename string
	switch *inputSize {
	case "339":
		urlFilename = "thumbnail"
	case "678", "1808", "5424":
		urlFilename = fmt.Sprintf("%sx%s", *inputSize, *inputSize)
	case "10848":
		urlFilename = "latest"
	default:
		log.Fatalf("%s is an invalid size. Choices include: %s", *inputSize, sizeChoices)
	}

	// Build the URL
	url := fmt.Sprintf("%s/%s.jpg", strings.TrimRight(URL_ROOT, "/"), urlFilename)

	// Download the image
	log.Printf("Downloading %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Create the output file
	filepath := strings.TrimRight(*inputFilepath, "/")
	filename := time.Now().Format("200601021504") + ".jpg"
	log.Printf("Creating %s\n", fmt.Sprintf("%s/%s", filepath, filename))
	err = os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(fmt.Sprintf("%s/%s", filepath, filename))
	if err != nil {
		log.Fatal(err)
	}

	// Copy the response body to the file
	log.Printf("Copying downloaded image to file")
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Change the desktop background
	log.Println("Changing desktop background")
	cmd := fmt.Sprintf("tell application \"Finder\" to set desktop picture to POSIX file \"%s/%s\"", filepath, filename)
	err = exec.Command("osascript", "-e", cmd).Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Success!")
}
