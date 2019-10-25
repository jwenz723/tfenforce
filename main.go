package main

import (
	"encoding/json"
	"flag"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	tfplan := flag.String("tfplan", "plan.json", "Path to tfplan json file to be analyzed")
	flag.Parse()

	// Read in the tfplan json
	l := log.New(os.Stderr, "", 0)
	b, err := ioutil.ReadFile(*tfplan)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}
	var plan tfjson.Plan
	json.Unmarshal(b, &plan)

	// Define the checks that will be enforced
	enforcements := map[string]func(plan tfjson.Plan) []string{
		"cannot declare '*' IAM permission": CheckIamPolicyStar,
		"cannot define aws_vpc resources":   CheckVpcDefinition,
	}

	// Evaluate all defined enforcements
	violations := map[string][]string{}
	for rule, f := range enforcements {
		violations[rule] = f(plan)
	}

	// Print out all violations
	for rule, addresses := range violations {
		for _, a := range addresses {
			l.Printf("%s (%s)\n", rule, a)
		}
	}
}

// CheckVpcDefinition will validate that no aws_vpc resources
// have been defined
func CheckVpcDefinition(plan tfjson.Plan) []string {
	var addresses []string
	for _, resourceChange := range plan.ResourceChanges {
		if resourceChange.Type == "aws_vpc" {
			addresses = append(addresses, resourceChange.Address)
		}
	}
	return addresses
}

// CheckIamPolicyStar will validate that no IAM policy is defined with
// `Effect`:`Allow` && `Action`:`*`
//
// 10/24/2019 (Jeff Wenzbauer) - Currently this func is configured to only be able to parse json
// 		policies because terraform does not currently support yaml policy definitions.
func CheckIamPolicyStar(plan tfjson.Plan) []string {
	var addresses []string
	for _, resourceChange := range plan.ResourceChanges {
		if resourceChange.Type == "aws_iam_policy" {
			addr := resourceChange.Address

			var after map[string]interface{}
			var ok bool
			if after, ok = (resourceChange.Change.After).(map[string]interface{}); !ok {
				continue
			}

			var p string
			if p, ok = after["policy"].(string); !ok {
				continue
			}

			var j Policy
			// Unmarshal the policy that is nested as an escaped json string
			json.Unmarshal([]byte(p), &j)

			for _, s := range j.Statements {
				if strings.ToUpper(s.Effect) == "ALLOW" {
					for _, a := range s.Action {
						if a == "*" {
							addresses = append(addresses, addr)
						}
					}
				}
			}
		}
	}

	return addresses
}

type Policy struct {
	Statements []Statement `json:"Statement"`
}

type Statement struct {
	Effect string `json:"Effect"`
	Action Value  `json:"Action"`
}

// Value is a custom slice of strings needed for unmarshalling IAM policies
type Value []string

// UnmarshalJSON is a custom unmarshaller that will allow a json Value to be defined either as a string or []string
func (value *Value) UnmarshalJSON(b []byte) error {

	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	var p []string
	//  value can be string or []string, convert everything to []string
	switch v := raw.(type) {
	case string:
		p = []string{v}
	case []interface{}:
		var items []string
		for _, item := range v {
			items = append(items, fmt.Sprintf("%v", item))
		}
		p = items
	default:
		return fmt.Errorf("invalid %s value element: allowed is only string or []string", value)
	}

	*value = p
	return nil
}
