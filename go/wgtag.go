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

import "github.com/aws/aws-sdk-go/service/athena"

// WGTags is a wrapper of []*athena.Tag.
type WGTags struct {
	tags []*athena.Tag
}

// NewWGTags is to create a new WGTags.
func NewWGTags() *WGTags {
	return &WGTags{tags: make([]*athena.Tag, 0, 2)}
}

// AddTag is to add tag.
func (t *WGTags) AddTag(k string, v string) {
	t.tags = append(t.tags, &athena.Tag{
		Key:   &k,
		Value: &v})
}

// Get is a getter.
func (t *WGTags) Get() []*athena.Tag {
	return t.tags
}
