package core

import (
	"ProjectWIND/dice_expr"
	"ProjectWIND/wba"
	"math"
	"math/rand/v2"
)

type CalculatorInfo struct{}

func (c *CalculatorInfo) NewVector(dimension int, value ...float64) wba.Vector {
	if dimension < len(value) {
		value = value[:dimension]
	}
	if dimension > len(value) {
		value = append(value, make([]float64, dimension-len(value))...)
	}
	return wba.Vector{Dimension: dimension, Value: value}
}

func (c *CalculatorInfo) NewVectorFromArray(array []float64) wba.Vector {
	return c.NewVector(len(array), array...)
}

func (c *CalculatorInfo) VectorZero(dimension int) wba.Vector {
	return c.NewVector(dimension)
}

func (c *CalculatorInfo) VectorAdd(v ...wba.Vector) wba.Vector {
	if len(v) == 0 {
		return c.NewVector(0, []float64{}...)
	}
	for i := 0; i < len(v); i++ {
		if v[i].Dimension != v[0].Dimension {
			return c.NewVector(0, []float64{}...)
		}
	}
	result := c.NewVector(v[0].Dimension, make([]float64, v[0].Dimension)...)
	for i := 0; i < v[0].Dimension; i++ {
		for j := 0; j < len(v); j++ {
			result.Value[i] += v[j].Value[i]
		}
	}
	return result
}

func (c *CalculatorInfo) VectorSub(v1, v2 wba.Vector) wba.Vector {
	if v1.Dimension != v2.Dimension {
		return c.NewVector(0, []float64{}...)
	}
	result := c.NewVector(v1.Dimension, make([]float64, v1.Dimension)...)
	for i := 0; i < v1.Dimension; i++ {
		result.Value[i] = v1.Value[i] - v2.Value[i]
	}
	return result
}

func (c *CalculatorInfo) VectorDotProduct(v1, v2 wba.Vector) float64 {
	if v1.Dimension != v2.Dimension {
		return 0
	}
	var product float64
	for i := 0; i < v1.Dimension; i++ {
		product += v1.Value[i] * v2.Value[i]
	}
	return product
}

func (c *CalculatorInfo) VectorCrossProduct(v1, v2 wba.Vector) wba.Vector {
	if v1.Dimension != 3 || v2.Dimension != 3 {
		return c.NewVector(0, []float64{}...)
	}
	result := c.NewVector(3, make([]float64, 3)...)
	result.Value[0] = v1.Value[1]*v2.Value[2] - v1.Value[2]*v2.Value[1]
	result.Value[1] = v1.Value[2]*v2.Value[0] - v1.Value[0]*v2.Value[2]
	result.Value[2] = v1.Value[0]*v2.Value[1] - v1.Value[1]*v2.Value[0]
	return result
}

func (c *CalculatorInfo) VectorAngle(v1, v2 wba.Vector) float64 {
	if v1.Dimension < v2.Dimension {
		v1.Value = append(v1.Value, make([]float64, v2.Dimension-v1.Dimension)...)
	} else if v2.Dimension < v1.Dimension {
		v2.Value = append(v2.Value, make([]float64, v1.Dimension-v2.Dimension)...)
	}
	dotProduct := c.VectorDotProduct(v1, v2)
	module1 := v1.GetMagnitude()
	module2 := v2.GetMagnitude()
	if module1*module2 == 0 {
		return 0
	}
	return math.Acos(dotProduct / (module1 * module2))
}

func (c *CalculatorInfo) VectorDistance(v1, v2 wba.Vector) float64 {
	if v1.Dimension != v2.Dimension {
		return 0
	}
	var distance float64
	for i := 0; i < v1.Dimension; i++ {
		distance += (v1.Value[i] - v2.Value[i]) * (v1.Value[i] - v2.Value[i])
	}
	return math.Sqrt(distance)
}

func (c *CalculatorInfo) Floor(value float64) float64 {
	return math.Floor(value)
}

func (c *CalculatorInfo) Ceil(value float64) float64 {
	return math.Ceil(value)
}

func (c *CalculatorInfo) Round(value float64) float64 {
	return math.Round(value)
}

func (c *CalculatorInfo) Abs(value float64) float64 {
	return math.Abs(value)
}

func (c *CalculatorInfo) Sqrt(value float64) float64 {
	return math.Sqrt(value)
}

func (c *CalculatorInfo) Sin(value float64) float64 {
	return math.Sin(value)
}

func (c *CalculatorInfo) Cos(value float64) float64 {
	return math.Cos(value)
}

func (c *CalculatorInfo) Tan(value float64) float64 {
	return math.Tan(value)
}

func (c *CalculatorInfo) Arcsin(value float64) float64 {
	return math.Asin(value)
}

func (c *CalculatorInfo) Arccos(value float64) float64 {
	return math.Acos(value)
}

func (c *CalculatorInfo) Arctan(value float64) float64 {
	return math.Atan(value)
}

func (c *CalculatorInfo) Log(input float64, base float64) float64 {
	return math.Log(input) / math.Log(base)
}

func (c *CalculatorInfo) Log10(input float64) float64 {
	return math.Log10(input)
}

func (c *CalculatorInfo) Ln(input float64) float64 {
	return math.Log(input) / math.Log10(math.Sqrt2)
}

func (c *CalculatorInfo) Exp(input float64) float64 {
	return math.Exp(input)
}

func (c *CalculatorInfo) Pow(base float64, exponent float64) float64 {
	return math.Pow(base, exponent)
}

func (c *CalculatorInfo) Min(values ...float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

func (c *CalculatorInfo) Max(values ...float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	max := values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func (c *CalculatorInfo) RandF() float64 {
	return rand.Float64()
}

func (c *CalculatorInfo) RandI(from int64, to int64) int64 {
	if from > to {
		from, to = to, from
	}
	return rand.Int64N(to-from) + from
}

func (c *CalculatorInfo) RandRange(from float64, to float64) float64 {
	if from > to {
		from, to = to, from
	}
	return rand.Float64()*(to-from) + from
}

func (c *CalculatorInfo) D(n, x int64, v ...int64) int64 {
	if len(v) < 1 {
		v = append(v, 1)
	}
	if len(v) < 3 {
		v = append(v, make([]int64, 3-len(v))...)
	}
	if len(v) > 3 {
		v = v[:3]
	}
	k, p, q := v[0], v[1], v[2]
	var res int64
	for i := 0; i < int(n); i++ {
		res += c.RandI(1, x)
	}
	res += p
	res *= k
	res += q
	return res
}

func (c *CalculatorInfo) Roll(expression string) wba.Result {
	result, err := dice_expr.Evaluate(expression)
	text := expression
	if err != nil {
		text = err.Error()
	}
	return wba.Result{
		Expression: text,
		Value:      int64(result),
	}
}

var CalculatorApi CalculatorInfo
