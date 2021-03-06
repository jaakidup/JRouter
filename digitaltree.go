package jaakitrouter

// Node ...
type Node struct {
	Path      string
	Name      string
	Value     string
	End       bool
	Delete    bool
	Handler   Handle
	Params    Params
	MaxParams int
	Parent    *Node
	Child     map[string]*Node
}

func (n *Node) hasChildren() bool {
	if len(n.Child) == 0 {
		return false
	}
	return true
}

// NewNode returns reference to new Node
func NewNode() *Node {
	return &Node{
		Child: make(map[string]*Node),
	}
}

// DigitalTree ...
type DigitalTree struct {
	Root *Node
}

// NewDigitalTree returns a ref to a new DigitalTree
func NewDigitalTree() *DigitalTree {
	return &DigitalTree{
		Root: &Node{
			Name:  "ROOT",
			Child: make(map[string]*Node),
		},
	}
}

// Add a word to the DigitalTree
func (dt *DigitalTree) Add(word string, handler Handle, params Params) {
	node := dt.Root
	var path string

	for _, letter := range word {
		char := string(letter)
		_, found := node.Child[char]
		if found {
			path += char
			node = node.Child[char]
		} else {
			newNode := NewNode()
			path += char
			newNode.Path = path
			newNode.Parent = node
			node.Child[char] = newNode
			node = node.Child[char]
		}
	}
	node.Handler = handler
	node.Params = params

	if len(params) >= 2 {
		node.Params = params
		// node.MaxParams = len(params) -1 // remove the function call
		node.MaxParams = len(params)
	} else {
		node.Params = nil
		node.MaxParams = 0
	}

	node.End = true
}

// Find by key
// @word string, key for lookup
// return: truthy if found, falsy if not found
// return handler if found
func (dt *DigitalTree) Find(word string) (bool, Handle, Params, int) {
	node := dt.Root

	for _, letter := range word {
		char := string(letter)
		_, found := node.Child[char]
		if found {
			node = node.Child[char]
		} else {
			return false, nil, nil, 0
		}
	}
	if !node.End {
		return false, nil, nil, 0
		// return false, nil
	}
	return true, node.Handler, node.Params, node.MaxParams
}

// Return the last node of the word
func (dt *DigitalTree) lastNodeOf(word string) *Node {
	node := dt.Root
	for _, letter := range word {
		node = node.Child[string(letter)]
	}
	return node
}

// Delete a word and handler
func (dt *DigitalTree) Delete(word string) {
	lastNode := dt.lastNodeOf(word)
	deleter(lastNode, word, true)
}

// return lastLetter
func lastLetter(word string) (bool, string) {
	wordLength := len(word)
	if wordLength >= 1 {
		return true, word[wordLength-1:]
	}
	return false, ""
}
func allButLastLetter(word string) (bool, string) {
	wordLength := len(word)
	if wordLength > 1 {
		return true, word[:wordLength-1]
	}
	return false, ""
}

// deleter
func deleter(node *Node, word string, first bool) {
	if !first && node.End {
		goto DONE
	}

	if node.hasChildren() {
		node.End = false
		node.Handler = nil
	}
	if !node.hasChildren() {
		node.End = false
		node.Handler = nil

		node = node.Parent

		found, char := lastLetter(word)
		if found {
			delete(node.Child, char)
			ok, nextWord := allButLastLetter(word)
			if ok {
				deleter(node, nextWord, false)
			}
		}
	}
DONE:
}

// ListKeys ...
func (dt *DigitalTree) ListKeys(name string) *ResultSet {

	resultSet := NewResultSet(name)
	Walk("", dt.Root, resultSet)

	resultSet.Count = len(resultSet.Results)
	return resultSet
}

// ResultSet ...
type ResultSet struct {
	Name    string        `json:"name,omitempty"`
	Count   int           `json:"count,omitempty"`
	Results []interface{} `json:"results,omitempty"`
}

// NewResultSet ...
func NewResultSet(name string) *ResultSet {
	return &ResultSet{
		Name:    name,
		Results: []interface{}{},
	}
}

// Walk ...
func Walk(word string, node *Node, results *ResultSet) {
	for char, child := range node.Child {
		if child.End {
			fullWord := word + char
			results.Results = append(results.Results, fullWord)
			if child.hasChildren() {
				Walk(word+char, child, results)
			}
		} else {
			Walk(word+char, child, results)
		}
	}
}
