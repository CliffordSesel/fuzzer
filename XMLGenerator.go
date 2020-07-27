package main

import (
	"os"
	"aqwari.net/xml/xmltree"
	"fmt"
	"log"
	"io/ioutil"
	"math/rand"
	"time"
)

/*
 * elemSet is a slice of pointers to all unique elements in the tree, 
 * allowing an element to be randomly selected in O(1) time. 
 * XMLTree is the root element.  
 */
type mXMLHolder struct {
	XMLTree	*xmltree.Element
	elemSet []*xmltree.Element
	description []string
}

func createXMLHolder(description string) mXMLHolder {
	s := mXMLHolder{}
	s.description = append(s.description, description)
	return s
}

/*
 * Reads the XML specified by file into the XMLHolder. 
 */
func (s *mXMLHolder) read(file string) {
	xmlFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	xmlBytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
	}

	s.XMLTree, err = xmltree.Parse(xmlBytes)
	if err != nil {
		log.Fatal(err)
	}

	operation := fmt.Sprintf("Read in XML from: %s\n", file)
	s.description = append(s.description, operation)

}

/*
 * Returns a slice of all Elements in the tree t
 */
func getElemSet(t *xmltree.Element) []*xmltree.Element {
	elems := t.Flatten()

	// add root of the tree to elems since the Flatten method
	// doesn't do it. 
	elems = append(elems, t)
	return elems
}

/*
 * randomly selects an element from the pool of elements in
 * a mXMLHolder
 */
func selectElement(s *mXMLHolder) *xmltree.Element {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	i := r.Intn(len(s.elemSet))
	return s.elemSet[i]
}

/*
 * Creates a new Element which contains the same StartElement and
 * and Content as the element e. The new Element has no children. 
 */
func childlessClone(e *xmltree.Element) *xmltree.Element {
	clone := new(xmltree.Element)
	clone.StartElement = e.StartElement.Copy()
	content := make([]byte, len(e.Content))
	copy(content, e.Content)
	clone.Content = content
	return clone
}

/*
 * Randomly selects two elements from the element pool. 
 * One element is chosen to be the parent, the other element is
 * chosen to be the child. Multiple childless copies of the child 
 * are added to the parent. 
 */
func (s *mXMLHolder) spamElementBreadthWise() {
	parent := selectElement(s)
	child := selectElement(s)

	for i := 0; i < 10; i++ {
		clone := childlessClone(child)
		parent.Children = append(parent.Children, *clone)
	}

}
