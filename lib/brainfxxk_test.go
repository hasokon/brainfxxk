package brainfxxk

import "testing"

func TestBrainfxxk(t *testing.T) {
	b := New()

	prog := []byte("+++++++++[>++++++++>+++++++++++>+++++<<<-]>.>++.+++++++..+++.>-.------------.<++++++++.--------.+++.------.--------.>+.>++++++++++.")
	for _, o := range prog {
		err := b.Add(o)
		if err != nil {
			t.Error(err.Error())
		}
	}

	b.Run()
}
