package main

import (
	"chainstack-core/rpc"
	"fmt"
	//"chainstack-monitor/chainstackscan/handler"
	"strings"
	"chainstack-monitor/util"
	//"math/big"
	"chainstack-core/common"
	//"time"
	"chainstack-core/core/accounts"
	//"chainstack-core/core/rpc-interface"
	"io/ioutil"
	//"reflect"
	//"reflect"
	"encoding/hex"
	"time"
	"math/big"
)

const (
	clusterConfFileName = "cluster.json"
	clusterConfigJson = "/home/qydev/tmp/default_cluster.json"
)

type NodeWallet struct {
	HttpPort int
	Account string
}

type NodeConf struct {
	NodeName string `json:"node_name"`
	P2PListener string `json:"p2p_listener"`
	HttpPort int `json:"http_port"`
	WsPort int `json:"ws_port"`
	Host string `json:"host"`
}

func GetClusterConfig() (configure []NodeConf,err error){

	fmt.Println("thg confFile path is:", clusterConfigJson)
	// 尝试从配置文件中加载配置
	if fb, err := ioutil.ReadFile(clusterConfigJson); err == nil {
		if err = util.ParseJsonFromBytes(fb, &configure); err != nil {
			return []NodeConf{},err
		}
	} else {
		return []NodeConf{},err
	}

	fmt.Println("the node number is:",len(configure))

	/*      verifierNumber := len(configure.VerifyNodes)
			normalNumber := len(configure.NormalNodes)
			minMasterNumber := len(configure.MineMasterNodes)
			MineWorker := len(configure.MineWorkers)

			log.Debug("[TestChangeVerifier]the cluster verifier number is:","number",verifierNumber)
			log.Debug("[TestChangeVerifier]the cluster normal number is:","number",normalNumber)
			log.Debug("[TestChangeVerifier]the cluster mineMaster number is:","number",minMasterNumber)
			log.Debug("[TestChangeVerifier]the cluster MineWorker number is:","number",MineWorker)*/
	return configure,nil
}

func newRpcClientConnect(port int) *rpc.Client {
	client, err := rpc.Dial(fmt.Sprintf("http://%v:%d", "127.0.0.1", port))
	if err != nil {
		panic(err.Error())
	}
	return client
}

//查看node address
func CheckAddressList() map[string]NodeWallet{
	configure, _ := GetClusterConfig()

	txAccounts := make(map[string]NodeWallet)
	for i:=0; i < len(configure); i++ {
		port := configure[i].HttpPort
		name := configure[i].NodeName
		client := newRpcClientConnect(port)
		path := fmt.Sprintf("/home/qydev/tmp/chainstack_apps/%v/CSWallet", name)
		WalletResp := &accounts.WalletIdentifier{
			WalletType: accounts.SoftWallet,
			Path:       path,
			WalletName: "CSWallet",
		}
		var wListResp []accounts.Account

		if err := client.Call(&wListResp, getRpcTXMethod("ListWalletAccount"),
			WalletResp);
			err != nil {
			fmt.Println("send transaction error==", err)
		} else {
			fmt.Println(name, " ", port, " ", wListResp)
			account := hex.EncodeToString(wListResp[0].Address[:])
			//fmt.Println(account)
			if name == "default_m0" || name == "default_n0" {
				txAccounts[name] = NodeWallet{port, account}
			}
			//_, ok := wListResp.([]accounts.Account)
			//fmt.Println(ok)
			//switch v := wListResp.(type) {
			//case []accounts.Account:
			//	fmt.Println(v)
			//default:
			//	fmt.Println(reflect.ValueOf(wListResp).Type())
			//	fmt.Println(reflect.ValueOf(wListResp))
			//	fmt.Println(v)
			//}
		}
	}
	fmt.Println(txAccounts)
	return txAccounts
}

//发交易RPC方法
func getRpcTXMethod(methodName string) string{
	return "chainstack_" + strings.ToLower(methodName[0:1]) + methodName[1:]
}

func TrySendTx(accounts map[string]NodeWallet, round int, doubleExpend bool) {
	clientB := newRpcClientConnect(accounts["default_m0"].HttpPort)

	var cNonce uint64

	if nonceErr := clientB.Call(&cNonce,getRpcTXMethod("GetTransactionNonce"),
		common.HexToAddress(accounts["default_m0"].Account)); nonceErr!=nil{
		fmt.Println("get current nonce error==", nonceErr)
	}else {
		fmt.Println("Current Nonce==",cNonce)
	}
	if doubleExpend && cNonce > 0 {
		//make nonce equal to last tx for double expend
		//cNonce--
		cNonce = 0
	}
	for i:=0; i<round; i++ {
		select {
		case <-time.After(time.Nanosecond*1000000):
			/*
			1秒=1e3毫秒=1e9纳秒
			TPS=10k，1e9/1e4=1e5，每隔1e5纳秒发一笔交易
			TPS=1k ，1e9/1e3=1e6，每隔1e6纳秒（1毫秒）发一笔交易
			以此类推
			*/

			var sendTxResp interface{}

			if err := clientB.Call(&sendTxResp, getRpcTXMethod("SendTransaction"),
				common.HexToAddress(accounts["default_m0"].Account),
				common.HexToAddress(accounts["default_n0"].Account),
				big.NewInt(100+int64(i*2)),
				big.NewInt(100+int64(i*2)),
				nil,
				cNonce+uint64(i));
				err != nil {
				fmt.Println("send transaction error==", err)
			}else {
				fmt.Println("B-A=====NO.",round,"--",sendTxResp)
			}
		}
	}

	fmt.Println("B-A交易发送完毕")
}

func main() {
	accounts := CheckAddressList()
	TrySendTx(accounts, 1, true)
}