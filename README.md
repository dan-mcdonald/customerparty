# CustomerParty

CustomerParty generates a report of customers within 100km of Intercom's office
in Dublin.

# Setup
## Install

```bash
$ go install github.com/dan-mcdonald/customerparty
```

## Usage

If you specify a file customerparty will read that file, otherwise it will read
from stdin. 

```bash
$ $(go env GOPATH)/bin/customerparty [filename]
```

To run against sample data:

```bash
$ $(go env GOPATH)/bin/customerparty $(go env GOPATH)/src/github.com/dan-mcdonald/customerparty/data/customerList.txt
```

# Input

Customer input data is expected to be one JSON object per line,
each object representing a customer. E.g.:

```JSON
{"latitude": "52.986375", "user_id": 12, "name": "Christina McArdle", "longitude": "-6.043701"}
{"latitude": "52.986375", "user_id": 9, "name": "Thomas", "longitude": "-6.043701"}
```

Each customer object should have at least the following properties:
 * `latitude` and `longitude` specified as JSON Strings representing signed decimal degrees
 * `user_id` integer JSON Number identifying the customer
 * `name` JSON String of the customer's name

Invalid objects will be ignored

# Output

The user_id and name of matching customers are output as tab separated values.
The output is sorted by user_id in ascending order.

Sample:
```
user_id	name
34	"Ian Kehoe"
58	"Nora Dempsey"
```

Names are ouput as quoted strings to make output parseable as the name may
contain arbitrary text and a newline character would be indistinguishable from
then end of the name.