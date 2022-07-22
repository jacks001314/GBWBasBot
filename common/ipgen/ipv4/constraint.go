package ipv4

/*
 Efficient address-space constraints  (AH 7/2013)

 This module uses a tree-based representation to efficiently
 manipulate and query constraints on the address space to be
 scanned.  It provides a value for every IP address, and these
 values are applied by setting them for network prefixes.  Order
 matters: setting a value replaces any existing value for that
 prefix or subsets of it.  We use this to implement network
 allowlisting and blocklisting.

 Think of setting values in this structure like painting
 subnets with different colors.  We can paint subnets black to
 exclude them and white to allow them.  Only the top color shows.
 This makes for potentially very powerful constraint specifications.

 Internally, this is implemented using a binary tree, where each
 node corresponds to a network prefix.  (E.g., the root is
 0.0.0.0/0, and its children, if present, are 0.0.0.0/1 and
 128.0.0.0/1.)  Each leaf of the tree stores the value that applies
 to every address within the leaf's portion of the prefix space.

 As an optimization, after all values are set, we look up the
 value or subtree for every /16 prefix and cache them as an array.
 This lets subsequent lookups bypass the bottom half of the tree.
**/

const  RADIX_LENGTH = 20

 type node struct {

 	l *node
 	r *node
 	value uint32
 	count uint64
 }

 type Constraint struct {

 	// root node of the tree
 	root *node

 	// array of prefixes (/RADIX_LENGTH) that are painted
 	radix []uint32

 	// have we precomputed counts for each node?
 	painted bool

 	// value for which we precomputed counts
 	paintValue uint32
 }



func newLeaf(value uint32) *node {

	return &node{
		l:     nil,
		r:     nil,
		value: value,
		count: 0,
	}
}

// Tree operations respect the invariant that every node that isn't a
// leaf has exactly two children.
func isLeaf(n *node) bool {

	return n.l == nil
}

// Convert from an internal node to a leaf.
func convertToLeaf(n *node) {

	if !isLeaf(n){

		n.l = nil
		n.r = nil
	}
}


// Recursive function to set value for a given network prefix within
// the tree.  (Note: prefix must be in host byte order.)
func setRecurse(n *node,prefix uint32, len int,value uint32) {


	var mask uint32 = 0x80000000

	if n == nil{
		return
	}

	if len == 0 {
		// We're at the end of the prefix; make this a leaf and set the
		// value.
		if !isLeaf(n) {
			convertToLeaf(n)
		}
		n.value = value
		return
	}

	if isLeaf(n) {
		// We're not at the end of the prefix, but we hit a leaf.
		if n.value == value {
			// A larger prefix has the same value, so we're done.
			return
		}
		// The larger prefix has a different value, so we need to
		// convert it into an internal node and continue processing on
		// one of the leaves.
		n.l = newLeaf(n.value)
		n.r = newLeaf(n.value)
	}

	// We're not at the end of the prefix, and we're at an internal
	// node.  Recurse on the left or right subtree.
	if (prefix &mask) != 0 {
		setRecurse(n.r,prefix << 1, len - 1, value)
	} else {
		setRecurse(n.l,prefix << 1, len - 1, value)
	}

	// At this point, we're an internal node, and the value is set
	// by one of our children or its descendent.  If both children are
	// leaves with the same value, we can discard them and become a left.
	if isLeaf(n.r) && isLeaf(n.l) && (n.r.value == n.l.value) {
		n.value = n.l.value
		convertToLeaf(n)
	}
}

// Return the value pertaining to an address, according to the tree
// starting at given root.  (Note: address must be in host byte order.)
func lookupIP(root *node,address uint32) uint32 {

	var nd *node = root
	var mask uint32 = 0x80000000

	for {

		if isLeaf(nd) {

			return nd.value
		}

		if address&mask != 0 {
			nd = nd.r
		} else {
			nd = nd.l
		}
		mask =mask>>1
	}
}


// Return the value pertaining to an address.
// (Note: address must be in host byte order.)

func (con *Constraint)LookupIP(address uint32) uint32 {

	return lookupIP(con.root, address)
}

// Set the value for a given network prefix, overwriting any existing
// values on that prefix or subsets of it.
// (Note: prefix must be in host byte order.)
func (con *Constraint)Set(prefix uint32, len int,value uint32) {

	setRecurse(con.root, prefix, len, value)

	con.painted = false;
}

// Return the nth painted IP address.
func lookupIndex(root *node, n uint64) uint64 {

	var nd *node = root
	var ipv uint64 = 0
	var mask uint64 = 0x80000000

	for {
		if isLeaf(nd) {
			return ipv | n
		}
		if n < nd.l.count {
			nd = nd.l
		} else {
			n = n-nd.l.count
			nd = nd.r
			ipv = ipv|mask
		}

		mask = mask >>1
	}
}


// Return a node that determines the values for the addresses with
// the given prefix.  This is either the internal node that
// corresponds to the end of the prefix or a leaf node that
// encompasses the prefix. (Note: prefix must be in host byte order.)
func lookupNode(root *node,prefix uint32, len int) *node {


	var nd *node = root
	var mask uint32 = 0x80000000

	for i := 0; i < len; i++ {

		if isLeaf(nd){
			return nd
		}

		if (prefix & mask)!=0 {

			nd = nd.r
		}else {
			nd = nd.l
		}

		mask = mask >>1
	}

	return nd
}

// Implement count_ips by recursing on halves of the tree.  Size represents
// the number of addresses in a prefix at the current level of the tree.
// If paint is specified, each node will have its count set to the number of
// leaves under it set to value.
// If excludeRadix is specified, the number of addresses will exclude prefixes
// that are a /RADIX_LENGTH or larger
func countIPSRecurse(nd *node,value uint32, size uint64, paint bool, excludeRadix bool) uint64 {

	var n uint64 = 0

	if isLeaf(nd) {

		if nd.value == value {
			n = size
			// Exclude prefixes already included in the radix
			if excludeRadix && size >= (1 << (32 - RADIX_LENGTH)) {
				n = 0
			}
		} else {
			n = 0
		}
	} else {

		n = countIPSRecurse(nd.l, value, size >> 1, paint,
			excludeRadix) +
			countIPSRecurse(nd.r, value, size >> 1, paint,
				excludeRadix)
	}

	if paint {

		nd.count = n
	}

	return n
}

// Return the number of addresses that have a given value.
func (con *Constraint)CountIPS(value uint32) uint64 {

	if con.painted && con.paintValue == value {

		return con.root.count +uint64(len(con.radix)) * (1 << (32 - RADIX_LENGTH))
	} else {

		return countIPSRecurse(con.root, value, uint64(1 << 32),false,false)
	}
}

// For each node, precompute the count of leaves beneath it set to value.
// Note that the tree can be painted for only one value at a time.

func (con *Constraint)PaintValue(value uint32) {

	var prefix uint32
	// Paint everything except what we will put in radix
	countIPSRecurse(con.root, value, uint64(1 << 32), true, true)

	// Fill in the radix array with a list of addresses
	for i := 0; i < (1 << RADIX_LENGTH); i++ {

		prefix = uint32(i << (32 - RADIX_LENGTH))
		nd := lookupNode(con.root, prefix, RADIX_LENGTH)

		if isLeaf(nd) && nd.value == value {
			// Add this prefix to the radix
			con.radix = append(con.radix,prefix)
		}

	}

	con.painted = true
	con.paintValue = value
}


// For a given value, return the IP address with zero-based index n.
// (i.e., if there are three addresses with value 0xFF, looking up index 1
// will return the second one).
// Note that the tree must have been previously painted with this value.
//

func (con *Constraint)LookupIndex(index uint64,value uint32) uint64 {

	var radixIndex uint64
	var radixOffset uint64

	if !con.painted || con.paintValue != value {
		con.PaintValue(value)
	}

	radixIndex = index / (1 << (32 - RADIX_LENGTH))

	if radixIndex < uint64(len(con.radix)) {
		// Radix lookup

		radixOffset = uint64(index % (1 << (32 - RADIX_LENGTH)))

		return uint64(con.radix[radixIndex]) | radixOffset
	}

	// Otherwise, do the "slow" lookup in tree.
	// Note that tree counts do NOT include things in the radix,
	// so we subtract these off here.
	index = index - uint64(len(con.radix)* (1 << (32 - RADIX_LENGTH)))
	return lookupIndex(con.root, index)
}




// Initialize the tree.
// All addresses will initially have the given value.
func NewConstraint(value uint32) *Constraint {

	return &Constraint{
		root:       newLeaf(value),
		radix:      make([]uint32, 0),
		painted:    false,
		paintValue: 0,
	}
}



