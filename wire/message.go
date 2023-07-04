package wire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/utility"
	"net"
	//"os"
)

const MsgIdentifierLength = 20

const (
	MsgIdentifierRequestHandshake = "REQUEST_HANDSHAKE"
	MsgIdentifierReplyHandshake   = "REPLY_HANDSHAKE"

	MsgIdentifierRequestMainchainLength = "REQUEST_MAINCHAINLENGTH"
	MsgIdentifierReplyMainchainLength   = "REPLY_MAINCHAINLENGTH"

	MsgIdentifierRequestMainheaders = "REQUEST_MAINHEADERS"
	MsgIdentifierReplyMainheaders   = "REPLY_MAINHEADERS"

	MsgIdentifierRequestMainblockTransactions = "REQUEST_MAINBLOCKTRANSACTIONS"
	MsgIdentifierReplyMainblockTransactions   = "REPLY_MAINBLOCKTRANSACTIONS"

	MsgIdentifierBroadcastMainblock   = "BROADCAST_MAINBLOCK"
	MsgIdentifierBroadcastTransaction = "BROADCAST_TRANSACTION"

	MsgIdentifierRequestData = "REQUEST_DATA"
	MsgIdentifierReplyData   = "REPLY_DATA"

	MsgIdentifierRequestDataFile = "REQUEST_DATAFILE"
	MsgIdentifierReplyDataFile   = "REPLY_DATAFILE"

	MainNetworkIdentifier = "9184"
)

func DecodeIdentifier(msgidentifier []byte) string {
	switch {
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestHandshake)):
		return MsgIdentifierRequestHandshake
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyHandshake)):
		return MsgIdentifierReplyHandshake

	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestMainchainLength)):
		return MsgIdentifierRequestMainchainLength
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyMainchainLength)):
		return MsgIdentifierReplyMainchainLength
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestMainheaders)):
		return MsgIdentifierRequestMainheaders
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyMainheaders)):
		return MsgIdentifierReplyMainheaders
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestMainblockTransactions)):
		return MsgIdentifierRequestMainblockTransactions
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyMainblockTransactions)):
		return MsgIdentifierReplyMainblockTransactions
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierBroadcastMainblock)):
		return MsgIdentifierBroadcastMainblock
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierBroadcastTransaction)):
		return MsgIdentifierBroadcastTransaction
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestData)):
		return MsgIdentifierRequestData
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyData)):
		return MsgIdentifierReplyData
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestDataFile)):
		return MsgIdentifierRequestDataFile
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyDataFile)):
		return MsgIdentifierReplyDataFile
		//case (RawCheckIdentifier(msgidentifier,  )):
		//	return
	}
	return ""
}
func GetContentMaxSize(msgidentifier []byte) int {
	switch {
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestHandshake)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyHandshake)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestMainchainLength)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyMainchainLength)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestMainheaders)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyMainheaders)):
		return 2000000
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestMainblockTransactions)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyMainblockTransactions)):
		return 2000000
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierBroadcastMainblock)):
		return 2000000
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierBroadcastTransaction)):
		return 2000
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestData)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyData)):
		return 2000000
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestDataFile)):
		return 500
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyDataFile)):
		return 2000000
		//case (RawCheckIdentifier(msgidentifier,  )):
		//	return
	}
	return 1000
}

//type Message struct {
//	Identifier []byte
//	Content []byte
//}
//Network 4 bytes
//Version 4 bytes
//Identifier 20 bytes
//Length 4 bytes

/*
type Message str {
	GetIdentifier() []byte
	PutIdentifier(string)
	CheckIdentifier(string) bool
	WriteBytes(connection net.Conn)
	ReadContent(connection net.Conn) error
	PutContent([]byte)
	GetContent() []byte
}
*/

////////////////////////////////
func NewMessage(identifier string) *Message {

	var msg Message

	msg.PutIdentifier(identifier)

	return &msg
}

type Message struct {
	Identifier []byte
	Content    []byte
	//Connection * net.Conn
	OriginPeer *Peer
}

func (msg *Message) GetIdentifier() []byte {
	return msg.Identifier
}
func (msg *Message) PutIdentifier(identifier string) {

	var identifierbuf [20]byte
	copy(identifierbuf[:], []byte(identifier))

	msg.Identifier = identifierbuf[:]
}
func (msg *Message) CheckIdentifier(identifier string) bool {
	if msg.Identifier == nil {
		return false
	}
	return RawCheckIdentifier(msg.Identifier, identifier)
}
func (msg *Message) PutContent(buf []byte) {
	msg.Content = buf
	return
}
func (msg *Message) GetContent() []byte {

	return msg.Content
}

/////////////////////////////////////////

func (msg *Message) WriteBytes(connection net.Conn) {

	tmpbw := utility.NewBufferWriter()
	tmpbw.PutBytes([]byte(MainNetworkIdentifier)) // the network
	tmpbw.PutUint32(1)                            // the message version
	tmpbw.PutBytes(msg.Identifier)
	tmpbw.PutUint32(uint32(len(msg.Content)))
	tmpbw.PutHash(utility.ComputeHash(msg.Content))
	tmpbw.PutBytes(msg.Content)

	var err error
	_, err = connection.Write(tmpbw.GetContent())
	if err != nil {
		applog.Trace("error")
	}
}
func (msg *Message) ReadContent(connection net.Conn, maxsize int) error {
	var err error
	buffcontentlength := make([]byte, 4)
	_, err = connection.Read(buffcontentlength)
	if err != nil {
		applog.Trace("warning unable to read content length%x", buffcontentlength)
		return err
	}

	contentlength := binary.LittleEndian.Uint32(buffcontentlength)
	if int(contentlength) > maxsize {
		applog.Warning("Max size exceeded while reading content of message - contentlength %d maxsize %d ", contentlength, maxsize)
		return fmt.Errorf("Max size exceeded")
	}

	buffcontenthashbytes := make([]byte, utility.HashSize)
	_, err = connection.Read(buffcontenthashbytes)
	if err != nil {
		applog.Trace("warning unable to read content hashbytes")
		return err
	}

	buffcontent := make([]byte, contentlength)
	_, err = connection.Read(buffcontent)
	if err != nil {
		applog.Trace("warning unable to read content")
		return err
	}

	//applog.Trace("****************content length %d content %x",contentlength, buffcontent)
	msg.PutContent(buffcontent)
	if utility.ComputeHash(buffcontent) != *utility.NewHash(buffcontenthashbytes) {
		return nil
	}
	return nil

}

////////////////////////////////

func ReadIdentifier(connection net.Conn) ([]byte, error) {
	var err error
	networkidentifier := make([]byte, 4)
	_, err = connection.Read(networkidentifier)
	if err != nil {
		return nil, err
	}
	//if !=tmpbw.PutBytes([]byte(MainNetworkIdentifier))// the network
	if !bytes.Equal(networkidentifier, ([]byte(MainNetworkIdentifier))) {
		return nil, nil
	}

	msgtype := make([]byte, 4)
	_, err = connection.Read(msgtype)
	if err != nil {
		return nil, err
	}
	msgidentifier := make([]byte, 20)
	_, err = connection.Read(msgidentifier)
	if err != nil {
		return nil, err
	}
	applog.Trace("message network %v type %v identifier %v", networkidentifier, msgtype, msgidentifier)
	return msgidentifier, nil
}

func ReadConnectionMessage(connection net.Conn) (*Message, error) {

	msgidentifier, err := ReadIdentifier(connection)
	if err != nil {
		return nil, err
	}

	if msgidentifier == nil {
		return nil, nil
	}
	buffcontentlength := make([]byte, 4)
	_, err = connection.Read(buffcontentlength)
	if err != nil {
		applog.Trace("warning unable to read content length%x", buffcontentlength)
		return nil, err
	}

	contentlength := binary.LittleEndian.Uint32(buffcontentlength)
	if contentlength != 0 {
		return nil, nil
	}
	//
	buffcontenthashbytes := make([]byte, utility.HashSize)
	_, err = connection.Read(buffcontenthashbytes)
	if err != nil {
		applog.Trace("warning unable to read content hashbytes")
		return nil, err
	}
	//
	var msg *Message

	switch {
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierRequestHandshake)):
		applog.Trace("Good request handshake")
		msg = NewMessage(MsgIdentifierRequestHandshake)
	case (RawCheckIdentifier(msgidentifier, MsgIdentifierReplyHandshake)):
		applog.Trace("Good reply handshake")
		msg = NewMessage(MsgIdentifierReplyHandshake)
	default:
		applog.Trace("Unknown connection message")
		return nil, nil
	}

	return msg, err
}
func (peer *Peer) ReadTCPMessage() (*Message, error) {
	connection := peer.Connection
	msgidentifier, err := ReadIdentifier(connection)
	if err != nil {
		return nil, err
	}

	if msgidentifier == nil {
		return nil, nil
	}

	var msg *Message

	msgidentifierstring := DecodeIdentifier(msgidentifier)
	if msgidentifierstring == "" {
		applog.Trace("error unknown message identifier")
		return nil, nil
	}

	msg = NewMessage(msgidentifierstring)
	rerr := msg.ReadContent(connection, GetContentMaxSize(msgidentifier))
	if rerr != nil {
		applog.Trace("error while reading message")
		return nil, nil
	}
	msg.OriginPeer = peer
	return msg, nil
}

func RawCheckIdentifier(buffer []byte, identifier string) bool {
	var identifierbuf [20]byte
	copy(identifierbuf[:], identifier)
	//applog.Trace("%v %v",buffer,identifierbuf)
	if len(identifierbuf) != 20 {
		return false
	}

	if bytes.Equal(buffer, identifierbuf[:]) {
		return true
	}
	return false
}
