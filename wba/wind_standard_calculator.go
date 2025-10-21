package wba

import "math"

// WindStandardCalculator Wind 标准计算接口，包含掷骰、随机数、取整、向量计算、统计计算等功能
type WindStandardCalculator interface {
	NewVector(dimension int, value ...float64) Vector
	NewVectorFromArray(array []float64) Vector
	VectorZero(dimension int) Vector
	VectorAdd(v ...Vector) Vector
	VectorSub(v1, v2 Vector) Vector
	VectorDotProduct(v1, v2 Vector) float64
	VectorCrossProduct(v1, v2 Vector) Vector
	VectorAngle(v1, v2 Vector) float64
	VectorDistance(v1, v2 Vector) float64
	Floor(value float64) float64
	Ceil(value float64) float64
	Round(value float64) float64
	Abs(value float64) float64
	Sqrt(value float64) float64
	Sin(value float64) float64
	Cos(value float64) float64
	Tan(value float64) float64
	Arcsin(value float64) float64
	Arccos(value float64) float64
	Arctan(value float64) float64
	Log(input float64, base float64) float64
	Log10(input float64) float64
	Ln(input float64) float64
	Exp(input float64) float64
	Pow(base float64, exponent float64) float64
	Max(values ...float64) float64
	Min(values ...float64) float64
	RandF() float64
	RandI(from int64, to int64) int64
	RandRange(from float64, to float64) float64
	D(n, x int64, v ...int64) int64
	Roll(expression string) Result
}

// Vector 向量
type Vector struct {
	Dimension int
	Value     []float64
}

func (v Vector) GetDimension() int {
	return v.Dimension
}

func (v Vector) GetMagnitude() float64 {
	var module float64
	for _, value := range v.Value {
		module += value * value
	}
	return math.Sqrt(module)
}

func (v Vector) GetUnitVector() Vector {
	module := v.GetMagnitude()
	unitVector := Vector{Dimension: v.Dimension, Value: make([]float64, v.Dimension)}
	for i := 0; i < v.Dimension; i++ {
		unitVector.Value[i] = v.Value[i] / module
	}
	return unitVector
}

func (v Vector) GetItem(index int) float64 {
	if math.Abs(float64(index)) > float64(v.Dimension) {
		return math.NaN()
	}
	if index < 0 {
		index += v.Dimension
	}
	return v.Value[index]
}

func (v Vector) Multiply(scalar float64) Vector {
	result := Vector{Dimension: v.Dimension, Value: make([]float64, v.Dimension)}
	for i := 0; i < v.Dimension; i++ {
		result.Value[i] = v.Value[i] * scalar
	}
	return result
}

type Result struct {
	Expression           string
	NormalizedExpression string
	Value                int64
}
