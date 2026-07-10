package cobblerclient

import (
	"testing"

	"github.com/go-test/deep"
)

func TestFindItemsPaged(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-items-paged")
	var items []interface{}
	var nilMap map[string]interface{}
	item1 := map[string]interface{}{
		"autoinstall_meta":    nilMap,
		"children":            []interface{}{},
		"comment":             nil,
		"ctime":               0.0,
		"depth":               int64(0),
		"display_name":        nil,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"ks_meta":             nilMap,
		"mtime":               0.0,
		"name":                "grub-menu",
		"owners":              []interface{}{"admin"},
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "000000000000000000000000000000d1",
	}
	item2 := map[string]interface{}{
		"autoinstall_meta":    nilMap,
		"children":            []interface{}{},
		"comment":             nil,
		"ctime":               0.0,
		"depth":               int64(0),
		"display_name":        nil,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"ks_meta":             nilMap,
		"mtime":               0.0,
		"name":                "testmenu",
		"owners":              []interface{}{"admin"},
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "000000000000000000000000000000d2",
	}
	item3 := map[string]interface{}{
		"autoinstall_meta":    nilMap,
		"children":            []interface{}{},
		"comment":             nil,
		"ctime":               0.0,
		"depth":               int64(0),
		"display_name":        nil,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"ks_meta":             nilMap,
		"mtime":               0.0,
		"name":                "testmenu1",
		"owners":              []interface{}{"admin"},
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "000000000000000000000000000000d3",
	}
	items = append(items, item1, item2, item3)

	expectedResult := PagedSearchResult{
		FoundItems: items,
		PageInfo: PageInfo{
			Page:             1,
			PrevPage:         -1,
			NextPage:         -1,
			Pages:            []int{1},
			NumPages:         1,
			NumItems:         3,
			StartItem:        0,
			EndItem:          3,
			ItemsPerPage:     5,
			ItemsPerPageList: []int{10, 20, 50, 100, 200, 500},
		},
	}

	criteria := make(map[string]interface{}, 1)
	criteria["display_name"] = ""
	result, err := c.FindItemsPaged("menu", criteria, "", 1, 5)
	FailOnError(t, err)
	if diff := deep.Equal(*result, expectedResult); diff != nil {
		t.Error(diff)
	}
}

func TestGetItem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item")
	// The fixture reflects a get_item call for a system named "test" that does not
	// exist on the server; Cobbler returns the "~" not-found marker, which GetItem
	// translates into an empty map.
	res, err := c.GetItem("system", "test", false, false)
	FailOnError(t, err)
	if len(res) != 0 {
		t.Error("expected an empty result for a non-existent item")
	}
}

func TestGetItemFlattened(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-flattened")
	res, err := c.GetItem("system", "testsys", true, false)
	FailOnError(t, err)
	if res["profile"] != "00000000000000000000000000000013" {
		t.Error("expected a different profile")
	}
	if res["kernel_options"] != "<<inherit>>" {
		t.Error("expected different kernel options")
	}
}

func TestGetItemResolved(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-resolved")
	res, err := c.GetItem("system", "testsys", false, true)
	FailOnError(t, err)
	if res["profile"] != "00000000000000000000000000000013" {
		t.Error("expected a different profile")
	}
	if res["redhat_management_key"] != nil {
		t.Error("expected a different redhat_management_key")
	}
	if kernelOptions, ok := res["kernel_options"].(map[string]interface{}); !ok || len(kernelOptions) != 0 {
		t.Error("expected different kernel options")
	}
}

func TestFindItems(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-items")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test*"
	res, err := c.FindItems("profile", criteria, "name", false)
	FailOnError(t, err)
	if len(res) != 2 {
		t.Error("expected two result profiles")
	}
}

func TestFindItemNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-item-names")
	expectedResult := []string{"testprof", "testprof1"}
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test*"
	res, err := c.FindItemNames("profile", criteria, "name")
	FailOnError(t, err)
	if diff := deep.Equal(res, expectedResult); diff != nil {
		t.Error(diff)
	}
}

func TestModifyItem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "modify-item")
	// modify_item takes a raw item uid (the current API's get_item_handle no longer
	// returns "type::name"-style handles), so we pass the real uid from the fixture.
	err := c.ModifyItem("distro_group", "000000000000000000000000000000df", []string{"comment"}, "hello")
	FailOnError(t, err)
}

func TestModifyItemInPlace(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"modify-item-in-place-get",
		"modify-item-in-place-handle",
		"modify-item-in-place-modify",
		"modify-item-in-place-save",
	})
	err := c.ModifyItemInPlace("profile", "testprof", "kernel_options", map[string]interface{}{"test": "1"})
	FailOnError(t, err)
}

func TestGetItemResolvedValue(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-resolved-value")
	result, err := c.GetItemResolvedValue("000000000000000000000000000000df", []string{"kernel_options"})
	FailOnError(t, err)
	if result == nil {
		t.Error("Expected non-nil result.")
	}
}

func TestSetItemResolvedValue(t *testing.T) {
	c := createStubHTTPClientSingle(t, "set-item-resolved-value")
	err := c.SetItemResolvedValue("000000000000000000000000000000df", []string{"comment"}, "hello")
	FailOnError(t, err)
}

func TestHasItem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "has-item")
	// The fixture's server response is "false": the "testtemplate" template was
	// deleted earlier in the recording sequence and a recent server-side fix
	// (template deletion / URI schema comparison) means it is no longer found here,
	// as it might have been with the old, buggy comparison.
	exists, err := c.HasItem("template", "testtemplate")
	FailOnError(t, err)
	if exists {
		t.Error("Expected item to not exist.")
	}
}

func TestNewItemClient(t *testing.T) {
	c := createStubHTTPClientSingle(t, "new-item-client")
	err := c.NewItem("template", false)
	FailOnError(t, err)
}
