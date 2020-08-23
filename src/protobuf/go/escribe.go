package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
	pb "gz.com/grpc/addressbook"
)

func promptForAddress(r io.Reader) (*pb.Person, error) {
	//Asigna p como una referencia a Person
	p := &pb.Person{}
	//Crea un reader
	rd := bufio.NewReader(r)

	//Informa el Id
	fmt.Print("Enter person ID number: ")
	if _, err := fmt.Fscanf(rd, "%d\n", &p.Id); err != nil {
		return p, err
	}
	fmt.Print("Enter name: ")
	name, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	// We trim the whitespace because rd.ReadString includes the trailing
	// newline character in its output.
	p.Name = strings.TrimSpace(name)

	fmt.Print("Enter email address (blank for none): ")
	email, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Email = &email

	for {
		fmt.Print("Enter a phone number (or leave blank to finish): ")
		phone, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}
		// The PhoneNumber message type is nested within the Person
		// message in the .proto file.  This results in a Go struct
		// named using the name of the parent prefixed to the name of
		// the nested message.  Just as with pb.Person, it can be
		// created like any other struct.
		//Define pn como una referencia a Person_PhoneNumber, y especifica el valor Number
		pn := &pb.Person_PhoneNumber{
			Number: phone,
		}

		fmt.Print("Is this a mobile, home, or work phone? ")
		ptype, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		ptype = strings.TrimSpace(ptype)

		// A proto enum results in a Go constant for each enum value.
		var tmp pb.Person_PhoneType
		switch ptype {
		case "mobile":
			tmp = pb.Person_MOBILE
		case "home":
			tmp = pb.Person_HOME
		case "work":
			tmp = pb.Person_WORK
		default:
			fmt.Printf("Unknown phone type %q.  Using default.\n", ptype)
		}
		pn.Type = &tmp

		// A repeated proto field maps to a slice field in Go.  We can
		// append to it like any other slice.
		p.Phones = append(p.Phones, pn)
	}

	return p, nil
}

func listPeople(w io.Writer, book *pb.AddressBook) {
	for _, p := range book.People {
		writePerson(w, p)
	}
}

func writePerson(w io.Writer, p *pb.Person) {
	fmt.Fprintln(w, "Person ID:", p.Id)
	fmt.Fprintln(w, "  Name:", p.Name)
	if *p.Email != "" {
		fmt.Fprintln(w, "  E-mail address:", *p.Email)
	}

	for _, pn := range p.Phones {
		switch *pn.Type {
		case pb.Person_MOBILE:
			fmt.Fprint(w, "  Mobile phone #: ")
		case pb.Person_HOME:
			fmt.Fprint(w, "  Home phone #: ")
		case pb.Person_WORK:
			fmt.Fprint(w, "  Work phone #: ")
		}
		fmt.Fprintln(w, pn.Number)
	}
}

// Main reads the entire address book from a file, adds one person based on
// user input, then writes it back out to the same file.
func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}

	fname := os.Args[1]

	// Read the existing address book.
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found.  Creating new file.\n", fname)
		} else {
			log.Fatalln("Error reading file:", err)
		}
	}

	// [START marshal_proto]
	book := &pb.AddressBook{}
	// [START_EXCLUDE]
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	// Add an address.
	addr, err := promptForAddress(os.Stdin)
	if err != nil {
		log.Fatalln("Error with address:", err)
	}
	book.People = append(book.People, addr)
	// [END_EXCLUDE]

	// Write the new address book back to disk.
	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
	// [END marshal_proto]

	// [START unmarshal_proto]
	// Read the existing address book.
	entrada, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	libro := &pb.AddressBook{}
	if err := proto.Unmarshal(entrada, libro); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}
	// [END unmarshal_proto]

	listPeople(os.Stdout, libro)

}
