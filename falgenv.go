// Package flagenv provides the ability to populate flags from
// environment variables.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Prefix Specify a prefix for environment variables.
var Prefix = ""

func contains(list []*flag.Flag, f *flag.Flag) bool {
	for _, i := range list {
		if i == f {
			return true
		}
	}
	return false
}

// ParseSet parses the given flagset. The specified prefix will be applied to
// the environment variable names.
func parseSet(prefix string, set *flag.FlagSet) error {
	var explicit []*flag.Flag
	var all []*flag.Flag
	set.Visit(func(f *flag.Flag) {
		explicit = append(explicit, f)
	})

	var err error
	set.VisitAll(func(f *flag.Flag) {
		if err != nil {
			return
		}
		all = append(all, f)
		if !contains(explicit, f) {
			name := strings.Replace(f.Name, ".", "_", -1)
			name = strings.Replace(name, "-", "_", -1)
			if prefix != "" {
				name = prefix + name
			}
			name = strings.ToUpper(name)
			val := os.Getenv(name)
			if val != "" {
				if ferr := f.Value.Set(val); ferr != nil {
					err = fmt.Errorf("failed to set flag %q with value %q", f.Name, val)
				}
			}
		}
	})
	return err
}

// Parse will set each defined flag from its corresponding environment
// variable . If dots or dash are presents in the flag name, they will be
// converted to underscores.
//
// If Parse fails, a fatal error is issued.
func parse() {
	if err := parseSet(Prefix, flag.CommandLine); err != nil {
		log.Fatalln(err)
	}
	flag.Parse()
}
