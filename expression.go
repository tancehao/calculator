package calculator

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrNotExpr = errors.New("the expression is not a mathmatical one")
)

const (
    O_RETURN = iota
    O_ADD
    O_MULTIPLE
)

type Expr struct {
	Operator int
	Operands []*Expr
	Value    int64

	data       []byte
	commaStack []int
}

func NewExpression(data []byte) *Expr {
	return &Expr{data: data}
}

func (e *Expr) Parse() error {
	var subExps [][]byte
	var prev int64
	//whether the expression is wholely in a comma
	//example: (1+2)
	var allInComma bool
	for i := 0; i < len(e.data); i++ {
		switch e.data[i] {
		case '(':
			e.commaStack = append(e.commaStack, i)
			if i == 0 {
				allInComma = true
			}
		case ')':
			l := len(e.commaStack)
			prevIndex := e.commaStack[l-1]
			if l == 0 || e.data[prevIndex] != '(' {
				return ErrNotExpr
			}
			if prevIndex == 0 {
				if i == len(e.data)-1 {
					allInComma = true
				} else {
					allInComma = false
				}
			}
			e.commaStack = e.commaStack[:l-1]
		case '+':
			if i == 0 || i == len(e.data)-1 {
				return ErrNotExpr
			}
			if len(e.commaStack) > 0 {
				//now in a subexpression
				continue
			}
			if e.Operator == O_MULTIPLE {
				/* a '*' was met before, the subExps stores the operands of a multiplication, now treat them as a whole subexpression */
				subExps = [][]byte{e.data[:i]}
			} else {
				subExps = append(subExps, e.data[prev:i])
			}
			e.Operator = O_ADD
			prev = int64(i + 1)
		case '*':
			if i == 0 || i == len(e.data)-1 {
				return ErrNotExpr
			}
			if len(e.commaStack) > 0 {
				//now in a subexpression
				continue
			}
			/* a '+' has been met before, so this multiplication should be a sub-expression */
			if e.Operator == O_ADD {
				continue
			}
			subExps = append(subExps, e.data[prev:i])
			e.Operator = O_MULTIPLE
			prev = int64(i + 1)
		}
	}
	if len(subExps) > 0 {
		subExps = append(subExps, e.data[prev:])
	}
	if len(e.commaStack) > 0 {
		return ErrNotExpr
	}
	if allInComma {
		e.data = e.data[1 : len(e.data)-1]
		return e.Parse()
	}
	if e.Operator == O_RETURN { //no operator found, this expression contains only number
		e.Value, _ = strconv.ParseInt(string(e.data), 10, 64)
		return nil
	}
	for _, subexp := range subExps {
		exp := NewExpression(subexp)
		err := exp.Parse()
		if err != nil {
			return err
		}
		e.Operands = append(e.Operands, exp)
	}
	return nil
}

func (e *Expr) Calculate() (out int64) {
	switch e.Operator {
    case O_RETURN:
        return e.Value
    case O_ADD:
		for i := 0; i < len(e.Operands); i++ {
			out += e.Operands[i].Calculate()
		}
		return
	case O_MULTIPLE:
		out = 1
		for i := 0; i < len(e.Operands); i++ {
			out *= e.Operands[i].Calculate()
		}
		return
	}
	return 0
}

func (e *Expr) Print() {
	if e.Value != 0 {
		fmt.Printf(" %d", e.Value)
	} else {
		fmt.Print(" (")
		defer fmt.Print(")")
		fmt.Printf("%s", string(e.Operator))
		for _, op := range e.Operands {
			op.Print()
		}
	}
}
