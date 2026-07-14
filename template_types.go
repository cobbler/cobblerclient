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

// TemplateSchema enumerates where a Template's content lives. Mirrors cobbler.enums.TemplateSchema. This is a plain
// string alias (not a Go int-backed enum) — see NetworkInterfaceType in network_interface_types.go for why: the
// XML-RPC codec picks its wire representation purely from reflect.Kind, so an int-kind type is always sent as
// <int>, which the corresponding Python enum field rejects with "0 must be a str or Enum".
type TemplateSchema = string

const (
	TemplateSchemaFile        = "file"
	TemplateSchemaImportlib   = "importlib"
	TemplateSchemaEnvironment = "environment"
)

// URIOption models the (schema, path) pair that locates a Template's content.
type URIOption struct {
	Schema TemplateSchema `mapstructure:"schema" json:"schema" yaml:"schema"`
	Path   string         `mapstructure:"path" json:"path" yaml:"path"`
}
