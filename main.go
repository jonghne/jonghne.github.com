package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"chainstack-core/crypto"
	"chainstack-core/p2p/discover"
	"chainstack-core/common/util"
	util2 "chainstack-core/common/util"
	"math/big"
	"chainstack-core/common"
	"chainstack-core/common/consts"
	"chainstack-core/core/chain-config"
	"chainstack-core/core/chain"
	"chainstack-core/core/chain/chaindb"
	"chainstack-core/core/chain/state_processor"
	"flag"
)

const (
	one = iota
	two
	three
)

func partition(buck []int, length int) int {
	pivot := buck[0]
	left := 0
	right := length - 1
	index := left
	for left < right {
		for ; left < right; right-- {
			if pivot > buck[right] {
				buck[index] = buck[right]
				left++
				index = right

				break
			}
		}

		for ; left < right; left++ {
			if pivot < buck[left] {
				buck[index] = buck[left]
				right--
				index = left
				break
			}
		}

	}
	buck[index] = pivot
	return index
}

func qsort(buck []int, length int) {
	fmt.Println(len(buck), buck)
	if length < 2 {
		fmt.Println("end")
		return
	}
	index := partition(buck, length)
	fmt.Println("index=", index)
	qsort(buck[:index], index)
	fmt.Println("===========")
	qsort(buck[index+1:], length-index-1)
}

func search(buck []int, w int) int {
	var cnt int = 0
	var length int = len(buck)
	var i int
	var invalid int = 0
	for invalid != length {
		var ary []int
		j := 0
		fmt.Println("*********")
		for i = 0; i < length; {
			fmt.Println("i=", i)
			if buck[i] != 0x7fffffff {
				for k := i + 1; k < length; k++ {
					fmt.Println("k=", k)
					if buck[k] != 0x7fffffff {
						if buck[k] == (buck[i] + 1) {
							fmt.Println("find ", i, k)
							ary = append(ary, buck[i])
							j++
							buck[i] = 0x7fffffff
							invalid++
							i = k
							break
						} else if buck[k] == buck[i] {
							i++
							continue
						} else {
							fmt.Println("nothing error ", ary, buck, i, buck[i], k, buck[k])
							return 0
						}
					}
				}
			} else {
				i++
			}
			if j == (w - 1) {
				cnt++
				ary = append(ary, buck[i])
				buck[i] = 0x7fffffff
				invalid++
				fmt.Println(ary)
				break
			}
		}
		if i >= length {
			fmt.Println("scan end error ", ary)
			return 0
		}
	}
	return cnt
}

type cat struct {
	race  string
	color string
}

type eat interface {
	eat()
	set(s string)
}

type rest interface {
	sleep()
}

func (a *cat) eat() {
	fmt.Printf("%s cat eat\n", a.race)
}

func (a *cat) set(s string) {
	a.race = s
}

func (a *cat) sleep() {
	fmt.Printf("cat %s sleep\n", a.color)
}

type home struct {
	eater eat
}

func setHomeCat(e eat) {
	hm := home{e}
	hm.eater.set("moxian")
}

func call(obj interface{}) {
	switch c := obj.(type) {
	case rest:
		c.sleep()
	case eat:
		c.eat()
	default:
		fmt.Println("unrecognized")
	}
}

func qq() {
	fmt.Println("qq")
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func CheckNodeIsAlive(info time.Time) bool {
	now := time.Now()
	nodeReportLastTime := info.Add(30*time.Second)
	return now.Before(nodeReportLastTime)
}

type Animal struct {
	Name  string
	Order string
	Pad   string
}

type genesisCfgFile struct {
	Nonce uint64 `json:"nonce"`
	Note string `json:"note"`
	Accounts []string `json:"accounts"`
	Balances []int64 `json:"balances"`
	Timestamp string  `json:"timestamp"`
	Difficulty string `json:"difficulty" gencodec:"required"`
	Verifiers  []common.Address `json:"verifiers" gencodec:"required"`
}

func GenesisBlockFromFile(chainDB chaindb.Database, accountStateProcessor state_processor.AccountStateProcessor) *chain.Genesis {

	ge, e := ioutil.ReadFile(filepath.Join(util2.HomeDir(), "softwares/chainstack_deploy/genesis.json"))
	if e != nil {
		return nil
	}

	var info genesisCfgFile
	fmt.Printf("ge: %v\n", string(ge))
	err := json.Unmarshal(ge, &info)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", info)

	var gTime time.Time
	if gTime, err = time.Parse("2006-01-02 15:04:05", info.Timestamp); err != nil {
		gTime, _ = time.Parse("2006-01-02 15:04:05", "2018-08-08 08:08:08")
	}

	alloc := make(map[common.Address]*big.Int)

	addrLen := len(info.Accounts)
	bLen := len(info.Balances)

	if addrLen != bLen {
		return nil
	}
	for i := range info.Accounts {
		alloc[common.HexToAddress(info.Accounts[i])] = big.NewInt(info.Balances[i] * consts.CSK)
	}

	return &chain.Genesis{
		//chainDB:               chainDB,
		//accountStateProcessor: accountStateProcessor,
		Config:                chain_config.GetChainConfig(),
		Nonce:      info.Nonce,
		Timestamp:  big.NewInt(gTime.UnixNano()),
		ExtraData:  []byte(info.Note),
		Difficulty: common.HexToDiff(info.Difficulty),
		Alloc: alloc,
		//verifiers: info.Verifiers,
	}
}

var (
	bootnodeIp = flag.String("ip", "127.0.0.1", "bootnode IP")
)
func main() {
	jc := 2
	jc1 := 3
	switch {

	case jc ==2 && jc1 == 3:
		fmt.Println("2 conditions")
	case jc == 2:
		fmt.Println("jc==2")
	}
	//genbootnodeid
	confB, e := ioutil.ReadFile(filepath.Join(util2.HomeDir(), "softwares/chainstack_deploy/bootnode_key"))
	if e == nil {
		fmt.Println(confB, string(confB))


		flag.Parse()
		fmt.Println("------", *bootnodeIp)
		pk, _ := crypto.HexToECDSA(string(confB))
		n := discover.PubKeyToNodeID(&pk.PublicKey)
		fmt.Println(n.String())
		idStr := n.String()
		host := *bootnodeIp
		ns := []string{}
		ns = append(ns, fmt.Sprintf("enode://%v@%v:30301", idStr, host))
		fmt.Println(ns)
		// write boot node info to datadir
		ioutil.WriteFile(filepath.Join(util2.HomeDir(), "softwares/chainstack_deploy/static_boot_nodes.json"), util.StringifyJsonToBytes(ns), 0755)
	}
	////////////////////////////////
	//genesis
	mg := GenesisBlockFromFile(nil, nil)
	fmt.Printf("%v\n", mg)

	gTime, _ := time.Parse("2006-01-02 15:04:05", "2018-08-08 08:08:08")
	fmt.Println(big.NewInt(gTime.UnixNano()))
	fmt.Println(common.HexToDiff("dead"))
	//////////////////////
	//generate account
	nodeKey, _ := crypto.GenerateKey()
	addr := crypto.GetNormalAddress(nodeKey.PublicKey)
	fmt.Println(addr, addr.String(), addr[:])
	fmt.Printf("%v %v %v\n",addr, addr.String(), addr.Str())

	var jsonBlob = []byte(`[  
    {"Name": "Platypus", "Order": "Monotremata", "pad":"jim"},  
    {"Name": "Quoll",    "Order": "Dasyuromorphia"}  
]`)

	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", animals)

	node := []byte(`["enode://9fb88505dca8bbed224e03e9c0c27f724bf12551003576c79aae38307d8a4b68a88ad6afa38ae87a0aadaeea4940da2b859ffb0ce6b4b94252de1acc035d2396@14.17.65.122:30301"]`)
	var nodesStr []string
	err = json.Unmarshal(node, &nodesStr)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(nodesStr)

	ml:=[]int{1,2,3,4,5,6,7,8,9,10}
	fromV := reflect.ValueOf(ml)
	tV := reflect.TypeOf(ml)
	fmt.Println(tV, fromV.Index(4).Interface())

	for j:=0;j<len(ml);j++ {
		if i:=j%2;i==0 {
			ml = append(ml[:j], ml[j+1:]...)
		}
	}
	//for j:= range ml {
	//	if i:=j%2;i==0 {
	//		ml = append(ml[:j], ml[j+1:]...)
	//	}
	//}
	fmt.Println(ml)

	strinit := []string{"172.16.110.236"}
	for i, v:= range strinit {
		fmt.Println("str list ", i, v)
	}
	emptycat := cat{}
	fmt.Println(emptycat)
	personSalary := map[string]int{
		"steve": 12000,
		"jamie": 15000,
	}
	personSalary["mike"] = 9000
	fmt.Println("All items of a map", len(personSalary))
	for key, value := range personSalary {
		fmt.Printf("personSalary[%s] = %d\n", key, value)
	}
	for key := range personSalary {
		fmt.Printf("personSalary[%s]\n", key)
	}
	fmt.Println(personSalary["lucy"])

	scratch := time.Now()
	time.Sleep(11*time.Second)
	if CheckNodeIsAlive(scratch) {
		fmt.Println("alive")
	} else {
		fmt.Println("timeout")
	}
	time.Sleep(20*time.Second)
	if CheckNodeIsAlive(scratch) {
		fmt.Println("alive")
	} else {
		fmt.Println("timeout")
	}
	p := cat{"shen", "red"}
	p.sleep()
	p.eat()
	call(&p)
	setHomeCat(&p)
	p.eat()

	fmt.Printf("%x\n", 1000000000)
	var ary [3]int
	a1 := ary[1:2]
	a1[0] = 12
	fmt.Println(a1, cap(a1), len(a1))
	fmt.Println(ary)
	fmt.Println(123)
	fmt.Println(one, two, three)

	fmt.Println(strconv.Atoi(os.Args[2]))
	fmt.Println(len(os.Args))
	var b interface{} = os.Args[3]
	switch q := b.(type) {
	case int:
		fmt.Println("integer", q)
	case string:
		fmt.Println("string", q)
	default:
		fmt.Println("other type")
	}

	var sc []int
	//sc=append(sc, 18)
	w, _ := strconv.Atoi(os.Args[1])
	for _, v := range os.Args[2:] {
		r, _ := strconv.Atoi(v)
		fmt.Println(r)
		sc = append(sc, r)
	}
	fmt.Println(sc, reflect.TypeOf(sc))
	qsort(sc, len(sc))
	fmt.Println(sc)
	if (len(sc) % w) > 0 {
		fmt.Println("not match, w error")
		//return
	}
	//count := search(sc, w)
	//fmt.Println(count)

	var sl []int
	fmt.Println("#########")
	print(sl)
	fmt.Println("#########")
	sl = append(sc)
	fmt.Println(sl)
	sl = append(sl[0:0])
	fmt.Println(sl)

	c1 := []byte("hello")
	c2 := []byte("hello12")
	if !bytes.Equal(c1[:], c2[:]) {
		fmt.Println("equal")
	}
	i := true
	if !i {
		fmt.Println("oppose")
	}

	fmt.Println("error: ", errors.New("coin base is invalid"))

	type INTARY [4]int
	ry := INTARY{12}
	fmt.Println(ry)

	str2 := "hello123"
	data2 := []byte(str2)
	fmt.Println(data2)

	select {}
	fmt.Println("________________")
	ch := make(chan int, 24)
	go func() {
		for {
			val := <-ch
			fmt.Printf("val=%d\n", val)
		}

	}()
	tick := time.NewTicker(1 * time.Second)
	for i := 0; i < 20; i++ {
		select {
		case ch <- i:
			//case <-tick.C:
			//	fmt.Printf("%d: case <-tick.C\n", i)
		}
		select {
		case <-tick.C:
			fmt.Printf("%d: case <-tick.C\n", i)
		default:
		}
		time.Sleep(200 * time.Millisecond)
	}
}
