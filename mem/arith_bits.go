package mem

import (
	"errors"
	"github.com/bits-and-blooms/bitset"
	"math"
)

var ErrUnderflow = errors.New("sub operation underflowed")

// ArithBits is a smart ass
type ArithBits struct {
	B *bitset.BitSet
}

func NewArithBits(i uint) *ArithBits {
	b := bitset.New(i)
	return &ArithBits{B: b}
}

// Set did nothing but call inner bitset function
// we could pack functions like this as much as possible
func (a *ArithBits) Set(i uint) {
	a.B.Set(i)
}

func (a *ArithBits) Clr(i uint) {
	a.B.Clear(i)
}

func (a *ArithBits) Test(i uint) bool {
	return a.B.Test(i)
}

func (a *ArithBits) Add(add uint64) {
	b := a.B.Bytes()
	var carry bool

	if add == 0 {
		return
	}

	for i := 0; i < len(b); i++ {
		if i != 0 {
			if carry == false {
				a.B = bitset.From(b)
				return
			}
			if b[i] != math.MaxUint64 {
				b[i] = b[i] + 1
				a.B = bitset.From(b)
				return
			}
			b[i] = 0
			carry = true
			continue
		}
		// i == 0
		// overflowd
		if add > math.MaxUint64-b[i] {
			carry = true
		}
		// if overflowed, we take the residual
		b[i] = b[i] + add
	}
	if carry {
		b = append(b, 1)
	}
	a.B = bitset.From(b)
}

func (a *ArithBits) Sub(sub uint64) error {
	b := a.B.Bytes()
	var carry bool

	if sub == 0 {
		return nil
	}

	for i := 0; i < len(b); i++ {
		if i != 0 {
			if carry == false {
				a.B = bitset.From(b)
				return nil
			}
			if b[i] < sub {
				carry = true
				b[i] = math.MaxUint64 - sub + 1 + b[i]
			} else {
				b[i] = b[i] - sub
				a.B = bitset.From(b)
				return nil
			}
			continue
		}
	}
	// underflowed
	if carry {
		return ErrUnderflow
	}
	return nil
}

