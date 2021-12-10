package main

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	
	m := mail.NewV3Mail()

	from := "dev@bychoice.ch"
	name := "Quick Task Creator"
	e := mail.NewEmail(name, from)
	m.SetFrom(e)

	m.SetTemplateID("d-17b13c4e50ac46d3bdcd21ca352d53c6")

	p := mail.NewPersonalization()
	tos := []*mail.Email{
	mail.NewEmail("Trello Board", "niklausmaurer+xrh1wnkmzvtdnm56k3au@boards.trello.com"),
	}
	p.AddTos(tos...)

	p.SetDynamicTemplateData("task", "roll over and survive")
	p.SetDynamicTemplateData("description", "clean the floor first")

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}  
}