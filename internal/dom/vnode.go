// Package dom provides virtual DOM structures and operations for HTML parsing and manipulation.
package dom

// VNodeType represents the type of a virtual DOM node.
type VNodeType string

const (
	// ElementNode represents an HTML element node.
	ElementNode VNodeType = "element"
	// TextNode represents a text node.
	TextNode VNodeType = "text"
)

// ReadabilityData stores readability-specific information for a node.
type ReadabilityData struct {
	ContentScore float64
}

// VNode is the interface for all virtual DOM nodes.
type VNode interface {
	// Type returns the node type.
	Type() VNodeType
	// Parent returns the parent element of this node, or nil if it has no parent.
	Parent() *VElement
	// SetParent sets the parent element of this node.
	SetParent(parent *VElement)
	// GetReadabilityData returns the readability-specific data for this node.
	GetReadabilityData() *ReadabilityData
	// SetReadabilityData sets the readability-specific data for this node.
	SetReadabilityData(data *ReadabilityData)
}

// baseNode implements common functionality for all node types.
type baseNode struct {
	nodeType       VNodeType
	parent         *VElement
	readabilityData *ReadabilityData
}

// Type returns the node type.
func (n *baseNode) Type() VNodeType {
	return n.nodeType
}

// Parent returns the parent element of this node.
func (n *baseNode) Parent() *VElement {
	return n.parent
}

// SetParent sets the parent element of this node.
func (n *baseNode) SetParent(parent *VElement) {
	n.parent = parent
}

// GetReadabilityData returns the readability-specific data for this node.
func (n *baseNode) GetReadabilityData() *ReadabilityData {
	return n.readabilityData
}

// SetReadabilityData sets the readability-specific data for this node.
func (n *baseNode) SetReadabilityData(data *ReadabilityData) {
	n.readabilityData = data
}

// VText represents a text node in the virtual DOM.
type VText struct {
	baseNode
	TextContent string
}

// NewVText creates a new text node with the given text content.
func NewVText(textContent string) *VText {
	return &VText{
		baseNode: baseNode{
			nodeType: TextNode,
		},
		TextContent: textContent,
	}
}

// VElement represents an element node in the virtual DOM.
type VElement struct {
	baseNode
	TagName    string
	Attributes map[string]string
	Children   []VNode
}

// NewVElement creates a new element node with the given tag name.
func NewVElement(tagName string) *VElement {
	return &VElement{
		baseNode: baseNode{
			nodeType: ElementNode,
		},
		TagName:    tagName,
		Attributes: make(map[string]string),
		Children:   make([]VNode, 0),
	}
}

// ID returns the id attribute of this element, or an empty string if it has no id.
func (e *VElement) ID() string {
	return e.Attributes["id"]
}

// ClassName returns the class attribute of this element, or an empty string if it has no class.
func (e *VElement) ClassName() string {
	return e.Attributes["class"]
}

// AppendChild adds a child node to this element.
func (e *VElement) AppendChild(child VNode) {
	child.SetParent(e)
	e.Children = append(e.Children, child)
}

// SetAttribute sets an attribute on this element.
func (e *VElement) SetAttribute(name, value string) {
	e.Attributes[name] = value
}

// GetAttribute gets the value of an attribute on this element.
func (e *VElement) GetAttribute(name string) string {
	return e.Attributes[name]
}

// HasAttribute checks if this element has the specified attribute.
func (e *VElement) HasAttribute(name string) bool {
	_, ok := e.Attributes[name]
	return ok
}

// VDocument represents a virtual DOM document.
type VDocument struct {
	DocumentElement *VElement
	Body            *VElement
	BaseURI         string
	DocumentURI     string
}

// NewVDocument creates a new virtual DOM document with the given document element and body.
func NewVDocument(documentElement, body *VElement) *VDocument {
	return &VDocument{
		DocumentElement: documentElement,
		Body:            body,
	}
}

// IsVElement checks if a node is a VElement.
func IsVElement(node VNode) bool {
	return node != nil && node.Type() == ElementNode
}

// AsVElement attempts to convert a VNode to a VElement.
// Returns the VElement and true if successful, otherwise nil and false.
func AsVElement(node VNode) (*VElement, bool) {
	if IsVElement(node) {
		return node.(*VElement), true
	}
	return nil, false
}

// IsVText checks if a node is a VText.
func IsVText(node VNode) bool {
	return node != nil && node.Type() == TextNode
}

// AsVText attempts to convert a VNode to a VText.
// Returns the VText and true if successful, otherwise nil and false.
func AsVText(node VNode) (*VText, bool) {
	if IsVText(node) {
		return node.(*VText), true
	}
	return nil, false
}