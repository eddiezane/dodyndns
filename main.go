package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"strings"

	"github.com/digitalocean/godo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

const ipInfoURL = "https://ipinfo.io/ip"

var domain = flag.String("domain", "", "The name of the domain")
var record = flag.String("record", "", "The name of the record")
var token = flag.String("token", "", "The DigitalOcean API token")

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	flag.Parse()
}

func main() {
	if *domain == "" {
		log.Fatal("--domain is required")
	}
	if *record == "" {
		log.Fatal("--record is required")
	}
	if *token == "" {
		log.Fatal("--record is required")
	}

	log.WithFields(log.Fields{
		"domain": *domain,
		"record": *record,
	}).Info("dodyndns starting run...")

	ip, err := fetchIP()
	if err != nil {
		log.Fatal(err)
	}
	if net.ParseIP(ip) == nil {
		log.Fatal("invalid ip: ", ip)
	}

	client := newGodoClient(*token)

	foundRecord, err := findRecordByName(client, *domain, *record)
	if err != nil {
		log.Fatal(err)
	}

	err = updateRecord(client, ip, *domain, foundRecord)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"ip":     ip,
		"domain": *domain,
		"record": *record,
	}).Info("record updated successfully")
}

func fetchIP() (string, error) {
	req, err := http.NewRequest(http.MethodGet, ipInfoURL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	ip := strings.TrimSpace(buf.String())
	return ip, nil
}

func updateRecord(client *godo.Client, ip, domain string, record *godo.DomainRecord) error {
	req := &godo.DomainRecordEditRequest{
		Data: ip,
	}

	_, _, err := client.Domains.EditRecord(context.Background(), domain, record.ID, req)
	if err != nil {
		return err
	}

	return nil
}

func findRecordByName(client *godo.Client, domain, recordName string) (*godo.DomainRecord, error) {
	records, _, err := client.Domains.Records(context.Background(), domain, nil)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.Name == recordName {
			return &record, nil
		}
	}

	return nil, errors.New("domain record not found")
}

func newGodoClient(token string) *godo.Client {
	ts := &tokenSource{
		AccessToken: token,
	}

	oauthClient := oauth2.NewClient(context.Background(), ts)
	client := godo.NewClient(oauthClient)

	return client
}
