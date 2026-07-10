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

import "testing"

func TestInputBoolean(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-boolean")
	b, isInherited, err := c.InputBoolean("yes")
	FailOnError(t, err)
	if isInherited {
		t.Errorf("expected isInherited=false")
	}
	if !b {
		t.Errorf("expected true")
	}
}

func TestInputBooleanInherited(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-boolean-inherited")
	b, isInherited, err := c.InputBoolean("<<inherit>>")
	FailOnError(t, err)
	if !isInherited {
		t.Errorf("expected isInherited=true")
	}
	if b {
		t.Errorf("expected zero value false, got %v", b)
	}
}

func TestInputInt(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-int")
	v, isInherited, err := c.InputInt("42")
	FailOnError(t, err)
	if isInherited {
		t.Errorf("expected isInherited=false")
	}
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestInputIntInherited(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-int-inherited")
	v, isInherited, err := c.InputInt("<<inherit>>")
	FailOnError(t, err)
	if !isInherited {
		t.Errorf("expected isInherited=true")
	}
	if v != 0 {
		t.Errorf("expected zero value 0, got %d", v)
	}
}

func TestInputStringOrList(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-list")
	v, isInherited, err := c.InputStringOrList("a,b,c")
	FailOnError(t, err)
	if isInherited {
		t.Errorf("expected isInherited=false")
	}
	if len(v) != 1 || v[0] != "a,b,c" {
		t.Errorf("unexpected list %v", v)
	}
}

func TestInputStringOrListInherited(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-list-inherited")
	v, isInherited, err := c.InputStringOrList("<<inherit>>")
	FailOnError(t, err)
	if !isInherited {
		t.Errorf("expected isInherited=true")
	}
	if v != nil {
		t.Errorf("expected zero value nil, got %v", v)
	}
}

func TestInputStringOrListNoInherit(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-list-no-inherit")
	v, isInherited, err := c.InputStringOrListNoInherit("a,b,c")
	FailOnError(t, err)
	if isInherited {
		t.Errorf("expected isInherited=false")
	}
	if len(v) != 1 || v[0] != "a,b,c" {
		t.Errorf("unexpected list %v", v)
	}
}

func TestInputStringOrListNoInheritInherited(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-list-no-inherit-inherited")
	v, isInherited, err := c.InputStringOrListNoInherit("<<inherit>>")
	FailOnError(t, err)
	if !isInherited {
		t.Errorf("expected isInherited=true")
	}
	if v != nil {
		t.Errorf("expected zero value nil, got %v", v)
	}
}

func TestInputStringOrDict(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-dict")
	v, isInherited, err := c.InputStringOrDict("k1=v1 k2=v2")
	FailOnError(t, err)
	if isInherited {
		t.Errorf("expected isInherited=false")
	}
	if v["k1"] != "v1" || v["k2"] != "v2" {
		t.Errorf("unexpected dict %v", v)
	}
}

func TestInputStringOrDictInherited(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-dict-inherited")
	v, isInherited, err := c.InputStringOrDict("<<inherit>>")
	FailOnError(t, err)
	if !isInherited {
		t.Errorf("expected isInherited=true")
	}
	if v != nil {
		t.Errorf("expected zero value nil, got %v", v)
	}
}

func TestInputStringOrDictNoInherit(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-dict-no-inherit")
	v, isInherited, err := c.InputStringOrDictNoInherit("k1=v1 k2=v2")
	FailOnError(t, err)
	if isInherited {
		t.Errorf("expected isInherited=false")
	}
	if v["k1"] != "v1" {
		t.Errorf("unexpected dict %v", v)
	}
}

func TestInputStringOrDictNoInheritInherited(t *testing.T) {
	c := createStubHTTPClientSingle(t, "input-string-or-dict-no-inherit-inherited")
	v, isInherited, err := c.InputStringOrDictNoInherit("<<inherit>>")
	FailOnError(t, err)
	if !isInherited {
		t.Errorf("expected isInherited=true")
	}
	if v != nil {
		t.Errorf("expected zero value nil, got %v", v)
	}
}
