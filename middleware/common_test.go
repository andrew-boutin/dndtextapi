// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"fmt"
	"net/http"
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

func TestQueryParamExtractor(t *testing.T) {
	testIO := []struct {
		desc          string
		queryParamSet bool
		queryParam    string
		expected      string
		expectedErr   error
	}{
		{
			desc:          "Valid query parameter.",
			queryParamSet: true,
			queryParam:    "id=1",
			expected:      "1",
			expectedErr:   nil,
		},
		{
			desc:          "Other query parameter set.",
			queryParamSet: true,
			queryParam:    "random=1",
			expected:      "",
			expectedErr:   ErrQueryParamNotFound,
		},
		{
			desc:          "No query parameter provided.",
			queryParamSet: false,
			queryParam:    "",
			expected:      "",
			expectedErr:   ErrQueryParamNotFound,
		},
	}

	for _, test := range testIO {
		t.Run(test.desc, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			url := "http://localhost:8080"

			if test.queryParamSet {
				url += fmt.Sprintf("?%s", test.queryParam)
			}

			c.Request, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:8080?%s", test.queryParam), nil)

			str, err := QueryParamExtractor(c, "id")

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expected, str)
		})
	}
}
