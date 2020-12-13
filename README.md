# gopastebin

![E2E](https://github.com/arminaaki/gopastebin/workflows/E2E/badge.svg)


Pastebin GO API

The work is currently in Progress, but PRs are welcome!

## Install
```sh
go get https://github.com/arminaaki/gopastebin
```

## Usage
```go
import "https://github.com/arminaaki/gopastebin"
```

## Examples

```go
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
```
