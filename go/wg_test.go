// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package athenadriver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWG(t *testing.T) {
	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber ID", "123456")
	wgTags.AddTag("Uber Role", "SDE")
	wg := NewWG("henry_wu", nil, wgTags)
	assert.Equal(t, wg.Name, "henry_wu")
	assert.Equal(t, len(wg.Tags.Get()), 3)
}

func TestGetWG(t *testing.T) {
	w, e := getWG(context.Background(), nil, "SELECT_OK")
	assert.Nil(t, w)
	assert.NotNil(t, e)

	athenaClient := newMockAthenaClient()
	w, e = getWG(context.Background(), athenaClient, "SELECT_OK")
	assert.Nil(t, w)
	assert.NotNil(t, e)

	athenaClient.GetWGStatus = true
	w, e = getWG(context.Background(), athenaClient, "SELECT_OK")
	assert.NotNil(t, w)
	assert.Nil(t, e)
}

func TestWorkgroup_CreateWGRemotely(t *testing.T) {
	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber ID", "123456")
	wgTags.AddTag("Uber Role", "SDE")
	wg := NewWG("henry_wu", nil, wgTags)
	athenaClient := newMockAthenaClient()
	e := wg.CreateWGRemotely(athenaClient)
	assert.NotNil(t, e)
	athenaClient.CreateWGStatus = true
	e = wg.CreateWGRemotely(athenaClient)
	assert.Nil(t, e)
}

func TestWorkgroup_CreateWGRemotely2(t *testing.T) {
	wgTags := NewWGTags()
	wg := NewWG("henry_wu", nil, wgTags)
	athenaClient := newMockAthenaClient()
	e := wg.CreateWGRemotely(athenaClient)
	assert.NotNil(t, e)
	athenaClient.CreateWGStatus = true
	e = wg.CreateWGRemotely(athenaClient)
	assert.Nil(t, e)
}
