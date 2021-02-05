package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/osm/dumpinen"
	"github.com/osm/flen"
)

func main() {
	// Setup the command line flags.
	addr := flag.String("addr", "https://dumpinen.com", "dumpinen address")
	cred := flag.String("cred", "", "credentials, format as user:pass")
	contentType := flag.String("content-type", "", "set a custom content type for the dump")
	deleteAfter := flag.String("delete-after", "", "delete after the given duration")
	id := flag.String("id", "", "id to retrieve from the dumpinen server")
	flen.SetEnvPrefix("DUMPINEN")
	flen.Parse()

	// Initialize the options slice.
	opts := []dumpinen.Option{dumpinen.WithAddr(*addr)}

	// We've got credentials, let's use em.
	if *cred != "" {
		opt, err := dumpinen.WithCredentials(*cred)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		opts = append(opts, opt)
	}

	// We've got a deleteAfter request, verify it and use it if everything
	// looks OK.
	if *deleteAfter != "" {
		opt, err := dumpinen.WithDeleteAfter(*deleteAfter)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		opts = append(opts, opt)
	}

	if *contentType != "" {
		opts = append(opts, dumpinen.WithContentType(*contentType))
	}

	// Initialize the client.
	d := dumpinen.NewClient(opts...)

	if *id != "" {
		// We've got an id set, this means that we don't want to upload
		// anything, instead we should try to fetch the dump from the server.
		d, err := d.Get(*id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", string(d))
	} else {
		// We want to upload something, read contents from stdin or a
		// file.
		var in io.Reader = os.Stdin
		if name := flag.Arg(0); name != "" && name != "-" {
			f, err := os.Open(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can't open file %s\n", name)
				return
			}
			defer f.Close()
			in = f
		}

		// Dump it.
		var dump string
		var err error
		if dump, err = d.Dump(in); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		fmt.Printf("%s", dump)
	}
}
