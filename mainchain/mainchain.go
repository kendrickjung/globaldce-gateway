package mainchain

import (
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/utility"
	"math/big"
	"math"
	"fmt"
	//"os"
)

const (
	GENESIS_BLOCK_REWARD uint64 = 50000000000000 // 50 000 000 globals * 1 000 000 
	BLOCK_TIME uint32= 600//in seconds// ten minutes

	DIFFICULTY_TUNING_INTERVAL uint32=2016//about two weeks 
	BLOCK_REWARD_TUNING_INTERVAL uint32=210000//blocks
	//BLOCK_CONFIRMATION_INTERVAL uint32=100//6//blocks
	ENGAGEMENT_REWARD_FINALIZATION_INTERVAL uint32=144//24*6//about one day
	BLOCK_MAX_SEIZE uint32=1024*1024//1 MO
	BLOCKTRANSACTIONS_MAX_SEIZE uint32=BLOCK_MAX_SEIZE-200
	OBSTRUCTED_MINING_TIME int64=60*600//10 hours
	//ENGAGEMENT_PUBLICPOST_MAXSTAKE uint64 = 1000000000
	TRACKING_ADDRESSES=true

)

func  GetMainblockReward(i uint32) uint64{
	return GENESIS_BLOCK_REWARD/uint64(math.Pow(2, float64(i/BLOCK_REWARD_TUNING_INTERVAL)))
}
func (mn *Maincore) GetTargetBits() uint32{
	var targetbits uint32
	mainchainlength:=int(mn.GetMainchainLength())
	if (mainchainlength % int ( DIFFICULTY_TUNING_INTERVAL)==0){
		
		targetbigint:=utility.BigIntFromCompact( mn.GetLastMainheader().Bits)
		idealtimeinterval:=int64 ( DIFFICULTY_TUNING_INTERVAL-1)*600
		realtimeinterval:=int64 ( mn.GetLastMainheader().Timestamp- mn.GetMainheader(mainchainlength-int ( DIFFICULTY_TUNING_INTERVAL)).Timestamp)
		
		if (realtimeinterval>=int64 (3)*idealtimeinterval) {
			targetbigint.Mul(targetbigint,big.NewInt(3))
			applog.Trace("bigger than 3  ")
		} else if (idealtimeinterval>=int64 (3) *realtimeinterval) {
			targetbigint.Div(targetbigint,big.NewInt(3))
			applog.Trace("smaller that 1/3 ")
		} else {
			targetbigint.Mul(targetbigint,big.NewInt(realtimeinterval))
			targetbigint.Div(targetbigint,big.NewInt(idealtimeinterval))
		}

		targetbits = utility.CompactFromBigInt(targetbigint)

		applog.Trace("\n realtime %d idealtime %d  targetbitint %d ",realtimeinterval,idealtimeinterval,targetbigint)
	} else {
		
		targetbits = mn.GetMainheader(mainchainlength-1).Bits

	}
	return targetbits
}
func  (mn *Maincore)  ValidatePropagatingMainblock(mb *Mainblock) (bool) {

	// if the mainblock height is < mn.GetConfirmedMainchainLength() than the mainblock is rejected -return false
	if mb.Height<mn.GetConfirmedMainchainLength(){
		applog.Warning("Rejected Propagating block - mb.Height (%d) < mn.GetConfirmedMainchainLength() (%d)",mb.Height,mn.GetConfirmedMainchainLength())
		return false
	}
	// if the mainblock height is < mn.GetMainchainLength() than the mainblock is a parallel block -
	if mb.Height<mn.GetMainchainLength() {
		//TODO check parallel block and even store it
		applog.Warning("Rejected Propagating block - parallel propagating block")
		return false
	}
	if mb.Height!=mn.GetMainchainLength() {
		applog.Warning("Rejected Propagating block - mb.Height != mn.GetConfirmedMainchainLength()")
		return false
		//
	}
	// check the header
	if !mn.CheckPropagatingMainheader(&(mb.Header),mb.Height){
		applog.Warning("Rejected Propagating block - inncorrect header")
		return false
	}
	// check the transactions
	if !CheckMainblockTransactions(&(mb.Transactions),mb.Header.Roothash){
		applog.Warning("Rejected Propagating block - incorrect transactions")
		return false
	}
	// validate every transaction
	if !mn.ValidateMainblockTransactions(mb.Height,&(mb.Transactions)){
		applog.Warning("Rejected Propagating block - invalid transactions")
		return false
	}
	return true
}

func (mn *Maincore) CheckPropagatingMainheader(header *Mainheader,height uint32) bool{
	prevmainheader :=mn.GetMainheader(int(height)-1)

	applog.Trace("Checking header %d %x",height,header)
	if (header.Prevblockhash!=prevmainheader.Hash){
		applog.Trace("\n error: blockheader %d - Prevblockhash do not match previous block hash ",height)
		return false
	}
	if (header.Timestamp<prevmainheader.Timestamp){
		applog.Trace("\n error: blockheader %d - Timestamp precede previous block timestamp ",height)
		return false
	}
	


	//targetbigint:=utility.CorrectTargetBigInt(utility.BigIntFromCompact(header.Bits),header.Timestamp,prevmainheader.Timestamp)
	targetbigint:=utility.BigIntFromCompact(header.Bits)
	if (!(utility.BigIntFromHash(&header.Hash).Cmp(targetbigint)<0)){
		applog.Trace("\n error: hash of header %d do not fall into its own target ",height)
		return false
	}
	if (height==mn.GetMainchainLength()) {
		if (header.Bits!=mn.GetTargetBits()){
			applog.Trace("\n error: blockheader %d - Propagating mainblock target do not match required block target ",height)
			return false
		}
	} else if (height<mn.GetMainchainLength()) {
		if (header.Bits!=mn.GetMainheader(int(height)).Bits){
			applog.Trace("\n error: blockheader %d - Propagating mainblock (parallel mainblock) target do not match required block target ",height)
			return false
		}
	} else {
		applog.Trace("\n error: blockheader %d - Propagating mainblock is too far ahead - mainchain length is %d",height,mn.GetMainchainLength())
		return false
	}


return true
}

func  (mn *Maincore)  CheckHeaderChain(pointerheaders *[]Mainheader) bool {
	var firstheaderid int =1
	if len(*pointerheaders)!=0{
		firstheaderid=len(mn.mainheaders)
	} 
	
	headers:=*pointerheaders
	headers=append(mn.mainheaders,headers ...)
	for i := firstheaderid; i < len(headers); i++ {

		applog.Trace("header %d %x",i,headers[i])
		if (headers[i].Prevblockhash!=headers[i-1].Hash){
			applog.Trace("\n error: blockheader %d - Prevblockhash do not match previous block hash ",i)
			return false
		}
		if (headers[i].Timestamp<headers[i-1].Timestamp){
			applog.Trace("\n error: blockheader %d - Timestamp precede previous block timestamp ",i)
			return false
		}
		//targetbigint:=utility.CorrectTargetBigInt(utility.BigIntFromCompact(headers[i].Bits),headers[i].Timestamp,headers[i-1].Timestamp)
		targetbigint:=utility.BigIntFromCompact(headers[i].Bits)
		if (!(utility.BigIntFromHash(&headers[i].Hash).Cmp(targetbigint)<0)){
			applog.Trace("\n error: hash of header %d do not fall into its own target ",i)
			return false
		}
		if ((i ) % int (DIFFICULTY_TUNING_INTERVAL)!=0) {
			if (headers[i].Bits!=headers[i-1].Bits){
				applog.Trace("\n error: blockheader %d - Block target do not match previous block target ",i)
				return false
			}
		} else {
			var targetbits uint32
			targetbigint:=utility.BigIntFromCompact(headers[i-1].Bits)
			idealtimeinterval:=int64 (DIFFICULTY_TUNING_INTERVAL-1)*600
			realtimeinterval:=int64 (headers[i-1].Timestamp-headers[i-int (DIFFICULTY_TUNING_INTERVAL)].Timestamp)
			
			if (realtimeinterval>=int64 (3)*idealtimeinterval) {
				targetbigint.Mul(targetbigint,big.NewInt(3))
				//applog.Trace("bigger than 3  ")
			} else if (idealtimeinterval>=int64 (3) *realtimeinterval) {
				targetbigint.Div(targetbigint,big.NewInt(3))
				//applog.Trace("smaller that 1/3 ")
			} else {
				targetbigint.Mul(targetbigint,big.NewInt(realtimeinterval))
				targetbigint.Div(targetbigint,big.NewInt(idealtimeinterval))
			}

			targetbits = utility.CompactFromBigInt(targetbigint)

			applog.Trace("\n realtime %d idealtime %d  targetbitint %d ",realtimeinterval,idealtimeinterval,targetbigint)
			if (targetbits!=headers[i].Bits){
				applog.Trace("\n error: blockheader %d - Block target do not match computed block target ",i)
				return false
			}
		}
	}
	applog.Trace("\n notice: headers verified and found correct ")
	return true
}

func CheckMainblockTransactions(pointertxs *[]utility.Transaction,root utility.Hash) bool {

	if len(*pointertxs)==0 {
		applog.Trace("\n Error: empty transactions array")
		return false
	}
	txs:=*pointertxs
	if len(txs[0].Vin)!=0{
		applog.Trace("\n Error: Reward Transactions must not have inputs")
		return false
	}
	if len(txs[0].Vout)!=1{
		applog.Trace("\n Error: Reward Transactions must a have reward output")
		return false
	}

	hashes:= make([]utility.Hash, len(txs))
	for i:=0;i<len(txs);i++{
		hashes[i]=txs[i].ComputeHash()
	}
	if root!=utility.ComputeRoot(&hashes){
		applog.Trace("\n Error: Header root does not match transactions root ")
		return false
	}

	//
	return true
}

func (mn *Maincore) ValidateMainblockTransactions(height uint32,pointertxs *[]utility.Transaction) bool {
	txs:=*pointertxs
	var totalfees uint64=0
	for i:=1;i<len(txs);i++{
		//
		validity,txfee:=mn.ValidateTransaction(&txs[i])
		if !validity{
			applog.Warning("invalid mainblock - invalid transaction %d",i)
			return false
		}
		totalfees+=txfee
	}

	if txs[0].Vout[0].Value!=totalfees+GetMainblockReward(height){
		applog.Warning("invalid mainblock - reward fees %d do not match total transactions fees %d plus block reward",txs[0].Vout[0].Value,uint64 (totalfees)+GetMainblockReward(height))
		return false
	}
	return true
}

func (mn *Maincore) ValidateNameRegistration(name []byte)(error){
	if len(name)>utility.RegistredNameMaxSize{
		return fmt.Errorf("Name length exceeds RegistredNameMaxSize - Name length %d RegistredNameMaxSize %d",len(name))
	}
	if !CheckNameBytes(name){
		return fmt.Errorf("Names can only have lower case letters and numbers")
	}
	if mn.GetNameState(name)==StateValueIdentifierActifNameRegistration{
		return fmt.Errorf("Names is already taken")
	}
	return nil
}


func CheckNameBytes(name []byte) bool {
	for _,n :=range name {
		// only lower case letters and numbers are allowed
		if !(((n>=97)&&(n<=122))||((n>=48)&&(n<=57))) {
			return false
		}
	}
	return true
}