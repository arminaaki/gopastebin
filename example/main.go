package main

import . "github.com/arminaaki/gopastebin"

import "os"

import "fmt"

func main() {

	//Create pastbin API Client
	client, _ := NewClient(&AccountRequest{
		APIDevKey:       os.Getenv("PASTEBIN_API_DEV_KEY"),
		APIUserName:     os.Getenv("PASTEBIN_API_USER_NAME"),
		APIUserPassword: os.Getenv("PASTEBIN_API_USER_PASSWORD"),
	})

	// Create a past
	_, _, err := client.Paste.Create(nil, &PasteRequest{
		APIPasteName:       "main.rb",
		APIPasteCode:       "puts \"Hello World!\"",
		APIPasteFormat:     "ruby",
		APIPastePrivate:    APIOptionPrivate,
		APIPasteExpireDate: "1M",
	})
	if err != nil {
		panic(err)
	}

	// List Pastes created by the User
	pastes, _, err := client.Paste.List(nil, &PasteListRequest{
		APIDevKey:  client.APIDevKey,
		APIUserKey: client.APIUserKey,
	})
	if err != nil {
		panic(err)
	}

	// Read a paste
	raw, _, err := client.Paste.GetRaw(nil, &PasteGetRawRequest{PasteKey: pastes[0].PasteKey})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))

	// Delete Pastes
	for _, paste := range pastes {
		fmt.Println("Name", paste.PasteTitle)
		fmt.Println("Deleting", paste.PasteKey, paste.PasteTitle)
		client.Paste.Delete(nil, &PasteDeleteRequest{APIPasteKey: paste.PasteKey})
	}

}
