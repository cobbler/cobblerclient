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
		"boot_files":          nilMap,
		"comment":             nil,
		"ctime":               1722927169.2400115,
		"depth":               int64(0),
		"display_name":        nil,
		"fetchable_files":     nilMap,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"mtime":               1722927169.2400115,
		"name":                "testmenu",
		"owners":              "<<inherit>>",
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "191d3a7a7dea425a9e0e5b8b51ddafab",
	}
	item2 := map[string]interface{}{
		"autoinstall_meta":    nilMap,
		"boot_files":          nilMap,
		"comment":             nil,
		"ctime":               1722932798.1049085,
		"depth":               int64(0),
		"display_name":        nil,
		"fetchable_files":     nilMap,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"mtime":               1722932798.1049085,
		"name":                "testmenu1",
		"owners":              "<<inherit>>",
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "f42a34f1f5564de49beba9e7a8fcd58f",
	}
	item3 := map[string]interface{}{
		"autoinstall_meta":    nilMap,
		"boot_files":          nilMap,
		"comment":             nil,
		"ctime":               1722932811.623434,
		"depth":               int64(0),
		"display_name":        nil,
		"fetchable_files":     nilMap,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"mtime":               1722932811.623434,
		"name":                "testmenu10",
		"owners":              "<<inherit>>",
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "80c03370ae9941928f2a64346f21a6ec",
	}
	item4 := map[string]interface{}{
		"autoinstall_meta":    nilMap,
		"boot_files":          nilMap,
		"comment":             nil,
		"ctime":               1722932799.4978173,
		"depth":               int64(0),
		"display_name":        nil,
		"fetchable_files":     nilMap,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"mtime":               1722932799.4978173,
		"name":                "testmenu2",
		"owners":              "<<inherit>>",
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "b63f624dcc8a42b19dae51f46c01bd91",
	}
	item5 := map[string]interface{}{
		"autoinstall_meta":    nilMap,
		"boot_files":          nilMap,
		"comment":             nil,
		"ctime":               1722932801.1373043,
		"depth":               int64(0),
		"display_name":        nil,
		"fetchable_files":     nilMap,
		"is_subobject":        false,
		"kernel_options":      nilMap,
		"kernel_options_post": nilMap,
		"mtime":               1722932801.1373043,
		"name":                "testmenu3",
		"owners":              "<<inherit>>",
		"parent":              nil,
		"template_files":      nilMap,
		"uid":                 "14b33123dc0b4a34af3eb83ad4344f27",
	}
	items = append(items, item1, item2, item3, item4, item5)

	expectedResult := PagedSearchResult{
		FoundItems: items,
		PageInfo: PageInfo{
			Page:             1,
			PrevPage:         -1,
			NextPage:         2,
			Pages:            []int{1, 2, 3},
			NumPages:         3,
			NumItems:         11,
			StartItem:        0,
			EndItem:          5,
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
	res, err := c.GetItem("system", "test", false, false)
	FailOnError(t, err)
	if res["profile"] != "Ubuntu-20.04-x86_64" {
		t.Error("expected a different profile")
	}
}

func TestGetItemFlattened(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-flattened")
	res, err := c.GetItem("system", "testsys", true, false)
	FailOnError(t, err)
	if res["profile"] != "testprof" {
		t.Error("expected a different profile")
	}
	if res["kernel_options"] != "test=1" {
		t.Error("expected different kernel options")
	}
}

func TestGetItemResolved(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-resolved")
	res, err := c.GetItem("system", "testsys", false, true)
	FailOnError(t, err)
	if res["profile"] != "testprof" {
		t.Error("expected a different profile")
	}
	if res["redhat_management_key"] != "test" {
		t.Error("expected a different redhat_management_key")
	}
	if len(deep.Equal(res["kernel_options"], map[string]interface{}{"test": "1"})) != 0 {
		t.Error("expected different kernel options")
	}
}

func TestFindItems(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-items")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test*"
	res, err := c.FindItems("profile", criteria, "name", false)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Error("expected a single result profile")
	}
}

func TestFindItemNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-item-names")
	expectedResult := []string{"testprof"}
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
	err := c.ModifyItem("distro_group", "distro_group::webservers", []string{"comment"}, "hello")
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
	result, err := c.GetItemResolvedValue("some-uuid", []string{"kernel_options"})
	FailOnError(t, err)
	if result == nil {
		t.Error("Expected non-nil result.")
	}
}

func TestSetItemResolvedValue(t *testing.T) {
	c := createStubHTTPClientSingle(t, "set-item-resolved-value")
	err := c.SetItemResolvedValue("some-uuid", []string{"comment"}, "hello")
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
