package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestPathParamExtractor(t *testing.T) {
	testIO := []struct {
		desc        string
		paramSet    bool
		param       gin.Param
		expected    string
		expectedErr error
	}{
		{
			desc:        "Valid path parameter.",
			paramSet:    true,
			param:       gin.Param{Key: "id", Value: "1"},
			expected:    "1",
			expectedErr: nil,
		},
		{
			desc:        "Other parameter set.",
			paramSet:    true,
			param:       gin.Param{Key: "random", Value: "1"},
			expected:    "",
			expectedErr: ErrPathParamNotFound,
		},
		{
			desc:        "No parameter set.",
			paramSet:    false,
			param:       gin.Param{},
			expected:    "",
			expectedErr: ErrPathParamNotFound,
		},
	}

	for _, test := range testIO {
		t.Run(test.desc, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			if test.paramSet {
				c.Params = gin.Params{test.param}
			}

			str, err := PathParamExtractor(c, "id")

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expected, str)
		})
	}
}
