package brainfxxk

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const (
	DEFAULT_DEST = -1
	DEBUG        = false
)

type Method struct {
	f    func()
	dest int
}

type Brainfxxk struct {
	pointer     int
	array       []int
	inputbuffer []byte
	reader      *bufio.Reader
	process     []Method
	eip         int
	looplabel   []int
}

func New() *Brainfxxk {
	return &Brainfxxk{
		pointer:     0,
		array:       []int{0},
		inputbuffer: []byte{},
		reader:      bufio.NewReader(os.Stdin),
		process:     []Method{},
		eip:         0,
		looplabel:   []int{},
	}
}

func (b *Brainfxxk) pinc() {
	b.pointer++
	if len(b.array)-1 < b.pointer {
		b.array = append(b.array, 0)
	}
	if DEBUG {
		fmt.Printf("%d : Pointer Increment (Point: %d,Data: %d)\n", b.eip, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) pdec() {
	b.pointer--
	if b.pointer < 0 {
		fmt.Fprint(os.Stderr, "Runtime Error : segment error")
		os.Exit(1)
	}
	if DEBUG {
		fmt.Printf("%d : Pointer Decrement (Point: %d,Data: %d)\n", b.eip, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) dinc() {
	b.array[b.pointer]++
	if DEBUG {
		fmt.Printf("%d : Data Increment (Point: %d,Data: %d)\n", b.eip, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) ddec() {
	b.array[b.pointer]--
	if DEBUG {
		fmt.Printf("%d : Data Decrement (Point: %d,Data: %d)\n", b.eip, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) put() {
	fmt.Print(string(b.array[b.pointer]))
	if DEBUG {
		fmt.Printf("%d : Print (Point: %d,Data: %d)\n", b.eip, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) get() {
	if len(b.inputbuffer) == 0 {
		input, _ := b.reader.ReadString('\n')
		b.inputbuffer = []byte(input)
	}

	b.array[b.pointer] = int(b.inputbuffer[0])
	b.inputbuffer = b.inputbuffer[1:]
	if DEBUG {
		fmt.Printf("%d : Get (Point: %d,Data: %d)\n", b.eip, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) loopstart() {
	d := b.process[b.eip].dest
	if b.array[b.pointer] == 0 {
		b.eip = d
	}
	if DEBUG {
		fmt.Printf("%d : if 0 then jump to %d (Point: %d,Data: %d)\n", b.eip, d, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) loopend() {
	d := b.process[b.eip].dest
	if b.array[b.pointer] != 0 {
		b.eip = d
	}
	if DEBUG {
		fmt.Printf("%d : if not 0 then jump to %d (Point: %d,Data: %d)\n", b.eip, d, b.pointer, b.array[b.pointer])
	}
}

func (b *Brainfxxk) Run() {
	processEnd := len(b.process)
	b.eip = 0
	for b.eip < processEnd {
		b.process[b.eip].f()
		b.eip++
	}
}

func (b *Brainfxxk) Add(order byte) error {
	switch order {
	case '>':
		b.process = append(b.process, Method{b.pinc, DEFAULT_DEST})
	case '<':
		b.process = append(b.process, Method{b.pdec, DEFAULT_DEST})
	case '+':
		b.process = append(b.process, Method{b.dinc, DEFAULT_DEST})
	case '-':
		b.process = append(b.process, Method{b.ddec, DEFAULT_DEST})
	case '.':
		b.process = append(b.process, Method{b.put, DEFAULT_DEST})
	case ',':
		b.process = append(b.process, Method{b.get, DEFAULT_DEST})
	case '[':
		b.process = append(b.process, Method{b.loopstart, DEFAULT_DEST})
		b.looplabel = append(b.looplabel, b.eip)
	case ']':
		llen := len(b.looplabel)
		if llen <= 0 {
			return errors.New("Invalid Order")
		}
		dest := b.looplabel[llen-1]
		b.looplabel = b.looplabel[:llen-1]
		b.process = append(b.process, Method{b.loopend, dest})
		b.process[dest].dest = b.eip
	default:
		return errors.New(fmt.Sprintf("Invalid Order : %s", string(order)))
	}

	b.eip++
	return nil
}
