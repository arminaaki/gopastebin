package main

import (
	"log"
	"os"

	"github.com/arminaaki/gopastebin"
)

func main() {

	//Create pastbin API Client
	client, _ := gopastebin.NewClient(&gopastebin.AccountRequest{
		APIDevKey:       os.Getenv("PASTEBIN_API_DEV_KEY"),
		APIUserName:     os.Getenv("PASTEBIN_API_USER_NAME"),
		APIUserPassword: os.Getenv("PASTEBIN_API_USER_PASSWORD"),
	})

	// Create a past
	log.Print("creating a paste")
	_, _, err := client.Paste.Create(nil, &gopastebin.PasteRequest{
		APIPasteName:       "main.rb",
		APIPasteCode:       "puts \"Hello World!\"",
		APIPasteFormat:     "ruby",
		APIPastePrivate:    gopastebin.APIOptionPrivate,
		APIPasteExpireDate: "1M",
	})
	if err != nil {
		panic(err)
	}

	// List Pastes created by the User
	log.Print("list all available pastes")
	pastes, _, err := client.Paste.List(nil, &gopastebin.PasteListRequest{
		APIDevKey:  client.APIDevKey,
		APIUserKey: client.APIUserKey,
	})
	if err != nil {
		panic(err)
	}

	// Read a paste
	log.Printf("reading paste wih paste key %s", pastes[0].PasteKey)
	raw, _, err := client.Paste.GetRaw(nil, &gopastebin.PasteGetRawRequest{PasteKey: pastes[0].PasteKey})
	if err != nil {
		panic(err)
	}
	log.Printf("paste data returned %s", string(raw))

	// Delete Pastes
	log.Print("deleting all available pastes")
	for _, paste := range pastes {
		log.Printf("deleting paste %s: %s", paste.PasteKey, paste.PasteTitle)
		client.Paste.Delete(nil, &gopastebin.PasteDeleteRequest{APIPasteKey: paste.PasteKey})
	}
}
