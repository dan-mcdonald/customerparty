package main

import (
	"errors"
	"sort"
	"strings"
	"testing"
)

func TestSortByUserID(t *testing.T) {
	cust0 := &Customer{UserID: 3}
	cust1 := &Customer{UserID: 7}
	cust2 := &Customer{UserID: 12}
	cust3 := &Customer{UserID: 31}
	customerList := []*Customer{
		cust3,
		cust1,
		cust0,
		cust2,
	}
	sort.Sort(SortByUserID(customerList))
	if customerList[0] != cust0 {
		t.Error("expected customer 0 first")
	}
	if customerList[1] != cust1 {
		t.Error("expected customer 1 second")
	}
	if customerList[2] != cust2 {
		t.Error("expected customer 2 third")
	}
	if customerList[3] != cust3 {
		t.Error("expected customer 3 fourth")
	}
}

func TestNearOffice(t *testing.T) {
	if NearOffice(&Customer{}) {
		t.Error("Can't assume they're nearby if coordinates broken")
	}
	if !NearOffice(&Customer{Latitude: "53.3394075", Longitude: "-6.2584701"}) {
		t.Error("The pub is just a stone's throw from the office!")
	}
	if NearOffice(&Customer{Latitude: "51.8960528", Longitude: "-8.4980692"}) {
		t.Error("C'mon now Cork is a wee bit far")
	}
}

func everyoneButThomas(c *Customer) bool {
	return c.Name != "Thomas"
}

func TestCustomersReport(t *testing.T) {
	const customerObjects = `
	{"latitude": "52.986375", "user_id": 12, "name": "Christina McArdle", "longitude": "-6.043701"}
	{"latitude": "52.986375", "user_id": 9, "name": "Thomas", "longitude": "-6.043701"}
	{"latitude": ".986375", "user_id": 74, "name": "
unexpectednewline", "longitude": "-6.043701"}
	{}
	[]
	42
	true
	{"latitude": "51.92893", "user_id": 1, "name": "Alice Cahill", "longitude": "-10.27699"}
foo`
	customers, _ := CustomersReport(strings.NewReader(customerObjects), everyoneButThomas)
	if customers[0].Name != "Alice Cahill" {
		t.Errorf("Expected Alice first but got %#v", *customers[0])
	}
	if customers[1].Name != "Christina McArdle" {
		t.Errorf("Expected Christina second but got %#v", *customers[1])
	}
	if len(customers) != 2 {
		t.Errorf("Predicate didn't exclude customer %#v", *customers[2])
	}
}

type ReaderWith func(p []byte) (n int, err error)

func (r ReaderWith) Read(p []byte) (n int, err error) {
	return r(p)
}

func TestCustomersReportError(t *testing.T) {
	errorReader := ReaderWith(func(p []byte) (n int, err error) {
		return 0, errors.New("TestError")
	})
	_, err := CustomersReport(errorReader, everyoneButThomas)
	if err == nil {
		t.Error("should have returned an error for reader that always errors")
	}
}
