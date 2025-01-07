package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Program struct {
	RegisterA          int
	RegisterB          int
	RegisterC          int
	InstructionPointer int
	Output             []string
	Instructions       []int
}

func (p *Program) eval() {
	mp := map[int]func(int){
		0: p.adv,
		1: p.bxl,
		2: p.bst,
		3: p.jnz,
		4: p.bxc,
		5: p.out,
		6: p.bdv,
		7: p.cdv,
	}

	for {
		i := p.InstructionPointer
		if i >= len(p.Instructions) {
			return
		}

		opcode, operand := p.Instructions[i], p.Instructions[i+1]
		op := mp[opcode]
		op(operand)

		if p.InstructionPointer == i {
			p.InstructionPointer += 2
		}
	}
}

func (p *Program) adv(operand int) {
	denominator := math.Pow(2, float64(p.combo(operand)))
	p.RegisterA /= int(denominator)
}

func (p *Program) bxl(operand int) {
	p.RegisterB ^= operand
}

func (p *Program) bst(operand int) {
	p.RegisterB = p.combo(operand) % 8
}

func (p *Program) jnz(operand int) {
	if p.RegisterA != 0 {
		p.InstructionPointer = operand
	}
}

func (p *Program) bxc(operand int) {
	p.RegisterB ^= p.RegisterC
}

func (p *Program) out(operand int) {
	p.Output = append(p.Output, strconv.Itoa(p.combo(operand)%8))
}

func (p *Program) bdv(operand int) {
	denominator := math.Pow(2, float64(p.combo(operand)))
	p.RegisterB = p.RegisterA / int(denominator)
}

func (p *Program) cdv(operand int) {
	denominator := math.Pow(2, float64(p.combo(operand)))
	p.RegisterC = p.RegisterA / int(denominator)
}

func (p Program) combo(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return p.RegisterA
	case 5:
		return p.RegisterB
	case 6:
		return p.RegisterC
	default:
		return 0
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) Program {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	rereg := regexp.MustCompile(`Register ([ABC]): (\d+)`)
	repgm := regexp.MustCompile(`Program: ([0-9,]+)`)
	pg := Program{}

	scanner.Scan()
	setRegister(&pg, rereg, scanner.Text())
	scanner.Scan()
	setRegister(&pg, rereg, scanner.Text())
	scanner.Scan()
	setRegister(&pg, rereg, scanner.Text())
	scanner.Scan()
	scanner.Scan()

	data := repgm.FindStringSubmatch(scanner.Text())
	pg.Instructions = parseInstructions(data[1])

	check(scanner.Err())
	return pg
}

func atoi(val string) int {
	i, err := strconv.Atoi(val)
	check(err)
	return i
}

func parseInstructions(str string) []int {
	data := strings.Split(str, ",")
	inst := make([]int, 0)
	for i := 0; i < len(data); i += 2 {
		inst = append(inst, atoi(data[i]), atoi(data[i+1]))
	}
	return inst
}

func setRegister(pg *Program, rg *regexp.Regexp, line string) {
	m := rg.FindStringSubmatch(line)
	switch m[1] {
	case "A":
		pg.RegisterA = atoi(m[2])
	case "B":
		pg.RegisterB = atoi(m[2])
	case "C":
		pg.RegisterC = atoi(m[2])
	}
}

func main() {
	c := readLines(os.Args[1])
	c.eval()
	fmt.Printf("Result: %v\n", strings.Join(c.Output, ","))
}
