package libkademlia

import (
	//"bytes"
	"fmt"
	"net"
	"strconv"
	"testing"
	//"container/heap"
	"time"
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

// // EXTRACREDIT
// // Check out the correctness of DoIterativeFindNode function
// func TestIterativeFindNode1(t *testing.T) {
// 	// tree structure;
// 	// A<--->B
// tree_node := make([]*Kademlia, 30)
// address := make([]string, 30)
// for i := 0; i < 30; i++ {
// 	address[i] = "localhost:" + strconv.Itoa(7000+i)
// 	tree_node[i] = NewKademlia(address[i])
// }
//
// //30 nodes ping each other
// for i := 0; i < 30 ; i++ {
// 	for j := 0; j < 30; j++ {
// 		host_number, port_number, _ := StringToIpPort(address[j])
// 		tree_node[i].DoPing(host_number, port_number)
// 	}
// }
//
// //find node[19], start from node 0
// contacts, _ := tree_node[0].DoIterativeFindNode(tree_node[19].SelfContact.NodeID)
// // contacts, _ := tree_node[0].DoFindNode(&tree_node[1].SelfContact,tree_node[19].SelfContact.NodeID)
// count := 0
// //check the result
// for i := 0; i < len(contacts); i++ {
// 	if(contacts[i].NodeID.Equals(tree_node[19].SelfContact.NodeID)) {
// 		count ++
// 	}
// 	//fmt.Print(contacts[i].NodeID)
// }
// if(count != 1) {
// 	t.Error("the result is not true, count="+strconv.Itoa(count) + "; len = " +strconv.Itoa(len(contacts)))
// }
// return
// }

// EXTRACREDIT
// Check out the correctness of DoIterativeFindNode function
// func TestIterativeFindNode2(t *testing.T) {
// 	// tree structure;
// 	// A<--->B
// 	total_num := 200											//total num mush >= 20
// 	tree_node := make([]*Kademlia, total_num)
// 	instance1 := NewKademlia("localhost:7399")										//starting node
// 	host1, port1, _ := StringToIpPort("localhost:7399")
//
// 	findId := total_num - 20
// 	//initialize the nodes
// 	for i := 0; i < total_num; i++ {
// 		address := "localhost:" + strconv.Itoa(7400+i)
// 		tree_node[i] = NewKademlia(address)
// 		tree_node[i].DoPing(host1, port1)														//every node ping instance1
// 	}
// 	target := tree_node[findId].SelfContact.NodeID								//id to be found
//
// 	//node less that findId ping the findId, not all ping the findId node
// 	for i := 0; i < findId; i++ {
//     fmt.Println("DoPing:"+strconv.Itoa(findId)+" to "+strconv.Itoa(i))
// 		tree_node[findId].DoPing(tree_node[i].SelfContact.Host, tree_node[i].SelfContact.Port)
//     fmt.Println("ok")
// 	}
//
// 	result, err := instance1.DoIterativeFindNode(target)						//start from instance1,find target
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
//
// 	if result == nil || len(result) == 0 {
// 		t.Error("No contacts were found")
// 	}
//
// 	//check the result
// 	count := 0
// 	for _, value := range result {
// 		if value.NodeID.Equals(target) {
// 			count++
// 		}
// 	}
// 	//if only Contains one same node, test pass
// 	if count != 1 {
// 		t.Error(count)
// 	}
// 	return
// }

// func TestGenerateID(t *testing.T) {
//   ret := SetBitID(9)
//   fmt.Println(ret)
//   // fmt.Println(id)
//   // if !ret.Equals(id) {
//   //   t.Error("TestGenerateID Failed")
//   // }
// }

// // EXTRACREDIT
// // Check out the correctness of DoIterativeFindNode function
// func TestIterativeFindNode(t *testing.T) {
// 	// tree structure;
// 	//  A--> B --> 5 nodes --> 5 nodes --> 5nodes
// 	/*
// 		                F
// 			             /
// 		          C --G    J
// 		         /    \  /
// 		       /        H --I
// 		   A-B -- D      \
// 		       \          K
// 		        E
// */
//   fmt.Println("TestIterativeFindNode")
//   begin_port := 6482
// 	// instance1 := NewKademlia("localhost:4480")
// 	instance2 := NewKademliaWithId("localhost:"+strconv.Itoa(begin_port-1),SetBitID(0))
// 	// host2, port2, _ := StringToIpPort("localhost:4481")
// 	// instance1.DoPing(host2, port2)
// 	// _, err := instance1.FindContact(instance2.NodeID)
// 	// if err != nil {
// 	// 	t.Error("Instance 2's contact not found in Instance 1's contact list")
// 	// 	return
// 	// }
//   //Build the Tree Structure
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

// // check correctness of
func TestIterativeStore(t *testing.T) {
	// tree structure;
	//  A--> B --> 5 nodes --> 5 nodes --> 5nodes
	/*
		                F
			             /
		          C --G    J
		         /    \  /
		       /        H --I
		   A-B -- D      \
		       \          K
		        E
*/
  fmt.Println("TestIterativeStore")
  time.Sleep(time.Second*1)
  begin_port := 3482
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
  res_v, err := tree_node[10].LocalFindValue(key)
  if err != nil ||  string(res_v) != string(value){
    t.Error("Error DoIterativeStore")
  }
	return
}

// // check correctness of DoIterativeFindValue
func TestIterativeFindValue(t *testing.T) {
	// tree structure;
	//  A--> B --> 5 nodes --> 5 nodes --> 5nodes
	/*
		                F
			             /
		          C --G    J
		         /    \  /
		       /        H --I
		   A-B -- D      \
		       \          K
		        E
*/
  fmt.Println("TestIterativeFindValue")
  time.Sleep(time.Second * 1)
  begin_port := 4482
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
  res_v, err := tree_node[0].DoIterativeFindValue(key)
  if err != nil ||  string(res_v) != string(value){
    t.Error("Error DoIterativeStore")
  }
	return
}

// // check correctness of DoIterativeFindValue
func TestIterativeFindValue2(t *testing.T) {
	// tree structure;
	//  A--> B --> 5 nodes --> 5 nodes --> 5nodes
	/*
		                F
			             /
		          C --G    J
		         /    \  /
		       /        H --I
		   A-B -- D      \
		       \          K
		        E
*/
  fmt.Println("TestIterativeFindValue2")
  time.Sleep(time.Second * 1)

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
