package main

import (
	"fmt"

	freshsalesclient "gomodules.xyz/freshsales-client-go"
)

func main() {
	c := freshsalesclient.DefaultFromEnv()

	result, err := c.LookupByEmail("tamal@appscode.com", freshsalesclient.EntityContact)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Contacts.Contacts[0].Email)

	//contacts, err := c.ListAllContacts()
	//if err != nil {
	//	panic(err)
	//}
	//for _, cc := range contacts {
	//	fmt.Println(cc.Email)
	//}
}
