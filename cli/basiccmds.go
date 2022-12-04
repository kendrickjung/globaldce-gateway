package cli

import (
    "github.com/globaldce/globaldce-gateway/applog"
    "os"
)


func createnewaddressCMD(){
    wlt:=Loadusermainwalletfile()
    address:=wlt.GenerateKeyPair()
    applog.Notice("Generated Address: %x",address)
    
    os.Exit(0)
}

func emptyCMD(){
    applog.Notice("%s command line interface: missing command",appName)
    applog.Notice("Try '%s help' for more information.",appName)
    os.Exit(0)
}
func helpCMD(){
    applog.Notice("\n%s command line interface",appName)
    applog.Notice("Copyright 2020-2022 globaldce developers\n")
    applog.Notice("Usage: %s COMMAND [OPTIONS]... ARGUMENTS...\n",appName)
    applog.Notice("BASIC COMMANDS:")
    applog.Notice("mine                                 Start mining")
    applog.Notice("managewallet                         Manage wallet")
    
    applog.Notice("createnewaddress                     Create a new address for receiving payments ")
    applog.Notice("                                     - by default the created address is a public key hash")
                                                        // TODO add support for more types of addresses and use of keyword '-addresstype:'
    //applog.Notice("createminingaddressesfile            Create mining addresses file")
    //applog.Notice("sendtoaddress                      send a transaction (create and broadcast a transaction)")
    //applog.Notice("sendnameregistration               send a name registration transaction")
    //applog.Notice("sendpublicpost                     send a public post transaction")
    //applog.Notice("getmainchaininfo                   .....")
    //applog.Notice("getwalletinfo                      .....")
    //applog.Notice("getblock                           ......")
    //applog.Notice("gettransaction                     ......")
    //applog.Notice("version or v                       version")
    applog.Notice("help or h               provide description of commands usage")
    applog.Notice("")
    applog.Notice("[OPTIONS]: for general use like with mine and managewallet commands")
    applog.Notice("-path=                               Sets appPath")
    applog.Notice("-port=                               Sets appLocalPort")
    applog.Notice("-miningaddressesfile=                Sets mining addresses file path")
    applog.Notice("")
    //

    os.Exit(0)
}