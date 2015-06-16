// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "crypto/tls"
    "errors"
    "os"
    "fmt"
    "log"

    "gopkg.in/ldap.v1"
)

var (
    LdapServer string = "192.168.99.100"
    LdapPort   uint16 = 389
    BaseDN     string = "dc=example,dc=org"
    BindDN     string = "cn=admin,dc=example,dc=org"
    BindPW     string = "admin"
    Filter     string = "(cn=admin)"
)

func search(l *ldap.Conn, filter string, attributes []string) (*ldap.Entry, error) {
    search := ldap.NewSearchRequest(
        BaseDN,
        ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
        filter,
        attributes,
        nil)

    sr, err := l.Search(search)
    if err != nil {
        log.Fatalf("ERROR: %s\n", err.Error())
        return nil, err
    }

    log.Printf(">>> Search: %s -> num of entries = %d\n", search.Filter, len(sr.Entries))
    if len(sr.Entries) == 0 {
        return nil, ldap.NewError(ldap.ErrorDebugging, errors.New(fmt.Sprintf("no entries found for: %s", filter)))
    }
    return sr.Entries[0], nil
}

func main() {
    
    log.Printf("This utility will attempt to bind to the specified ldap server, search and print the record for the binding user.")
    log.Printf("example: searchtest ldap.example.com ou=People,o=example.com uid=username@example.com,ou=People,o=example.com password")
    
    if len(os.Args) >=5 {
        LdapServer = os.Args[1]
        BaseDN = os.Args[2]
        BindDN = os.Args[3]
        BindPW = os.Args[4]
        Filter = fmt.Sprintf("(%s)", BindDN)
    }


    log.Printf("LdapServer: ", fmt.Sprintf("%s", LdapServer))
    log.Printf("LdapPort  : ", fmt.Sprintf("%d", LdapPort))
    log.Printf("BaseDN    : ", fmt.Sprintf("%s", BaseDN))
    log.Printf("BindDN    : ", fmt.Sprintf("%s", BindDN))
    log.Printf("BindPW    : (hidden)")
    log.Printf("Filter    : ", fmt.Sprintf("%s", Filter))
        
    log.Printf(">>> Attempting to connect...")
    //l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", LdapServer, LdapPort), &tls.Config{InsecureSkipVerify: true})
    l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", LdapServer, LdapPort))
    if err != nil {
        log.Fatalf("ERROR: %s\n", err.Error())
    }
    defer l.Close()
    log.Printf(">>> Connected. Turning Debug on...")
    l.Debug = true

    // Then startTLS
    log.Printf(">>> Attempting StartTLS...")
    err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
    if err != nil {
        log.Fatalf("ERROR: %s\n", err)
    }
    log.Printf(">>> StartTLS successful...")

    log.Printf(">>> Attempting Bind...")
    l.Bind(BindDN, BindPW)
    log.Printf(">>> Bind complete...")

    log.Printf(">>> Searching for record ... %s\n", Filter)
    entry, err := search(l, Filter, []string{})
    if err != nil {
        log.Fatal("could not get entry")
    }
    entry.PrettyPrint(0)
    log.Printf(">>> Search complete.")
}
