package main

import (
	"fmt"
	freshsalesclient "gomodules.xyz/freshsales-client-go"
)

func main() {
	c := freshsalesclient.DefaultFromEnv()

	email := "tamal@appscode.com"

	result, err := c.LookupByEmail(email, freshsalesclient.EntityContact)
	if err != nil {
		panic(err)
	}
	if len(result.Contacts.Contacts) == 0 {
		fmt.Println("Contact not found!")
	} else {
		fmt.Println("Contact found!")
		fmt.Println(result.Contacts.Contacts[0].Email)
		return
	}

	c2, err := c.CreateContact(&freshsalesclient.Contact{
		FirstName: "Tamal",
		LastName:  "Saha",
		JobTitle:  "CEO",
		Email:     email,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(c2.Email)

	//contacts, err := c.ListAllContacts()
	//if err != nil {
	//	panic(err)
	//}
	//for _, cc := range contacts {
	//	fmt.Println(cc.Email)
	//}
}
