package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := make([]byte, 5)
	for i := range id {
		id[i] = charset[r.Intn(len(charset))]
	}
	return string(id)
}

func sendDNSQuery(domain, dnsServer string) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	client := new(dns.Client)
	response, _, err := client.Exchange(msg, dnsServer)
	if err != nil {
		log.Printf("Failed to send DNS query %s: %v", domain, err)
		return
	}

	if response.Rcode != dns.RcodeSuccess {
		log.Printf("DNS query failed for %s: %s", domain, dns.RcodeToString[response.Rcode])
		return
	}

	log.Printf("Sent %s", domain)
}

// Encode UTF-8 text to numeric string representation
func utf8ToNumericString(text string) string {
	var encoded strings.Builder
	for _, char := range text {
		encoded.WriteString(fmt.Sprintf("%03d", char))
	}
	return encoded.String()
}

// splitIntoChunks splits the input string into chunks of the specified size.
func splitIntoChunks(data string, chunkSize int) []string {
	var chunks []string

	for i := 0; i < len(data); i += chunkSize {
		// Ensure we don't exceed the length of the string
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}

	return chunks
}

func main() {
	// Set up command-line flags
	inputFilePath := flag.String("file", "", "Path to the input file")
	dnsServer := flag.String("dns", "127.0.0.1:53", "DNS server to query")
	flag.Parse()

	if *inputFilePath == "" {
		log.Fatal("Please provide an input file using the -file flag.")
	}

	randomID := generateRandomID()
	log.Printf("Generated random ID: %s", randomID)

	// Step 1: Send start query
	startDomain := fmt.Sprintf("start.%s.transfer", randomID)
	sendDNSQuery(startDomain, *dnsServer)

	// Step 2: Send data queries
	data, err := os.ReadFile(*inputFilePath)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}
	rawData := string(data)

	// Split into chunks that are less than 10 characters
	chunks := splitIntoChunks(rawData, 10)

	for _, chunk := range chunks {
		dataDomain := fmt.Sprintf("%s.%s.transfer", utf8ToNumericString(chunk), randomID)
		sendDNSQuery(dataDomain, *dnsServer)
	}

	// Step 3: Send end query
	endDomain := fmt.Sprintf("end.%s.transfer", randomID)
	sendDNSQuery(endDomain, *dnsServer)

	log.Println("All queries sent.")
}
