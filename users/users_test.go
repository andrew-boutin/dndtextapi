// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserIDs(t *testing.T) {
	testIO := []struct {
		desc     string
		uc       UserCollection
		expected []int
	}{
		{
			desc:     "No users returns empty slice",
			uc:       UserCollection{},
			expected: []int{},
		},
		{
			desc:     "Single user returns slice with single id",
			uc:       UserCollection{&User{ID: 1}},
			expected: []int{1},
		},
		{
			desc:     "Multiple users returns slice with the user ids",
			uc:       UserCollection{&User{ID: 1}, &User{ID: 2}},
			expected: []int{1, 2},
		},
	}

	for _, test := range testIO {
		t.Run(test.desc, func(t *testing.T) {
			actual := test.uc.GetUserIDs()
			assert.Equal(t, test.expected, actual)
		})
	}
}
