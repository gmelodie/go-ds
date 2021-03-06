package trie

import (
	"errors"
	"fmt"
	"strings"
)

// Data is the data type the trie holds
type Data int

// Trie is the root of the tree.
// It doesn't have any data or prefixes.
type Trie struct {
	root *Node
}

// Node is a node of the Trie.
type Node struct {
	data       Data
	isTerminal bool
	prefix     rune
	children   []*Node
}

// CreateTrie creates an empty Trie tree and returns it.
func CreateTrie() *Trie {
	return &Trie{&Node{}}
}

// CreateNode creates a new node.
// Returns a pointer to that node.
func CreateNode(data Data, isTerminal bool, prefix rune) *Node {
	return &Node{data, isTerminal, prefix, nil}
}

// Insert appends a node (n) containing data and prefix to the trie.
// Returns an error if the node already exists.
func (t *Trie) Insert(data Data, prefix []rune) error {

	if len(prefix) == 0 {
		return errors.New("Can't insert node with empty prefix")
	}

	if t == nil {
		return errors.New("Can't insert on nil trie")
	}

	n := t.root
	lastChar := 0
	// Move to the last existing node
	for i, c := range prefix {
		aux := n.hasChildWithPrefix(c)
		if aux == nil { // doesn't have child, insert it
			lastChar = i
			break // Stops when child with prefix char doesn't exist
		}
		n = aux // has child, go down in tree
	}

	remainingPrefix := []rune(prefix)[lastChar:]
	return n.createSubTree(data, remainingPrefix)
}

// createSubTree inserts the rest of a prefix beginning
// in the Node n.
func (n *Node) createSubTree(data Data, prefix []rune) error {
	var newNode *Node
	var newNodeDad = n

	// Node already exists, make terminal
	if len(prefix) == 1 && n.prefix == prefix[0] {
		n.isTerminal = true
		n.data = data
		return nil
	}

	for _, c := range prefix {
		newNode = CreateNode(-1, false, c)
		(*newNodeDad).children = append((*newNodeDad).children, newNode)
		newNodeDad = newNode
	}

	// Insert data in last node
	newNode.isTerminal = true
	newNode.data = data

	return nil
}

// hasChildWithPrefix returns a *Node containing
// the child of n that has a prefix of c.
// Returns false otherwise.
func (n *Node) hasChildWithPrefix(c rune) *Node {

	if n.children == nil {
		return nil
	}

	for _, child := range n.children {
		if child.prefix == c {
			return child
		}
	}

	return nil
}

// Delete searches for a prefix in the Trie.
// Removes the node and rearranges the tree if prefix exists.
// Returns an error if prefix doesn't exist.
func (t *Trie) Delete(prefix []rune) error {

	if t == nil {
		return errors.New("Can't delete in nil trie")
	}

	if prefix == nil {
		return errors.New("Can't delete nil prefix")
	}

	lookup := t.root
	lastUseful := t.root
	deleteIDX := -1

	// Get last useful node (which we can't delete)
	for i, c := range prefix {
		// We can't delete lookup if it has more than 1 child
		// or if it is terminal for another prefix
		if len(lookup.children) > 1 ||
			(lookup.isTerminal && i != len(prefix)-1) {
			lastUseful = lookup
		}
		// Go down in tree
		for j, n := range lookup.children {
			if n.prefix == c {
				if lastUseful == lookup {
					deleteIDX = j
				}
				lookup = n
				break
			}
			if j == len(lookup.children)-1 {
				return errors.New("Didn't find prefix")
			}
		}
	}

	if len(lookup.children) > 0 {
		lookup.isTerminal = false
		lookup.data = -1
	} else {
		// Remove subtree below lastUseful at index deleteIDX
		lastUseful.children[deleteIDX] =
			lastUseful.children[len(lastUseful.children)-1]

		(*lastUseful).children =
			(*lastUseful).children[:len((*lastUseful).children)-1]
	}

	return nil
}

// Update searches for a prefix in the Trie.
// Updates the node if prefix exists.
// Returns an error if prefix doesn't exist.
func (t *Trie) Update(prefix []rune, data Data) error {
	n, err := t.searchNode(prefix)

	if n != nil {
		n.data = data
		return nil
	}

	return err
}

// Search looks for the node indexed by prefix.
// Returns a string containing the data if prefix exists.
// Returns an empty string and error if prefix doesn't exist.
func (t *Trie) Search(prefix []rune) (Data, error) {
	if t == nil {
		return 0, errors.New("Can't search in nil trie")
	}

	if prefix == nil {
		return 0, errors.New("Can't search nil prefix")
	}

	n, err := t.searchNode(prefix)

	if n != nil {
		return n.data, nil
	}

	return 0, err
}

// searchNode returns the node containing the data for prefix.
// Returns an error if the prefix doesn't exist in the tree.
func (t *Trie) searchNode(prefix []rune) (*Node, error) {
	// Lookup node starts at root
	lookup := t.root

	for _, c := range prefix {
		for i, n := range lookup.children {
			if n.prefix == c { // found prefix, update lookup (go down)
				lookup = n
				break
			}
			if i == len(lookup.children)-1 { // at the last child of lookup
				return nil, errors.New("Didn't find prefix")
			}
		}
	}

	if lookup.isTerminal {
		return lookup, nil
	}

	return nil, errors.New("Didn't find prefix")
}

// PrintTrie prints trie showing parent-child relationships
func (t *Trie) PrintTrie() error {
	if t.root.children == nil {
		return nil
	}

	for _, n := range t.root.children {
		fmt.Printf("%s: %d\n", string(n.prefix), int(n.data))
		n.printSubTree(1)
	}

	return nil
}

func (n *Node) printSubTree(tabs int) error {
	if n.children == nil {
		return nil
	}

	for _, aux := range n.children {
		fmt.Printf(strings.Repeat(" ", tabs))
		fmt.Printf("%s: %d\n", string(aux.prefix), int(aux.data))
		aux.printSubTree(tabs + 1)
	}

	return nil
}
