package wire

import (
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/mainchain"
	"github.com/globaldce/globaldce-gateway/utility"
	"math/rand"
	"time"
)

func (sw *Swarm) BroadcastMainblock(mb *mainchain.Mainblock) {
	//applog.Trace("our Mainchainlength is %d", currentmainchainlength)
	applog.Notice("Broadcasting Mainblock %d", mb.Height)

	//blockmsg:=NewMessage(MsgIdentifierBroadcastMainblock)
	blockmsg := EncodeBroadcastMainblock(uint32(0), mb)
	//blockmsg.WriteBytes(p.Connection)
	//
	for address, p := range sw.Peers {
		p.WriteMessage(blockmsg)
		applog.Trace("Writing to peer %v", address)
	}
}
func (sw *Swarm) BroadcastTransaction(tx *utility.Transaction) {
	//applog.Trace("our Mainchainlength is %d", currentmainchainlength)
	applog.Notice("Broadcasting Transaction !")
	//blockmsg:=NewMessage(MsgIdentifierBroadcastMainblock)
	txmsg := EncodeBroadcastTransaction(uint32(0), tx)
	//blockmsg.WriteBytes(p.Connection)
	//
	for address, p := range sw.Peers {

		p.WriteMessage(txmsg)
		applog.Trace("Writing to peer %v", address)
	}
}

func (sw *Swarm) RequestData(hash utility.Hash) {
	//applog.Trace("our Mainchainlength is %d", currentmainchainlength)
	applog.Notice("Broadcasting Request Data !")
	//blockmsg:=NewMessage(MsgIdentifierBroadcastMainblock)
	rdatamsg := EncodeRequestData(hash)
	//blockmsg.WriteBytes(p.Connection)
	//
	var keys []string
	for address, _ := range sw.Peers {
		keys = append(keys, address)
	}
	if len(keys) == 0 {
		return
	}
	rand.Seed(time.Now().UnixNano())
	randomkey := keys[rand.Intn(len(keys))]
	p := sw.Peers[randomkey]
	p.WriteMessage(rdatamsg)
	applog.Trace("Writing to peer %v", randomkey)

}

func (sw *Swarm) RequestDataFile(hash utility.Hash) {
	//applog.Trace("our Mainchainlength is %d", currentmainchainlength)
	applog.Notice("Broadcasting Request Data File !")
	//blockmsg:=NewMessage(MsgIdentifierBroadcastMainblock)
	rdatamsg := EncodeRequestDataFile(hash)
	//blockmsg.WriteBytes(p.Connection)
	//
	var keys []string
	for address, _ := range sw.Peers {
		keys = append(keys, address)
	}
	if len(keys) == 0 {
		return
	}
	rand.Seed(time.Now().UnixNano())
	randomkey := keys[rand.Intn(len(keys))]
	p := sw.Peers[randomkey]
	p.WriteMessage(rdatamsg)
	applog.Trace("Writing to peer %v", randomkey)

}

func (sw *Swarm) RelayMessage(msg *Message, originpeer *Peer) {

	applog.Trace("Relaying message ")

	for address, p := range sw.Peers {

		if address != originpeer.Address {
			p.WriteMessage(msg)
		}
	}
}

func (sw *Swarm) ReplyMessage(msg *Message, originpeer *Peer) {

	applog.Trace("Replying message ")

	//for address, p := range sw.Peers {

	//if address!=originpeer.Address{
	originpeer.WriteMessage(msg)
	//}
	//}
}
