package fspatterns

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSizeWithinLimits(t *testing.T) {
	tests := []struct {
		name           string
		st             SizeThreshold
		actualSize     int64
		expectedResult bool
	}{
		{
			name:           "Exact size as threshold and condition is GreaterEqualThan returns true",
			st:             SizeThreshold{SizeInBytes: 100, Condition: GreaterEqualThan},
			actualSize:     100,
			expectedResult: true,
		},
		{
			name:           "SizeInBytes above threshold and condition is GreaterEqualThan returns true",
			st:             SizeThreshold{SizeInBytes: 100, Condition: GreaterEqualThan},
			actualSize:     150,
			expectedResult: true,
		},
		{
			name:           "SizeInBytes below threshold and condition is GreaterEqualThan returns false",
			st:             SizeThreshold{SizeInBytes: 100, Condition: GreaterEqualThan},
			actualSize:     50,
			expectedResult: false,
		},
		{
			name:           "Exact size as threshold and condition is LessThan returns false",
			st:             SizeThreshold{SizeInBytes: 100, Condition: LessThan},
			actualSize:     100,
			expectedResult: false,
		},
		{
			name:           "SizeInBytes above threshold and condition is LessThan returns false",
			st:             SizeThreshold{SizeInBytes: 100, Condition: LessThan},
			actualSize:     150,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.st.IsSizeWithinThreshold(tt.actualSize)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
