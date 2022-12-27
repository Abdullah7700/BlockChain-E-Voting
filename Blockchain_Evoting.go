package main

import (
	"crypto/sha256"
	"fmt"
	"math"
	
	"context"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"os"
	"time"
)

type treenode struct {
	left  *treenode
	right *treenode
	data  string
}

type blockdata struct {
	merkleroot string
	phash      string
	nonce      int
}

type Blockchain struct {
	prev         *Blockchain
	data         blockdata
	chash        string
	root         *treenode
	difficulty   int
	transactions string
}

// this function gets nil chainhead
func genesisblock(chainhead *Blockchain) *Blockchain {
	var x *Blockchain
	x = new(Blockchain)
	x.prev = chainhead
	x.data.phash = "0"
	x.data.nonce = 0
	x.data.merkleroot = "0"
	x.chash = CalculateHash(x)
	return x
}

func NewBlock(transactions []string, chainhead *Blockchain, diff int) *Blockchain {
	var x *Blockchain
	x = new(Blockchain)
	x.prev = chainhead
	//because 1st block is empty it will not caclulate hash of previous block
	if chainhead != nil {
		x.data.phash = CalculateHash(chainhead)
	}
	x.difficulty = diff
	x.root, x.data.merkleroot = createMerkleTree(transactions)
	x.data.nonce = MineBlock(x)
	x.chash = CalculateHash(x)
	
	//Connecting to bootstrap node/network part
	//connect_to_Bootstrap()	
	
	return x
}

func MineBlock(bl *Blockchain) int {
	count := 0
	for i := 1; i < math.MaxInt64; i++ {
		bl.data.nonce = i
		wew := CalculateHash(bl)
		for i := 0; i < bl.difficulty; i++ {
			if string(wew[i]) == "0" {
				count++
			}
		}
		if count == bl.difficulty {
			return i
		} else {
			count = 0
		}
	}
	return 0
}

func createMerkleTree(fortree []string) (*treenode, string) {
	queue := make([]*treenode, 0)
	for _, transaction := range fortree {
		var leaf *treenode
		leaf = new(treenode)
		leaf.data = transaction
		queue = append(queue, leaf)
	}

	if len(queue)%2 != 0 { //if odd element, copy last element again to make tree balanced
		var leaf *treenode
		leaf = new(treenode)
		leaf.data = queue[len(queue)-1].data
		queue = append(queue, leaf)
	}

	for len(queue) > 1 {
		var parent *treenode
		parent = new(treenode)
		parent.left, queue = queue[0], queue[1:]
		parent.right, queue = queue[0], queue[1:]

		tohash := make([]string, 0)
		tohash = append(tohash, parent.left.data)
		tohash = append(tohash, parent.right.data)
		parent.data = CalculateHashForString(tohash)

		queue = append(queue, parent)
	}

	return queue[0], queue[0].data
}

func CalculateHash(inpblock *Blockchain) string {

	newhash := sha256.New()
	newhash.Write([]byte(fmt.Sprintf("%v", inpblock.data)))
	return fmt.Sprintf("%x", newhash.Sum(nil))
}

func CalculateHashForString(inp []string) string {

	newhash := sha256.New()
	newhash.Write([]byte(fmt.Sprintf("%v", inp)))
	return fmt.Sprintf("%x", newhash.Sum(nil))
}

func DisplayBlocks(inpblock *Blockchain) {
	h := inpblock
	for h != nil {
		fmt.Println("-------BLOCK------------------------------------------------------")
		fmt.Println("nonce =", h.data.nonce)
		fmt.Println("previous block hash=", h.data.phash)
		fmt.Println("current hash=", h.chash)
		fmt.Println("-------DISPLAYING MERKLE TREE-------------------------------------")
		DisplayMerkle(h.root)
		fmt.Println("-------END OF TREE------------------------------------------------")
		fmt.Println("------------------------------------------------------------------")
		h = h.prev
	}
}

func DisplayMerkle(node *treenode) {
	if node == nil {
		return
	}
	DisplayMerkle(node.left)
	DisplayMerkle(node.right)
	if node.left == nil {
		fmt.Println("transaction =", node.data)
	} else {
		fmt.Println("hash =", node.data)
	}
}

func VerifyChain(inpblock *Blockchain) {
	var x *Blockchain
	x = inpblock
	for x != nil {
		if x.prev != nil {
			if x.data.phash != x.prev.chash {
				fmt.Println("Blockchain compromised")
				return
			}
		}
    x = x.prev
	}
	fmt.Println("Blockchain integrity safe")
	return
}

func ChangeBlock(ref *Blockchain, changed []string) *Blockchain {
	ref.root, _ = createMerkleTree(changed)
	return ref
}

func main() {

	var check *Blockchain // this is equal to nil rn
	realtransactions := []string{"Alice to Asim 1 vote", "Bob to Asim 1 vote", "John to Bilal 1 vote", "Mark to Asim 1 vote"}
	realtransactions2:=[]string{ "M to B 1 vote","A to B 1 vote", "B to A 1 vote", "R to A 1 vote"}

	check = GenesisBlock(check)

	check = NewBlock(realtransactions, check, 3)
  	check = NewBlock(realtransactions,check, 2)
	VerifyChain(check)
  	ChangeBlock(check, realtransactions2)

	DisplayBlocks(check)
}
