package model

import (
	"testing"
	"runtime"
	"fmt"
	"time"
	"math/big"
	"github.com/zcheng9/nipo/crypto"
	"chainstack-core/common/util"
)

type addRet struct {
	v int
	ok bool
}

func addInt(i interface{}) interface{} {
	ret := addRet{0, false}
	if v, ok := i.(int); ok {
		//fmt.Println("---", v)
		ret.v = v+1
		ret.ok = true
	}
	time.Sleep(time.Microsecond*1)
	return ret
}

func Test_Map_Work(t *testing.T) {
	LENGTH := 40000

	mWorkMap := InitWorkMap(runtime.NumCPU())

	mWorkMap.SetOperation(addInt)

	raw := make([]interface{}, LENGTH)
	after := make([]addRet, LENGTH)

	for i:=0; i<LENGTH; i++ {
		raw[i] = i
	}


	retCh, err := mWorkMap.StartWorks(raw)

	j:=0

	st := time.Now()
	if err == nil {
		pass:
		for {
			select {
			case r := <- retCh:
				after[j] = r.(addRet)
				j++
				//fmt.Println(r, j)
				if j == LENGTH {
					//fmt.Println("work processed")
					break pass
				}
			}
		}
	}

	fmt.Println(time.Now().Sub(st))

	st = time.Now()
	for i:=0; i<LENGTH; i++ {
		r := addInt(i)
		after[i] = r.(addRet)
	}
	fmt.Println(time.Now().Sub(st))
}

func fakeVerify(item interface{}) interface{} {
	ms := NewNakamotoSigner(big.NewInt(1))
	tx := item.(*Transaction)
	ret := false
	if ms.VerifySender(*tx) {
		ret = true
	}
	return ret
}

func fakeFilter(inTxs []*Transaction, valids []bool) (txs []*Transaction) {
	for i, v := range valids {
		if v {
			txs = append(txs, inTxs[i])
		}
	}
	return
}

func TestTxVerify(t *testing.T) {
	txs:=createTxList(1000)

	key1, key2 := createKey()
	fs1 := NewNakamotoSigner(big.NewInt(1))
	alice := crypto.GetNormalAddress(key1.PublicKey)
	bob := crypto.GetNormalAddress(key2.PublicKey)
	badtx := NewTransaction(10, alice,bob, big.NewInt(10000), big.NewInt(10), []byte{})
	badtx.SignTx(key2, fs1)
	txs=append(txs,badtx)

	LENGTH := len(txs)
	after := make([]bool, LENGTH)

	mWorkMap := InitWorkMap(runtime.NumCPU())

	mWorkMap.SetOperation(fakeVerify)
	j := 0
	st := time.Now()
	//convert to interface slice
	ifTxs := make([]interface{}, len(txs))
	util.InterfaceSliceCopy(ifTxs, txs)
	//start calculate
	retCh, err := mWorkMap.StartWorks(ifTxs)

	if err == nil {
	pass:
		for {
			select {
			case r := <- retCh:
				after[j] = r.(bool)
				j++
				//fmt.Println(r, j)
				if j == LENGTH {
					//fmt.Println("work processed")
					break pass
				}
			}
		}
	}

	fv := fakeFilter(txs, after)
	fmt.Println(time.Now().Sub(st))
	fmt.Println(len(fv), LENGTH)
}