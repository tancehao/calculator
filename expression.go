package calculator

import (
    "errors"
    "strconv"
    "fmt"
)

var (
    ErrNotExpr = errors.New("the expression is not a mathmatical one")
)

type Expr struct {
    Operator byte
    Operands []*Expr
    Value int64
    Parent *Expr

    data []byte
    commaStack []byte
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
        switch (e.data[i]) {
        case '(':
            e.commaStack = append(e.commaStack, '(')
            if i == 0 {
                allInComma = true
            }
        case ')':
            l := len(e.commaStack)
            if l == 0 || e.commaStack[l-1] != '(' {
                return ErrNotExpr
            }
            if i != len(e.data) -1 {
                allInComma = false
            }
            e.commaStack = e.commaStack[:l - 1]
        case '+':
            if i == 0 || i == len(e.data) - 1 {
                return ErrNotExpr
            }
            if len(e.commaStack) > 0 {
                //now in a subexpression
                continue
            }
            if e.Operator == '*' {
                /* a '*' was met before, the subExps stores the operands of a multiplication, now treat them as a whole subexpression */
                subExps = [][]byte{e.data[:i]}
            } else {
                subExps = append(subExps, e.data[prev:i])
            }
            e.Operator = '+'
            prev = int64(i + 1)
        case '*':
            if i == 0 || i == len(e.data) - 1 {
                return ErrNotExpr
            }
            if len(e.commaStack) > 0 {
                //now in a subexpression
                continue
            }
            /* a '+' was met before, so this multiplication should be a sub-expression */
            if e.Operator == '+' {
                continue
            }
            subExps = append(subExps, e.data[prev:i])
            e.Operator = '*'
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
        e.data = e.data[1:len(e.data)-1]
        return e.Parse()
    }
    if e.Operator == byte(0) {  //no operator found, this expression contains only number
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
    fmt.Println("Operator:", string(e.Operator), "Operands", e.Operands)
    return nil
}

//execute an expression defined by an operator and a list of operands
/*
func Exec(operator func(*Expr, *Expr) int64, operands ...*Expr) int64 {
    l := len(operands)
    if l == 0 {
        return int64(0)
    } else if l == 1 {
        return operands[0].Value
    } else {
        return operator(Exec(operator, operands[:l-1]...), operands[l-1])
    }
}
*/

func (e *Expr) Calculate() (out int64) {
    if e.Value > 0 {
        return e.Value
    }
    if len(e.Operands) == 0 {
        return int64(0)
    } else if len(e.Operands) == 1 {
        return e.Operands[0].Calculate()
    }

    switch (e.Operator) {
    case '+':
        for i := 1; i < len(e.Operands); i++ {
            out = e.Operands[i-1].Calculate() + e.Operands[i].Calculate()
        }
        return
    case '*':
        for i := 1; i < len(e.Operands); i++ {
            out = e.Operands[i-1].Calculate() * e.Operands[i].Calculate()
        }
        return
    }
    return int64(0)
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
