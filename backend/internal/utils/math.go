package utils

// 整数最小值
func MinInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// 整数最大值
func MaxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// 浮点数最小值
func MinFloat64(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// 浮点数最大值
func MaxFloat64(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
} 