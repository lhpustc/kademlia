package libkademlia

import (
	//"bytes"
	"fmt"
	"net"
	"strconv"
	"testing"
	//"container/heap"
	// "time"
)

func StringToIpPort(laddr string) (ip net.IP, port uint16, err error) {
	hostString, portString, err := net.SplitHostPort(laddr)
	if err != nil {
		return
	}
	ipStr, err := net.LookupHost(hostString)
	if err != nil {
		return
	}
	for i := 0; i < len(ipStr); i++ {
		ip = net.ParseIP(ipStr[i])
		if ip.To4() != nil {
			break
		}
	}
	portInt, err := strconv.Atoi(portString)
	port = uint16(portInt)
	return
}

/*
	Test Iterative find node

                       1    1
                     /    /
			        /    /
		          B -- 2 -- 2
		        /   \  . \  .
		       /     \ .  \ .
		     A -- C    5    5
		       \  .
		        \ .
                  F

	We build such architecture of out network.
	We let every node save its neibour in its k-bucket.
	We just set one bit as 1 to distinguish every nodeID, at most, there is
	160 nodeIDs, 0 - 159

	We just set A as 0, B - F are A * 5 + i(i = 1-5)
	Following the order, If we want to find the the neighbour list of Node 10 from Node 0
	It will find:
	0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19

	We just check whether the shortlist returned contains such nodes, if it is, test passed
	or if not, faield.


*/

// func TestIterativeFindNode(t *testing.T) {
//
//   fmt.Println("TestIterativeFindNode")
//   begin_port := 6482
// 	instance2 := NewKademliaWithId("localhost:"+strconv.Itoa(begin_port-1),SetBitID(0))
// 	tree_node := make([]*Kademlia,5)
// 	count := 0
// 	for i := 1; i < 6; i++ {
// 		address := "localhost:" + strconv.Itoa(begin_port+count)
// 		count++
// 		tree_node[i-1] = NewKademliaWithId(address,SetBitID(i))
// 		host_number, port_number, _ := StringToIpPort(address)
// 		instance2.DoPing(host_number, port_number)
//
// 		for j := 1; j < 6; j++ {
// 			address_v := "localhost:" + strconv.Itoa(begin_port+count)
//
// 			count++
// 			instance_temp := NewKademliaWithId(address_v,SetBitID((5*i+j)))
// 			host_number_v, port_number_v, _ := StringToIpPort(address_v)
// 			tree_node[i-1].DoPing(host_number_v, port_number_v)
// 			for m := 1; m < 6; m++ {
// 				address_r := "localhost:" + strconv.Itoa(begin_port+count)
// 				count++
// 				 NewKademliaWithId(address_r,SetBitID((5*(5*i+j)+m)))
// 				host_number_r, port_number_r, _ := StringToIpPort(address_r)
// 				instance_temp.DoPing(host_number_r, port_number_r)
// 			}
// 		}
// 	}
//
// 	//Implement DoIterativeFindNode
// 	key := SetBitID(10)
// 	contacts, err := instance2.DoIterativeFindNode(key)
//   if err != nil || len(contacts) != 20{
//   	t.Error("Error doing DoIterativeFindNode")
//   }
//
//   contacts = SortKNodes(SetBitID(0),contacts)
//   fmt.Println("print result")
//   for i:=0;i<len(contacts);i++ {
//     fmt.Println(contacts[i].NodeID)
//     if !contacts[i].NodeID.Equals(SetBitID(i)) {
//       t.Error("Error Result of DoIterativeFindNode")
//     }
//   }
// 	return
// }
//
// /*
// 	Test Iterative store
//
//                        1    1
//                      /    /
// 			        /    /
// 		          B -- 2 -- 2
// 		        /   \  . \  .
// 		       /     \ .  \ .
// 		     A -- C    5    5
// 		       \  .
// 		        \ .
//                   F
//
// 	We buile such architecture of out network.
// 	We let every node save its neibour in its k-bucket.
// 	We just set one bit as 1 to distinguish every nodeID, at most, there is
// 	160 nodeIDs, 0 - 159
//
// 	We just set A as 0, B - F are A * 5 + i(i = 1-5)
// 	We just save a {value: helloworld, ID: nodeID +1}
// 	We just let node 10 find the data locally
//
// 	We just check whether the value returned is "hello world", if it is, test passed
// 	or if not, faield.
//
// */
// func TestIterativeStore(t *testing.T) {
//
//   fmt.Println("TestIterativeStore")
//   // time.Sleep(time.Second*1)
//   begin_port := 3482
//   tree_node := make([]*Kademlia,160)
// 	tree_node[0] = NewKademliaWithId("localhost:"+strconv.Itoa(begin_port-1),SetBitID(0))
//   //Build the Tree Structure
// 	count := 0
// 	for i := 1; i < 6; i++ {
// 		address := "localhost:" + strconv.Itoa(begin_port+count)
// 		count++
// 		tree_node[i] = NewKademliaWithId(address,SetBitID(i))
// 		host_number, port_number, _ := StringToIpPort(address)
// 		tree_node[0].DoPing(host_number, port_number)
//
// 		for j := 1; j < 6; j++ {
// 			address_v := "localhost:" + strconv.Itoa(begin_port+count)
//
// 			count++
// 			tree_node[5*i+j] = NewKademliaWithId(address_v,SetBitID((5*i+j)))
// 			host_number_v, port_number_v, _ := StringToIpPort(address_v)
// 			tree_node[i].DoPing(host_number_v, port_number_v)
// 			for m := 1; m < 6; m++ {
// 				address_r := "localhost:" + strconv.Itoa(begin_port+count)
// 				count++
// 				tree_node[5*(5*i+j)+m] = NewKademliaWithId(address_r,SetBitID((5*(5*i+j)+m)))
// 				host_number_r, port_number_r, _ := StringToIpPort(address_r)
// 				tree_node[5*i+j].DoPing(host_number_r, port_number_r)
// 			}
// 		}
// 	}
//
// 	//Implement DoIterativeFindNode
// 	key := SetBitID(10)
//   key[19] = key[19] ^ byte(1)
//   value := []byte("Hello world")
// 	contacts, err := tree_node[0].DoIterativeStore(key,value)
//   if contacts == nil {
//     t.Error("Failed DoIterativeStore: No Contacts")
//   }
//   // PrintList(contacts)
//   res_v, err := tree_node[10].LocalFindValue(key)
//   if err != nil ||  string(res_v) != string(value){
//     t.Error("Error DoIterativeStore")
//   }
// 	return
// }

/*

TestIterativeFindValue:


                       1    1
                     /    /
			        /    /
		          B -- 2 -- 2
		        /   \  . \  .
		       /     \ .  \ .
		     A -- C    5    5
		       \  .
		        \ .
                  F



Create a three-layer tree structure, and each node(server) have five children (156 nodes in total).
Each node that is not leaf node DoPing to its children.
The i-th node have a NodeID of which the i-th bit is set
(e.g. tree_node[10].NodeID = [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 0]),
and its children is tree_node[5i+1:5i+5].
Then we do a DoIterativeStore at tree_node[0] with
key=[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 1] and value = “Hello World”,
and perform a DoIterativeFindValue to check if the value is stored in tree_node[10]
(NodeID=[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 0]).
Result: “Hello World”. Correct.




*/
// func TestIterativeFindValue(t *testing.T) {
//
//   fmt.Println("TestIterativeFindValue")
//   // time.Sleep(time.Second * 1)
//   begin_port := 4482
//   tree_node := make([]*Kademlia,160)
// 	tree_node[0] = NewKademliaWithId("localhost:"+strconv.Itoa(begin_port-1),SetBitID(0))
//   //Build the Tree Structure
// 	count := 0
// 	for i := 1; i < 6; i++ {
// 		address := "localhost:" + strconv.Itoa(begin_port+count)
// 		count++
// 		tree_node[i] = NewKademliaWithId(address,SetBitID(i))
// 		host_number, port_number, _ := StringToIpPort(address)
// 		tree_node[0].DoPing(host_number, port_number)
//
// 		for j := 1; j < 6; j++ {
// 			address_v := "localhost:" + strconv.Itoa(begin_port+count)
//
// 			count++
// 			tree_node[5*i+j] = NewKademliaWithId(address_v,SetBitID((5*i+j)))
// 			host_number_v, port_number_v, _ := StringToIpPort(address_v)
// 			tree_node[i].DoPing(host_number_v, port_number_v)
// 			for m := 1; m < 6; m++ {
// 				address_r := "localhost:" + strconv.Itoa(begin_port+count)
// 				count++
// 				tree_node[5*(5*i+j)+m] = NewKademliaWithId(address_r,SetBitID((5*(5*i+j)+m)))
// 				host_number_r, port_number_r, _ := StringToIpPort(address_r)
// 				tree_node[5*i+j].DoPing(host_number_r, port_number_r)
// 			}
// 		}
// 	}
//
// 	//Implement DoIterativeFindNode
// 	key := SetBitID(10)
//   key[19] = key[19] ^ byte(1)
//   value := []byte("Hello world")
// 	contacts, err := tree_node[0].DoIterativeStore(key,value)
//   if contacts == nil {
//     t.Error("Failed DoIterativeStore: No Contacts")
//   }
//   // PrintList(contacts)
//   res_v, err := tree_node[0].DoIterativeFindValue(key)
//   if err != nil ||  string(res_v) != string(value){
//     t.Error("Error DoIterativeStore")
//   }
// 	return
// }

/*

TestIterativeFindValue2:


                       1    1
                     /    /
			        /    /
		          B -- 2 -- 2
		        /   \  . \  .
		       /     \ .  \ .
		     A -- C    5    5
		       \  .
		        \ .
                  F

Create a three-layer tree structure, and each node(server) have five children (156 nodes in total).
Each node that is not leaf node DoPing to its children.
The i-th node have a NodeID of which the i-th bit is set
(e.g. tree_node[10].NodeID = [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 0]),
and its children is tree_node[5i+1:5i+5].
Then we do a DoIterativeStore at tree_node[0] with key=[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 1] and
value = “Hello World”, and perform a DoIterativeFindValue with key =[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 2]).
Result: Value Not Found, Closest node:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 0]. Correct.

*/
func TestIterativeFindValue2(t *testing.T) {

  fmt.Println("TestIterativeFindValue2")
  // time.Sleep(time.Second * 1)

  begin_port := 5482
  tree_node := make([]*Kademlia,160)
	tree_node[0] = NewKademliaWithId("localhost:"+strconv.Itoa(begin_port-1),SetBitID(0))
  //Build the Tree Structure
	count := 0
	for i := 1; i < 6; i++ {
		address := "localhost:" + strconv.Itoa(begin_port+count)
		count++
		tree_node[i] = NewKademliaWithId(address,SetBitID(i))
		host_number, port_number, _ := StringToIpPort(address)
		tree_node[0].DoPing(host_number, port_number)

		for j := 1; j < 6; j++ {
			address_v := "localhost:" + strconv.Itoa(begin_port+count)

			count++
			tree_node[5*i+j] = NewKademliaWithId(address_v,SetBitID((5*i+j)))
			host_number_v, port_number_v, _ := StringToIpPort(address_v)
			tree_node[i].DoPing(host_number_v, port_number_v)
			for m := 1; m < 6; m++ {
				address_r := "localhost:" + strconv.Itoa(begin_port+count)
				count++
				tree_node[5*(5*i+j)+m] = NewKademliaWithId(address_r,SetBitID((5*(5*i+j)+m)))
				host_number_r, port_number_r, _ := StringToIpPort(address_r)
				tree_node[5*i+j].DoPing(host_number_r, port_number_r)
			}
		}
	}

	//Implement DoIterativeFindNode
	key := SetBitID(10)
  key[19] = key[19] ^ byte(1)
  value := []byte("Hello world")
	contacts, err := tree_node[0].DoIterativeStore(key,value)
  if contacts == nil {
    t.Error("Failed DoIterativeStore: No Contacts")
  }
  // PrintList(contacts)
  key[19] <<= 1
  res_v, err := tree_node[0].DoIterativeFindValue(key)
  if !(res_v == nil && err != nil) {
    t.Error("Error DoIterativeStore")
  }
  fmt.Println(err)
	return
}
