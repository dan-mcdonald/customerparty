package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// CustomerJSON corresponds to the JSON format for customer objects. We need
// user_id to be a pointer so we can detect when it was not present. Otherwise
// we couldn't tell it apart from the default value of 0. This isn't a problem
// for the other fields because invalid Lat/Lon will be ignored and Name isn't
type CustomerJSON struct {
	Latitude, Longitude, Name string
	UserID                    *int `json:"user_id"`
}

// Customer represents a valid customer, no pointers so easy to work with
type Customer struct {
	Latitude, Longitude, Name string
	UserID                    int
}

// asCustomer builds a *Customer object if the object is valid, i.e. it has a
// user_id set
func (cj *CustomerJSON) asCustomer() *Customer {
	if cj.UserID == nil {
		return nil
	}
	return &Customer{cj.Latitude, cj.Longitude, cj.Name, *cj.UserID}
}

// getCoord is a convenience function to extract the coordinate from a Customer
func (c *Customer) getCoord() (*Coord, error) {
	return CoordFromLatLon(c.Latitude, c.Longitude)
}

// NearOffice determines whether customer is close enough to invite to the party
func NearOffice(customer *Customer) bool {
	// This function hardcodes a lot of business logic that could be generalized
	// Office location and distance could be moved into a struct and this could
	// become a receiver function. I haven't done this yet because YAGNI
	// https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it
	officeLocation, _ := CoordFromLatLon("53.339428", "-6.257664")
	customerLocation, err := customer.getCoord()
	if err != nil {
		fmt.Printf("warning could not parse coordinates for user %#v, error: %#v\n", customer, err)
		return false
	}
	return DistanceKm(officeLocation, customerLocation) <= 100
}

// CustomerPredicate is a type for predicate function over the Customer type
type CustomerPredicate func(*Customer) bool

// SortByUserID implements sort.Interface to sort array of *Customer by user_id
type SortByUserID []*Customer

func (array SortByUserID) Len() int           { return len(array) }
func (array SortByUserID) Swap(i, j int)      { array[i], array[j] = array[j], array[i] }
func (array SortByUserID) Less(i, j int) bool { return array[i].UserID < array[j].UserID }

// CustomersReport returns array of customer objects read from each line
// of customerData that satisfy the predicate function sorted by user_id (asc)
func CustomersReport(customerData io.Reader, pred CustomerPredicate) ([]*Customer, error) {
	matchingCustomers := make([]*Customer, 0)
	lineScanner := bufio.NewScanner(customerData)
	for lineScanner.Scan() { // runs once for each input line
		var cj CustomerJSON

		line := lineScanner.Bytes()

		// attempt to decode as CustomerJSON
		if err := json.Unmarshal(line, &cj); err != nil {
			fmt.Printf("Couldn't parse customer JSON: %s\n", err)
			continue
		}
		// try to build a Customer from the CustomerJSON
		c := cj.asCustomer()
		// append if Customer was built successfully and predicate matches
		if c != nil && pred(c) {
			matchingCustomers = append(matchingCustomers, c)
		}
	}
	if err := lineScanner.Err(); err != nil {
		return nil, err
	}
	// Sort by user_id ascending
	sort.Sort(SortByUserID(matchingCustomers))
	return matchingCustomers, nil
}

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s [filename]\n", os.Args[0])
		fmt.Println("If no file specified reads from stdin")
		return
	}
	file := os.Stdin
	if len(os.Args) == 2 {
		var err error
		file, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Couldn't open customer list file: %#v\n", err)
			return
		}
	}

	customers, err := CustomersReport(file, NearOffice)
	if err != nil {
		fmt.Printf("Failure generating report: %#v\n", err)
		return
	}
	fmt.Println("user_id\tname")
	for _, customer := range customers {
		fmt.Printf("%d\t%#v\n", customer.UserID, customer.Name)
	}
}
