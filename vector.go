package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func CosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, errors.New("vector dimensions do not match")
	}

	var dotProduct float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0, nil
	}

	similarity := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
	return float32(similarity), nil
}

func ParseVector(s string) ([]float32, error) {
	parts := strings.Split(s, ",")
	vec := make([]float32, len(parts))

	for i := 0; i < len(parts); i++ {
		f, err := strconv.ParseFloat(strings.TrimSpace(parts[i]), 32)
		if err != nil {
			return nil, err
		}
		vec[i] = float32(f)
	}

	return vec, nil
}