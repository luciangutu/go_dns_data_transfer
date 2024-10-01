package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var ipAddr string
var transferIDs []string

func removeFromList(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// Decode numeric string representation back to UTF-8 text
func utf8FromNumericString(encoded string) (string, error) {
	var decoded strings.Builder
	for i := 0; i < len(encoded); i += 3 {
		numStr := encoded[i : i+3]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return "", err
		}
		decoded.WriteRune(rune(num))
	}
	return decoded.String(), nil
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	clientIP := w.RemoteAddr().String()
	log.Printf("Received DNS query from %s for domain %s", clientIP, r.Question[0].Name)

	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		switch q.Qtype {
		case dns.TypeA:
			if len(q.Name) > 253 {
				log.Println("Error: domain is too long!")
				return
			}

			domainParts := strings.Split(strings.TrimSpace(q.Name), ".")

			if len(domainParts) == 4 && domainParts[2] == "transfer" {
				transferData := domainParts[0]
				transferID := domainParts[1]

				if transferData == "start" {
					idx := slices.Index(transferIDs, transferID)
					if idx < 0 {
						transferIDs = append(transferIDs, transferID)
						log.Printf("Register new transfer with ID %s", transferID)
					}
				} else if transferData == "end" {
					idx := slices.Index(transferIDs, transferID)
					if idx >= 0 {
						transferIDs = removeFromList(transferIDs, transferID)
						log.Printf("Ended transfer with ID %s", transferID)
					}
				} else {
					idx := slices.Index(transferIDs, transferID)
					if idx >= 0 {
						log.Printf("Received %s for transfer ID %s", transferData, transferID)

						decoded, err := utf8FromNumericString(transferData)
						if err != nil {
							fmt.Println("Error:", err)
						}

						file, fileErr := os.OpenFile("data/"+transferID+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if fileErr != nil {
							fmt.Println("Error opening file:", fileErr)
							return
						}
						defer file.Close()

						if _, writeErr := file.WriteString(decoded); writeErr != nil {
							fmt.Println("Error writing to file:", writeErr)
							return
						}
					}
				}
			}

			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ipAddr))
			if err != nil {
				log.Printf("Failed to generate an A record for the DNS response (%s): %v)", q.Name, err)
				continue
			}
			msg.Answer = append(msg.Answer, rr)
		}
	}

	if err := w.WriteMsg(&msg); err != nil {
		log.Printf("Failed to send DNS response to %s: %v", clientIP, err)
	}
}

func main() {
	flag.StringVar(&ipAddr, "ip", "", "IP address to return for A records")
	flag.Parse()

	if ipAddr == "" {
		log.Fatal("Please provide an IP address using the -ip flag.")
	}

	if net.ParseIP(ipAddr) == nil {
		log.Fatalf("Invalid IP address: %s", ipAddr)
	}

	listen_addr := "0.0.0.0:53"

	dns.HandleFunc(".", handleDNSRequest)
	server := &dns.Server{Addr: listen_addr, Net: "udp"}
	log.Printf("Wildcard DNS server on %s will respond with %s", listen_addr, ipAddr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start DNS server: %s", err.Error())
	}
}
