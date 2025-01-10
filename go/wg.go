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
	"github.com/aws/aws-sdk-go-v2/aws"
	v2athena "github.com/aws/aws-sdk-go-v2/service/athena"
	v2athenatypes "github.com/aws/aws-sdk-go-v2/service/athena/types"
)

// Workgroup is a wrapper of Athena Workgroup.
type Workgroup struct {
	Name   string
	Config *v2athenatypes.WorkGroupConfiguration
	Tags   *WGTags
}

// NewDefaultWG is to create new default Workgroup.
func NewDefaultWG(name string, config *v2athenatypes.WorkGroupConfiguration, tags *WGTags) *Workgroup {
	wg := Workgroup{
		Name:   name,
		Config: config,
	}
	if config == nil {
		wg.Config = GetDefaultWGConfig()
	}
	if tags != nil {
		wg.Tags = tags
	} else {
		wg.Tags = NewWGTags()
	}
	return &wg
}

// NewWG is to create a new Workgroup.
func NewWG(name string, config *v2athenatypes.WorkGroupConfiguration, tags *WGTags) *Workgroup {
	return &Workgroup{
		Name:   name,
		Config: config,
		Tags:   tags,
	}
}

// getWG is to get Athena Workgroup from AWS remotely.
func getWG(ctx context.Context, client AthenaClient, Name string) (*v2athenatypes.WorkGroup, error) {
	if client == nil {
		return nil, ErrAthenaNilClient
	}
	getWorkGroupOutput, err := client.GetWorkGroup(ctx,
		&v2athena.GetWorkGroupInput{
			WorkGroup: aws.String(Name),
		})
	if err != nil {
		return nil, err
	}
	return getWorkGroupOutput.WorkGroup, nil
}

// CreateWGRemotely is to create a Workgroup remotely.
func (w *Workgroup) CreateWGRemotely(ctx context.Context, athenaClient AthenaClient) error {
	tags := w.Tags.Get()
	var err error
	if len(tags) == 0 {
		_, err = athenaClient.CreateWorkGroup(ctx, &v2athena.CreateWorkGroupInput{
			Configuration: w.Config,
			Name:          aws.String(w.Name),
		})
	} else {
		_, err = athenaClient.CreateWorkGroup(ctx, &v2athena.CreateWorkGroupInput{
			Configuration: w.Config,
			Name:          aws.String(w.Name),
			Tags:          w.Tags.Get(),
		})
	}
	return err
}
