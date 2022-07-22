package ipv4

import (
	"crypto/aes"
	"encoding/binary"
	"math/big"
	"time"
)

/*
 * cyclic provides an inexpensive approach to iterating over the IPv4 address
 * space in a random(-ish) manner such that we connect to every host once in
 * a scan execution without having to keep track of the IPs that have been
 * scanned or need to be scanned and such that each scan has a different
 * ordering. We accomplish this by utilizing a cyclic multiplicative group
 * of integers modulo a prime and generating a new primitive root (generator)
 * for each scan.
 *
 * We know that 3 is a generator of (Z mod 2^32 + 15 - {0}, *)
 * and that we have coverage over the entire address space because 2**32 + 15
 * is prime and ||(Z mod PRIME - {0}, *)|| == PRIME - 1. Therefore, we
 * just need to find a new generator (primitive root) of the cyclic group for
 * each scan that we perform.
 *
 * Because generators map to generators over an isomorphism, we can efficiently
 * find random primitive roots of our mult. group by finding random generators
 * of the group (Zp-1, +) which is isomorphic to (Zp*, *). Specifically the
 * generators of (Zp-1, +) are { s | (s, p-1) == 1 } which implies that
 * the generators of (Zp*, *) are { d^s | (s, p-1) == 1 }. where d is a known
 * generator of the multiplicative group. We efficiently find
 * generators of the additive group by precalculating the psub1_f of
 * p - 1 and randomly checking random numbers against the psub1_f until
 * we find one that is coprime and map it into Zp*. Because
 * totient(totient(p)) ~= 10^9, this should take relatively few
 * iterations to find a new generator.
 */

// Represents a multiplicative cyclic group (Z/pZ)*
type CyClicGroup struct {

	prime           uint64
	knownPrimroot   uint64
	primeFactors    []uint64
}

type Cycle struct {

	Group *CyClicGroup
	Generator uint64
	Order uint64
	Offset uint32
}

// We will pick the first cyclic group from this list that is
// larger than the number of IPs in our allowlist. E.g. for an
// entire Internet scan, this would be cyclic32
// Note: this list should remain ordered by size (primes) ascending.

var groups []CyClicGroup = []CyClicGroup{

	// 2^8 + 1
	{prime: 257,knownPrimroot: 3,primeFactors: []uint64{2}},
	// 2^16 + 1
	{prime: 65537,knownPrimroot: 3,primeFactors: []uint64{2}},
	// 2^24 + 43
	{prime: 16777259,knownPrimroot: 2,primeFactors: []uint64{2, 23, 103, 3541}},
	// 2^28 + 3
	{prime: 268435459,knownPrimroot: 2,primeFactors: []uint64{2, 3, 19, 87211}},
	// 2^32 + 15
	{prime: 4294967311,knownPrimroot: 3,primeFactors: []uint64{2, 3, 5, 131, 364289}},
}

func makeKey() []byte {

	key := make([]byte,16)
	send := time.Now().UnixNano()

	binary.LittleEndian.PutUint64(key, uint64(send))

	return key
}

func getRandInt() uint64 {

	var mask uint64 = 0xFFFFFFFF

	key := makeKey()

	cip,_:= aes.NewCipher(key)
	out := make([]byte,len(key))

	cip.Encrypt (out, key)

	return binary.BigEndian.Uint64(out[:8])&mask
}

func GetGroup(minSize uint64) *CyClicGroup {

	for _, g := range groups {
		if (g.prime > minSize) {
			return &g
		}
	}

	return nil
}

// Check whether an integer is coprime with (p - 1)
func (g *CyClicGroup) isCoprime(check uint64) bool {


	if check == 0 || check == 1 {
		return false
	}

	n := len(g.primeFactors)

	for i := 0; i < n; i++ {

		pf := g.primeFactors[i]

		if pf > check && (pf % check)==0 {
			return false
		} else if pf < check && (check % pf) == 0 {
			return false
		} else if pf == check {
			return false
		}
	}

	return true
}

func (g *CyClicGroup) GetOmorphism(additiveElt uint64) uint64{

	base := new(big.Int).SetUint64(g.knownPrimroot)
	power := new(big.Int).SetUint64(additiveElt)
	prime := new(big.Int).SetUint64(g.prime)

	primeRoot := new(big.Int)


	primeRoot = primeRoot.Exp(base,power,prime)

	return primeRoot.Uint64()
}

// Return a (random) number coprime with (p - 1) of the group,
// which is a generator of the additive group mod (p - 1)

func (g *CyClicGroup) findPrimroot() uint64 {

	var retv uint64
	var candidate uint64

	// The maximum primitive root we can return needs to be small enough such
	// that there is no overflow when multiplied by any element in the largest
	// group , which currently has p = 2^{32} + 15.
	var maxRoot uint64 = (uint64(1) << 32) - 14

	candidate = getRandInt() % g.prime

	// Repeatedly find a generator until we hit one that is small enough. For
	// the largest group, we have a very low probability of ever executing this
	// loop more than once, and for small groups it will only execute once.
	for {
		// Find an element that is coprime in the additive group
		for !g.isCoprime(candidate) {
			candidate = candidate+1
			candidate = candidate%g.prime
		}
		// Given a coprime element, apply the isomorphism.
		retv = g.GetOmorphism(candidate)

		if retv <= maxRoot {
			break
		}
	}

	return retv
}

func (g *CyClicGroup) MakeCycle() *Cycle {

	return &Cycle{
		Group:     g,
		Generator: g.findPrimroot(),
		Order:     g.prime -1,
		Offset:    uint32(getRandInt()%g.prime),
	}

}




