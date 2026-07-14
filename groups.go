/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

import (
	"fmt"
)

// Group is the shared structure for the three 4.0.0 group item types. Each
// concrete group (DistroGroup, ProfileGroup, SystemGroup) embeds Group with
// `mapstructure:",squash"` so callers see typed signatures while the wire
// representation stays uniform.
type Group struct {
	Item    `mapstructure:",squash" yaml:",inline"`
	Members []string `mapstructure:"members" json:"members" yaml:"members"`
}

func newGroup() Group {
	return Group{Item: NewItem(), Members: []string{}}
}

func convertRawGroup(what, name string, xmlrpcResult interface{}, dest interface{}) error {
	if xmlrpcResult == "~" {
		return fmt.Errorf("%s %s not found", what, name)
	}
	_, err := decodeCobblerItem(xmlrpcResult, dest)
	return err
}
