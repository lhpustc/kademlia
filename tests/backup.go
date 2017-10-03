func (k *Kademlia) RPC_Search(contact Contact, searchKey ID, results chan rpcResult) {
	// TODO: Implement
	hostname := contact.Host.String()
	port_str := strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTPPath("tcp", hostname+":"+port_str,rpc.DefaultRPCPath+contact.Host.String()+port_str)
	if err != nil {
		fmt.Println(err)
		results <- rpcResult{contact, nil, &CommandFailed{"DialHTTPPath failed!"}}
		return
	}

	req := FindNodeRequest{k.SelfContact, NewRandomID(), searchKey}
	res := new(FindNodeResult)

	callChan := make(chan error,1)
	go func() {callChan <- client.Call("KademliaRPC.FindNode", req, &res)}()
	defer func() {
		client.Close()
	}()

  select {
	case err := <- callChan:
		if err != nil {
			results <- rpcResult{contact, nil, &CommandFailed{"RPC Error"}}
			return
		}
		k.Senders <- contact
		if res.Err != nil {
			results <- rpcResult{contact, nil, &CommandFailed{"FindNodeError"}}
			return
		}
		// fmt.Println("RPC_Search")
		// fmt.Println(contact.NodeID)
		for i:=0;i<len(res.Nodes);i++ {
			k.Senders <- res.Nodes[i]
		}
		results <- rpcResult{contact, res.Nodes, nil}
		return
	case <-time.After(time.Second * 1):
		results <- rpcResult{contact, nil, &CommandFailed{"Timeout"}}
		return
	}
}

func (k *Kademlia) RPC_SearchValue(contact *Contact, searchKey ID, results chan rpcResultValue) {
	// TODO: Implement
	hostname := contact.Host.String()
	port_str := strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTPPath("tcp", hostname+":"+port_str,rpc.DefaultRPCPath+contact.Host.String()+port_str)
	if err != nil {
		fmt.Println(err)
		results <- rpcResultValue{contact, nil, nil, &CommandFailed{"DialHTTPPath failed!"}}
		return
	}

	req := FindValueRequest{k.SelfContact, NewRandomID(), searchKey}
	res := new(FindValueResult)


	callChan := make(chan error,1)
	go func() {callChan <- client.Call("KademliaRPC.FindValue", req, &res)}()
	defer func() {
		client.Close()
	}()

  select {
	case err := <- callChan:
		if err != nil {
			results <- rpcResultValue{contact, nil, nil, &CommandFailed{"RPC Error"}}
			return
		}
		k.Senders <- *contact
		if res.Err != nil {
			results <- rpcResultValue{contact, nil, nil, &CommandFailed{"FindNodeError"}}
			return
		}
		results <- rpcResultValue{contact, res.Value, res.Nodes, nil}
		return
	case <-time.After(time.Second * 1):
		results <- rpcResultValue{contact, nil, nil, &CommandFailed{"Timeout"}}
		return
	}
}

func (k *Kademlia) DoIterativeFindNode(id ID) ([]Contact, error) {
	var shortList []Contact
	var closestNode Contact
	stopIter := false
	nodesChan := make(chan []Contact)
	k.FindNodeChan <- FindNodeReq{id, nodesChan}
	tempList := <-nodesChan
	results := make(chan rpcResult, alpha)

  // fmt.Println("tempList len=" + strconv.Itoa(len(tempList)))
	if len(tempList)>0 {
		closestNode = tempList[0]
	}
	for len(shortList)<20 && !stopIter {
		stopIter = true
		minLen := Min(alpha, len(tempList))
		// results := make(chan rpcResult, minLen)
		for i:=0;i<minLen;i++ {
			// fmt.Println("Out send RPC_Search")
			go k.RPC_Search(tempList[i], id, results)
		}
		beforeList := tempList
		for i:=0;i<minLen;i++ {
			result := <- results
			if(result.err == nil) {
				shortList = append(shortList, result.contact)
				clearList(shortList)
				// remove result.contact from tempList
				for j:=0;j<len(tempList);j++ {
					// fmt.Println("minLen="+strconv.Itoa(minLen)+"; tempList len=" + strconv.Itoa(len(tempList)))
					if tempList[j].NodeID.Equals(result.contact.NodeID) {
						tempList = append(tempList[:j], tempList[j+1:]...)
						break
					}
				}
				ret := removeDuplicates(result.Nodes, shortList)
				tempList = mergeSort(tempList, ret, id)
					if len(tempList)>0 && tempList[0].NodeID.Xor(id).Less(closestNode.NodeID.Xor(id)) {
						closestNode = tempList[0]
					}
			} else { // timeout error
				fmt.Println("timeout")
				for j:=0;j<minLen;j++ {
					if tempList[j].NodeID.Equals(result.contact.NodeID) {
						tempList = append(tempList[:j], tempList[j+1:]...)
					}
				}
			}
		}
		stopIter = !differ(beforeList, tempList)
	}
	return SortKNodes(id, shortList), nil
}
