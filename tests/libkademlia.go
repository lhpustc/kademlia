package libkademlia

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	// "encoding/binary"
)

const (
	alpha = 3
	b     = 8 * IDBytes
	K     = 20
)

type KBucket []Contact

type kv_pair struct {
	key ID
	value []byte
}

type FindNodeReq struct {
	NodeID ID
	NodesChan chan []Contact
}

type ValueRpl struct {
	value []byte
	ok bool
}
type FindValueReq struct {
	Key ID
	ValueChan chan ValueRpl
}

// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID      ID
	SelfContact Contact
	buckets [b]KBucket
	DHT map[ID][]byte
	Senders chan Contact
	kvPairs chan kv_pair
	UpdateCH     chan bool
	FindNodeChan chan FindNodeReq
	FindValueChan chan FindValueReq
}

func NewKademliaWithId(laddr string, nodeID ID) *Kademlia {
	// fmt.Println("NewKademlia:")
	k := new(Kademlia)
	k.NodeID = nodeID
	k.DHT = make(map[ID][]byte)
	k.Senders = make(chan Contact)
	k.kvPairs = make(chan kv_pair)
	k.UpdateCH = make(chan bool, 100)
	k.FindNodeChan = make(chan FindNodeReq)
	k.FindValueChan = make(chan FindValueReq)

	// TODO: Initialize other state here as you add functionality.

	// Set up RPC server
	// NOTE: KademliaRPC is just a wrapper around Kademlia. This type includes
	// the RPC functions.

	s := rpc.NewServer()
	s.Register(&KademliaRPC{k})
	hostname, port, err := net.SplitHostPort(laddr)
	// fmt.Println("hostname, port:",hostname,port)
	// h, err := net.LookupHost(hostname)
	// fmt.Println("hostname",h)
	if err != nil {
		return nil
	}
	s.HandleHTTP(rpc.DefaultRPCPath+port,
		rpc.DefaultDebugPath+port)
	// s.HandleHTTP(rpc.DefaultRPCPath+"127.0.0.1"+port,
	// 	rpc.DefaultDebugPath+"127.0.0.1"+port)
	// fmt.Println("the lader is ", laddr)
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		fmt.Println("Kademlia server failed")
		log.Fatal("Listen: ", err)
	}

	// Run RPC server forever.
	go http.Serve(l, nil)

    // fmt.Println("Kademlia server succeeed")
	// Add self contact
	hostname, port, _ = net.SplitHostPort(l.Addr().String())
	port_int, _ := strconv.Atoi(port)
	ipAddrStrings, err := net.LookupHost(hostname)
	var host net.IP
	for i := 0; i < len(ipAddrStrings); i++ {
		// fmt.Println(ipAddrStrings[i])
		host = net.ParseIP(ipAddrStrings[i])
		if host.To4() != nil {
			break
		}
	}
	// fmt.Println("!!!-------------")
	k.SelfContact = Contact{k.NodeID, host, uint16(port_int)}
	go k.Update()

	return k
}

func NewKademlia(laddr string) *Kademlia {
	return NewKademliaWithId(laddr, NewRandomID())
}


func inBucket(bucket KBucket, nodeId ID) int {
	for idx, e := range bucket {
		if e.NodeID == nodeId {
			return idx
		}
	}
	return -1
}


// func (k *Kademlia) Update() {
// 	for {
// 		select {
// 		case Sender := <- k.Senders:
// 			distance := k.NodeID.Xor(Sender.NodeID)
// 			i := distance.PrefixLen()
// 			idx := inBucket(k.buckets[i], Sender)
// 			if  idx != -1 {
// 				k.buckets[i] = append(k.buckets[i], k.buckets[i][idx])
// 				k.buckets[i] = append(k.buckets[i][:idx], k.buckets[i][idx+1:]...)
// 			} else {
// 				if len(k.buckets[i]) < K {
// 					k.buckets[i] = append(k.buckets[i], Sender)
// 				} else {
// 					least_recent_contact, err := k.DoPing(k.buckets[i][0].Host, k.buckets[i][0].Port)
// 					if err != nil {
// 						k.buckets[i] = append(k.buckets[i][1:], Sender)
// 					} else {
// 						k.buckets[i] = append(k.buckets[i][1:], *least_recent_contact)
// 					}
// 				}
// 			}
// 		}
// 	}
// }
func (k *Kademlia) UpdateBuckets(Sender Contact) {
	// fmt.Println("Thread: ", k.SelfContact.Port)
	distance := k.NodeID.Xor(Sender.NodeID)
	i := distance.PrefixLen()
	// fmt.Println("Thread: ", k.SelfContact.Port,"-**-",i)
	idx := inBucket(k.buckets[i], Sender.NodeID)
	// fmt.Println("Thread: ", k.SelfContact.Port,Sender.Port)
	if  idx != -1 {

		// fmt.Println("Thread: ", k.SelfContact.Port,"-1-",Sender.Port)
		c := k.buckets[i][idx]
		k.buckets[i] = append(k.buckets[i][:idx], k.buckets[i][idx+1:]...)
		k.buckets[i] = append(k.buckets[i], c)
	} else {
		if len(k.buckets[i]) < K {
			// fmt.Println("Thread: ", k.SelfContact.Port,"-2-",Sender.Port)
			k.buckets[i] = append(k.buckets[i], Sender)
		} else {
			// fmt.Println("Thread: ", k.SelfContact.Port,"-3-",Sender.Port)
			least_recent_contact, err := k.DoPing(k.buckets[i][0].Host, k.buckets[i][0].Port)
			if err != nil {
				k.buckets[i] = append(k.buckets[i][1:], Sender)
			} else {
				k.buckets[i] = append(k.buckets[i][1:], *least_recent_contact)
			}
		}
	}
	// fmt.Println("Thread: ", k.SelfContact.Port,"-END-")
	k.UpdateCH <- true;
}

func (k *Kademlia) Update() {
	for {
		select {
		case Sender := <- k.Senders:
			k.UpdateBuckets(Sender)
		case kv := <-k.kvPairs:
			k.DHT[kv.key] = kv.value
		case nodereq := <- k.FindNodeChan:
			nodereq.NodesChan <- k.FindClosestKNodes(nodereq.NodeID)
		case valuereq := <- k.FindValueChan:
			v,ok := k.DHT[valuereq.Key]
			valuereq.ValueChan <- ValueRpl{v, ok}
		}
	}
}

// func TestConnection(node1 *Kademlia, node2 *Kademlia) {
// 	contact2, err := node1.FindContact(node2.NodeID)
// 	if err != nil {
// 		fmt.Println(strconv.Itoa(int(node1.SelfContact.Port)) + " does not contain " + strconv.Itoa(int(contact2.Port)))
// 		return
// 	} else {
// 		fmt.Println(strconv.Itoa(int(node1.SelfContact.Port)) + " contain " + strconv.Itoa(int(contact2.Port)))
// 	}
// 	contact1, err2 := node2.FindContact(node1.NodeID)
// 	if err2 != nil {
// 		fmt.Println(strconv.Itoa(int(node2.SelfContact.Port)) + " does not contain " + strconv.Itoa(int(contact1.Port)))
// 		return
// 	} else {
// 		fmt.Println(strconv.Itoa(int(node2.SelfContact.Port)) + " contain " + strconv.Itoa(int(contact1.Port)))
// 	}
// 	fmt.Println(strconv.Itoa(int(node1.SelfContact.Port))+ "-" + strconv.Itoa(int(node2.SelfContact.Port)) + " connected")
// 	return
// }



type ContactNotFoundError struct {
	id  ID
	msg string
}

func (e *ContactNotFoundError) Error() string {
	return fmt.Sprintf("%x %s", e.id, e.msg)
}

func (k *Kademlia) FindContact(nodeId ID) (*Contact, error) {
	// TODO: Search through contacts, find specified ID
	// Find contact with provided ID
	if nodeId == k.SelfContact.NodeID {
		// fmt.Println(nodeId, "Found Myself")
		return &k.SelfContact, nil
	}
	distance := k.NodeID.Xor(nodeId)
	i := distance.PrefixLen()
	for _,contact := range k.buckets[i]{
		if contact.NodeID == nodeId{
			// fmt.Println("%v %v\n", contact.Host, contact.Port)
			return &contact, nil
		}
	}
	return nil, &ContactNotFoundError{nodeId, "Not found"}
}

type CommandFailed struct {
	msg string
}

func (e *CommandFailed) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

func min(x, y int) int {
	if x<y {
		return x
	}
	return y
}

func SortKNodes(nodeID ID, nodes []Contact) (ret []Contact){
	var distances []ID
	for _, node := range nodes {
		distances = append(distances, nodeID.Xor(node.NodeID))
	}
	for i:=0; i< min(len(nodes),K); i++ {
		for j:=len(nodes)-1; j>i; j-- {
			if distances[j].Less(distances[j-1]) {
				temp := nodes[j]
				nodes[j] = nodes[j-1]
				nodes[j-1] = temp
			}
		}
	}
	if len(nodes) > K {
		ret = nodes[:K]
	} else {
		ret = nodes
	}
	return
}

func (k *Kademlia) FindClosestKNodes(nodeID ID) (ret []Contact){
  var nodes []Contact
	for _, bucket := range k.buckets {
		if len(bucket) > 0 {
			nodes = append(nodes, bucket...)
		}
	}
	ret = SortKNodes(nodeID, nodes)
	return
}

func (k *Kademlia) DoPing(host net.IP, port uint16) (*Contact, error) {
	// TODO: Implement
	fmt.Println("DoPing:")
	// fmt.Println(rpc.DefaultRPCPath+host.String()+strconv.Itoa(int(port)))

	// fmt.Println("I(" + strconv.Itoa(int(k.SelfContact.Port))+") want to dial: " + host.String() + ":" + strconv.Itoa(int(port)))
    // client, err := rpc.DialHTTP("tcp", host.String() + ":" + strconv.Itoa(int(port)))
    client, err := rpc.DialHTTPPath("tcp", host.String() + ":" + strconv.Itoa(int(port)),
		rpc.DefaultRPCPath+strconv.Itoa(int(port)))
    // client, err := rpc.DialHTTP("tcp",  "localhost:7891")
    if err != nil {
    	fmt.Println(err)
    	fmt.Println("DOPING: Client setup failed!")
        return nil, &CommandFailed{"DOPING: Client setup failed!"}
    }
    // fmt.Println("Dial succeed!")
    ping := new(PingMessage)

    ping.Sender.NodeID = k.NodeID
    ping.Sender.Host = k.SelfContact.Host
    ping.Sender.Port = k.SelfContact.Port
    ping.MsgID = NewRandomID()


    var pong PongMessage

    err = client.Call("KademliaRPC.Ping", ping, &pong)

		defer func() {
			client.Close()
		}()

    if err != nil {
    	fmt.Println(err)
        return nil, &CommandFailed{"DOPING: Client call failed!"}
    }
    // fmt.Println("call succeed!")
    if pong.MsgID.Equals(ping.MsgID){
    	// fmt.Println("Preparing to update")
		k.Senders <- pong.Sender
		_ = <- k.UpdateCH
		return nil, &CommandFailed{fmt.Sprintf("Succeed!")}
    }
	return nil, &CommandFailed{
		"Unable to ping " + fmt.Sprintf("%s:%v", host.String(), port)}
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) error {
	// TODO: Implement
	client, err := rpc.DialHTTPPath("tcp", contact.Host.String() + ":" + strconv.Itoa(int(contact.Port)),rpc.DefaultRPCPath+strconv.Itoa(int(contact.Port)))

  if err != nil{
   log.Fatal("DialHTTP: ", err)
  }

  storeRequest := new (StoreRequest)
  storeResult := new (StoreResult)
  storeRequest.MsgID = NewRandomID()
  storeRequest.Sender = k.SelfContact
  storeRequest.Key = key
  storeRequest.Value = value

  err = client.Call("KademliaRPC.Store", storeRequest, &storeResult)
  if err != nil{
   log.Fatal("Call: ", err)
   return &CommandFailed{"Unable to store {"+key.AsString()+"}"}
  }


  if(storeResult.MsgID.Compare(storeRequest.MsgID) == 0){
  	// k.Sender <- storeResult
   	return nil
  }
	return &CommandFailed{"Not implemented"}
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) ([]Contact, error) {
	// TODO: Implement
	hostname := contact.Host.String()
	port_str := strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTPPath("tcp", hostname+":"+port_str,rpc.DefaultRPCPath+port_str)
	if err != nil {
		fmt.Println(err)
		return nil, &CommandFailed{"DoFindNode failed!"}
	}

	req := FindNodeRequest{k.SelfContact, NewRandomID(), searchKey}
	res := new(FindNodeResult)
	err = client.Call("KademliaRPC.FindNode", req, &res)
	defer func() {
		client.Close()
	}()
	if err != nil {
		return nil, &CommandFailed{"RPC Error"}
	}
	k.Senders <- *contact
	if res.Err != nil {
		return nil, &CommandFailed{"FindNodeError"}
	}
	return res.Nodes, nil
}

func (k *Kademlia) DoFindValue(contact *Contact,
	searchKey ID) (value []byte, contacts []Contact, err error) {
	// TODO: Implement
	client, err := rpc.DialHTTPPath("tcp", contact.Host.String() + ":" + strconv.Itoa(int(contact.Port)),
	rpc.DefaultRPCPath+strconv.Itoa(int(contact.Port)))
	if err != nil {
		fmt.Println(err)
		fmt.Println("DOFINDVALUE: Client setup failed!")
			// return nil, &CommandFailed{"DOFINDVALUE: Client setup failed!"}
	}
	fmt.Println("Dial succeed!")
	fvr := new(FindValueRequest)

	fvr.Sender.NodeID = k.NodeID
	fvr.Sender.Host = k.SelfContact.Host
	fvr.Sender.Port = k.SelfContact.Port
	fvr.MsgID = NewRandomID()
	fvr.Key = searchKey

	var fvrl FindValueResult

	err = client.Call("KademliaRPC.FindValue", fvr, &fvrl)
	if err != nil {
		fmt.Println(err)
			return nil, nil, &CommandFailed{"DOFINDVALUE: Client call failed!"}
	}

	fmt.Println("call succeed!")
	if fvrl.MsgID.Equals(fvr.MsgID){
		fmt.Println("Preparing to return")
	  return fvrl.Value, fvrl.Nodes, nil
	}
	return nil, nil, &CommandFailed{"Not implemented"}
}

func (k *Kademlia) LocalFindValue(searchKey ID) ([]byte, error) {
	// TODO: Implement
	// value, ok := k.DHT[searchKey]
	// if ok {
	// 	return value, nil
	// }
	vrplChan := make(chan ValueRpl)
	k.FindValueChan <- FindValueReq{searchKey, vrplChan}
	vrpl := <- vrplChan
	if vrpl.ok {
		return vrpl.value, nil
	}

	return []byte(""), &CommandFailed{"Not implemented"}
}

// For project 2!
func (k *Kademlia) DoIterativeFindNode(id ID) ([]Contact, error) {
	return nil, &CommandFailed{"Not implemented"}
}
func (k *Kademlia) DoIterativeStore(key ID, value []byte) ([]Contact, error) {
	return nil, &CommandFailed{"Not implemented"}
}
func (k *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
	return nil, &CommandFailed{"Not implemented"}
}

// For project 3!
func (k *Kademlia) Vanish(data []byte, numberKeys byte,
	threshold byte, timeoutSeconds int) (vdo VanashingDataObject) {
	return
}

func (k *Kademlia) Unvanish(searchKey ID) (data []byte) {
	return nil
}
