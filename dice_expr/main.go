package dice_expr

import (
	"fmt"
	"math/rand"
	"unicode"
)

// 测试用例
//func main() {
//	fmt.Println(Evaluate("2d6+3"))                  // 正常骰子
//	fmt.Println(Evaluate("d"))                      // 默认1d100
//	fmt.Println(Evaluate("2dddd5"))                 // 连续d
//	fmt.Println(Evaluate("2*3+4"))                  // 纯算术
//	fmt.Println(Evaluate("2d2@+11"))                // 含有非法字符
//	fmt.Println(Evaluate("(2d6+3)*4+(1d100-50)/2")) //含有多级括号的表达式
//	fmt.Println(Evaluate("1dd(1d(d100))"))          //某故意刁难的测试用例
//	fmt.Println(Evaluate("d100+(d10d)d(d)"))        //某故意刁难的测试用例
//}

func Evaluate(dice_expr string, default_dice_size ...int) (int, error) {
	var dice_size int
	totalRolls := 0
	if len(default_dice_size) > 0 {
		dice_size = default_dice_size[0]
	} else {
		dice_size = 100
	}

	dice_expr = preprocessDiceExpr(dice_expr, dice_size)

	opnd := NewStack[int](100)  // 操作数栈
	oprt := NewStack[rune](100) // 运算符栈
	oprt.Push('#')              // 栈底标识

	i := 0
loop:
	for i < len(dice_expr) || oprt.Peek() != '#' {
		if i < len(dice_expr) && (dice_expr[i] == ' ' || dice_expr[i] == '\t') {
			i++ // 跳过空白字符
			continue
		}

		if i < len(dice_expr) && isDigit(dice_expr[i]) {
			// 处理数字
			num := 0
			for i < len(dice_expr) && isDigit(dice_expr[i]) {
				num = num*10 + int(dice_expr[i]-'0')
				i++
			}
			opnd.Push(num)
		} else {
			// 处理运算符
			var op rune
			if i < len(dice_expr) {
				op = rune(dice_expr[i])

				// 检查字符合法性
				if !isValidOperator(op) {
					// 遇到非法字符，视为结束
					op = '#'
					i = len(dice_expr) // 跳过剩余字符
				} else {
					i++
				}
			} else {
				op = '#'
			}
			if i >= len(dice_expr) && op == '#' && oprt.Peek() == '#' {
				break
			}

			switch PrecedenceTable.GetPrecedence(oprt.Peek(), op) {
			case '<':
				oprt.Push(op)
			case '=':
				popped := oprt.Pop()
				if popped == '#' && op == '#' {
					break loop // 跳出整个循环
				}
			case '>':
				if opnd.Size() < 2 {
					return 0, fmt.Errorf("缺少操作数")
				}
				b := opnd.Pop()
				a := opnd.Pop()
				operator := oprt.Pop()
				result, err := calculate(a, b, operator, &totalRolls) // 传递计数器
				if err != nil {
					return 0, err
				}
				opnd.Push(result)
				if op != '#' {
					i-- // 仅对非结束符回退
				}
			case '!':
				return 0, fmt.Errorf("语法错误: 不支持的运算符组合 %c 和 %c", oprt.Peek(), op)
			}
		}
	}

	if opnd.Size() != 1 {
		return 0, fmt.Errorf("表达式不完整")
	}
	return opnd.Pop(), nil
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isValidOperator(op rune) bool {
	validOps := []rune{'+', '-', '*', '/', 'd', '(', ')', '#'}
	for _, v := range validOps {
		if op == v {
			return true
		}
	}
	return false
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isExprEnd(r rune) bool {
	return r == ')' || r == '+' || r == '-' || r == '*' || r == '/' || r == 'd'
}

func calculate(a, b int, op rune, totalRolls *int) (int, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, fmt.Errorf("除零错误")
		}
		return a / b, nil
	case 'd':
		if b <= 0 {
			return 0, fmt.Errorf("骰子面数必须为正数")
		}
		return d(a, b, totalRolls)
	default:
		return 0, fmt.Errorf("未知运算符: %c", op)
	}
}

func d(a, b int, totalRolls *int) (int, error) {
	// 检查掷骰次数限制
	if *totalRolls+a > 100000000 {
		return 0, fmt.Errorf("掷骰次数超过上限（100,000,000次）")
	}

	// 更新计数器
	*totalRolls += a

	// 执行掷骰
	sum := 0
	for i := 0; i < a; i++ {
		sum += rand.Intn(b) + 1
	}
	return sum, nil
}

func preprocessDiceExpr(expr string, dice_size int) string {
	var processed []rune
	exprRunes := []rune(expr)

	for i := 0; i < len(exprRunes); i++ {
		current := exprRunes[i]

		if current == 'd' {
			// ===== 左侧检查 =====
			leftHasDigit := false
			leftIsExprEnd := false

			if i > 0 {
				// 向左跳过空格
				leftIndex := i - 1
				for leftIndex >= 0 && isSpace(exprRunes[leftIndex]) {
					leftIndex--
				}

				if leftIndex >= 0 {
					leftChar := exprRunes[leftIndex]
					if unicode.IsDigit(leftChar) {
						leftHasDigit = true
					}
					if isExprEnd(leftChar) {
						leftIsExprEnd = true
					}
				}
			}

			// 添加默认次数
			if !leftHasDigit && !leftIsExprEnd {
				processed = append(processed, '1')
			}
			processed = append(processed, 'd')

			// ===== 右侧检查 =====
			rightHasValidContent := false
			rightIndex := i + 1

			// 跳过空格
			for rightIndex < len(exprRunes) && isSpace(exprRunes[rightIndex]) {
				rightIndex++
			}

			// 检查右侧有效内容
			if rightIndex < len(exprRunes) {
				nextChar := exprRunes[rightIndex]
				// 有效内容：数字或左括号
				if unicode.IsDigit(nextChar) || nextChar == '(' {
					rightHasValidContent = true
				}
			}

			// 添加默认面数
			if !rightHasValidContent {
				processed = append(processed, []rune(fmt.Sprintf("%d", dice_size))...)
			}
		} else {
			processed = append(processed, current)
		}
	}

	return string(processed)
}

type Table struct {
	data map[rune]map[rune]rune
}

func (t *Table) GetPrecedence(op1, op2 rune) rune {
	return t.data[op1][op2]
}

var PrecedenceTable = Table{
	data: map[rune]map[rune]rune{
		'#': {'#': '=', '(': '<', '+': '<', '-': '<', ')': '!', '*': '<', '/': '<', 'd': '<'},
		'+': {'#': '>', '(': '<', '+': '>', '-': '>', ')': '>', '*': '<', '/': '<', 'd': '<'},
		'-': {'#': '>', '(': '<', '+': '>', '-': '>', ')': '>', '*': '<', '/': '<', 'd': '<'},
		'(': {'#': '!', '(': '!', '+': '<', '-': '<', ')': '=', '*': '<', '/': '<', 'd': '<'},
		')': {'#': '>', '(': '!', '+': '>', '-': '>', ')': '>', '*': '>', '/': '>', 'd': '>'},
		'*': {'#': '>', '(': '<', '+': '>', '-': '>', ')': '>', '*': '>', '/': '>', 'd': '<'},
		'/': {'#': '>', '(': '<', '+': '>', '-': '>', ')': '>', '*': '>', '/': '>', 'd': '<'},
		'd': {'#': '>', '(': '<', '+': '>', '-': '>', ')': '>', '*': '>', '/': '>', 'd': '>'},
	},
}

// ... 其他代码保持不变 ...

type Stack[T comparable] struct {
	data []T
	top  int
}

func NewStack[T comparable](size int) *Stack[T] {
	return &Stack[T]{
		data: make([]T, size),
		top:  -1,
	}
}

func (s *Stack[T]) Push(value T) {
	if s.top == len(s.data)-1 {
		s.data = append(s.data, value)
	} else {
		s.data[s.top+1] = value
	}
	s.top++
}

func (s *Stack[T]) Pop() T {
	if s.top == -1 {
		var zero T
		return zero
	}
	s.top--
	return s.data[s.top+1]
}

func (s *Stack[T]) Peek() T {
	if s.top == -1 {
		var zero T
		return zero
	}
	return s.data[s.top]
}

func (s *Stack[T]) Empty() bool {
	return s.top == -1
}

func (s *Stack[T]) Size() int {
	return s.top + 1
}
