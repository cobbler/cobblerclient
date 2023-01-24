package item

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cobbler/cobblerclient/client"
	"github.com/cobbler/cobblerclient/internal/item"
	rawItem "github.com/cobbler/cobblerclient/internal/item/raw"
	resolvedItem "github.com/cobbler/cobblerclient/internal/item/resolved"
)

// Item general fields
type Item struct {
	// Internal fields
	raw      *rawItem.Item
	resolved *resolvedItem.Item

	// Internal XML-RPC handle
	Handle string

	// Item fields
	Parent            item.Property[string]
	Depth             item.Property[int]
	Children          item.Property[[]string]
	CTime             item.Property[float64]
	MTime             item.Property[float64]
	Uid               item.Property[string]
	Name              item.Property[string]
	Comment           item.Property[string]
	KernelOptions     item.InheritableProperty[map[string]string]
	KernelOptionsPost item.InheritableProperty[map[string]string]
	AutoinstallMeta   item.InheritableProperty[map[string]string]
	FetchableFiles    item.InheritableProperty[map[string]string]
	BootFiles         item.Property[map[string]string]
	TemplateFiles     item.Property[map[string]string]
	Owners            item.InheritableProperty[[]string]
	MgmtClasses       item.InheritableProperty[[]string]
	MgmtParameters    item.InheritableProperty[map[string]interface{}]

	// Connectivity to the server
	Client client.Client
}

func BuildItem(client client.Client) Item {
	i := Item{
		Client:   client,
		raw:      &rawItem.Item{},
		resolved: &resolvedItem.Item{},
	}
	refreshItemPointers(&i)
	return i
}

func refreshItemPointers(i *Item) {
	i.Parent = Parent{i}
	i.Depth = Depth{i}
	i.Children = Children{i}
	i.CTime = CTime{i}
	i.MTime = MTime{i}
	i.Uid = Uid{i}
	i.Name = Name{i}
	i.Comment = Comment{i}
	i.KernelOptions = KernelOptions{i}
	i.KernelOptionsPost = KernelOptionsPost{i}
	i.AutoinstallMeta = AutoinstallMeta{i}
	i.FetchableFiles = FetchableFiles{i}
	i.BootFiles = BootFiles{i}
	i.TemplateFiles = TemplateFiles{i}
	i.Owners = Owners{i}
	i.MgmtClasses = MgmtClasses{i}
	i.MgmtParameters = MgmtParameters{i}
}

// UpdateCobblerFields updates all fields in a Cobbler Item structure.
func (it *Item) UpdateCobblerFields(what string, item reflect.Value, id string) error {
	method := fmt.Sprintf("modify_%s", what)

	typeOfT := item.Type()
	// In Cobbler v3.3.0, if profile name isn't created first, an empty child gets written to the distro, which causes
	// a ValueError: "calling find with no arguments"  TO-DO: figure a more efficient way of targeting name.
	for i := 0; i < item.NumField(); i++ {
		v := item.Field(i)
		tag := typeOfT.Field(i).Tag
		field := tag.Get("mapstructure")
		if method == "modify_profile" && field == "name" {
			var value interface{}
			switch v.Type().String() {
			case "string", "bool", "int64", "int":
				value = v.Interface()
			case "[]string":
				value = strings.Join(v.Interface().([]string), " ")
			}
			_, err := it.Client.Call(method, id, field, value, it.Client.Token)
			if err != nil {
				return err
			}
		}
	}

	for i := 0; i < item.NumField(); i++ {
		v := item.Field(i)
		tag := typeOfT.Field(i).Tag
		field := tag.Get("mapstructure")
		cobblerTag := tag.Get("cobbler")

		if cobblerTag == "noupdate" {
			continue
		}

		if field == "" {
			continue
		}
		var value interface{}
		switch v.Type().String() {
		case "string", "bool", "int64", "int":
			value = v.Interface()
		case "[]string":
			value = strings.Join(v.Interface().([]string), " ")
		}
		if result, err := it.Client.Call(method, id, field, value, it.Client.Token); err != nil {
			return err
		} else {
			if result.(bool) == false && value != false {
				// It's possible this is a new field that isn't available on
				// older versions.
				if cobblerTag == "newfield" {
					continue
				}
				return fmt.Errorf("error updating %s to %s", field, value)
			}
		}
	}
	return nil
}

type Parent struct {
	item *Item
}

func (p Parent) Get() string {
	return p.item.raw.Parent
}

func (p Parent) Set(name string) {
	p.item.raw.Parent = name
}

type Depth struct {
	item *Item
}

func (d Depth) Get() int {
	return d.item.raw.Depth
}

func (d Depth) Set(depth int) {
	d.item.raw.Depth = depth
}

type Children struct {
	item *Item
}

func (c Children) Get() []string {
	return c.item.raw.Children
}

func (c Children) Set(children []string) {
	c.item.raw.Children = children
}

type CTime struct {
	item *Item
}

func (ct CTime) Get() float64 {
	return ct.item.raw.CTime
}

func (ct CTime) Set(ctime float64) {
	ct.item.raw.CTime = ctime
}

type MTime struct {
	item *Item
}

func (mt MTime) Get() float64 {
	return mt.item.raw.MTime
}

func (mt MTime) Set(mtime float64) {
	mt.item.raw.MTime = mtime
}

type Uid struct {
	item *Item
}

func (u Uid) Get() string {
	return u.item.raw.Uid
}

func (u Uid) Set(uid string) {
	u.item.raw.Uid = uid
}

type Name struct {
	item *Item
}

func (n Name) Get() string {
	return n.item.raw.Name
}

func (n Name) Set(name string) {
	n.item.raw.Name = name
}

type Comment struct {
	item *Item
}

func (c Comment) Get() string {
	return c.item.raw.Comment
}

func (c Comment) Set(comment string) {
	c.item.raw.Parent = comment
}

type KernelOptions struct {
	item *Item
}

func (ko KernelOptions) Get() map[string]string {
	return ko.item.resolved.KernelOptions
}

func (ko KernelOptions) GetRaw() interface{} {
	return ko.item.raw.KernelOptions
}

func (ko KernelOptions) Set(kernelOptions map[string]string) {
	ko.item.resolved.KernelOptions = kernelOptions
}

func (ko KernelOptions) SetRaw(kernelOptions interface{}) {
	ko.item.raw.KernelOptions = kernelOptions
}

func (ko KernelOptions) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(ko.item.raw.KernelOptions)
}

func (ko KernelOptions) IsInherited() (bool, error) {
	switch ko.item.KernelOptions.GetRaw().(type) {
	case map[string]string:
		// Maps are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}

type KernelOptionsPost struct {
	item *Item
}

func (ko KernelOptionsPost) Get() map[string]string {
	return ko.item.resolved.KernelOptionsPost
}

func (ko KernelOptionsPost) GetRaw() interface{} {
	return ko.item.raw.KernelOptionsPost
}

func (ko KernelOptionsPost) Set(kernelOptions map[string]string) {
	ko.item.resolved.KernelOptionsPost = kernelOptions
}

func (ko KernelOptionsPost) SetRaw(kernelOptions interface{}) {
	ko.item.raw.KernelOptionsPost = kernelOptions
}

func (ko KernelOptionsPost) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(ko.item.raw.KernelOptionsPost)
}

func (ko KernelOptionsPost) IsInherited() (bool, error) {
	switch ko.item.KernelOptionsPost.GetRaw().(type) {
	case map[string]string:
		// Maps are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}

type AutoinstallMeta struct {
	item *Item
}

func (ko AutoinstallMeta) Get() map[string]string {
	return ko.item.resolved.AutoinstallMeta
}

func (ko AutoinstallMeta) GetRaw() interface{} {
	return ko.item.raw.AutoinstallMeta
}

func (ko AutoinstallMeta) Set(autoinstallmeta map[string]string) {
	ko.item.resolved.AutoinstallMeta = autoinstallmeta
}

func (ko AutoinstallMeta) SetRaw(autoinstallmeta interface{}) {
	ko.item.raw.AutoinstallMeta = autoinstallmeta
}

func (ko AutoinstallMeta) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(ko.item.raw.AutoinstallMeta)
}

func (ko AutoinstallMeta) IsInherited() (bool, error) {
	switch ko.item.AutoinstallMeta.GetRaw().(type) {
	case map[string]string:
		// Maps are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}

type FetchableFiles struct {
	item *Item
}

func (f FetchableFiles) Get() map[string]string {
	return f.item.resolved.KernelOptionsPost
}

func (f FetchableFiles) GetRaw() interface{} {
	return f.item.raw.FetchableFiles
}

func (f FetchableFiles) Set(fetchableFiles map[string]string) {
	f.item.resolved.FetchableFiles = fetchableFiles
}

func (f FetchableFiles) SetRaw(fetchableFiles interface{}) {
	f.item.raw.FetchableFiles = fetchableFiles
}

func (f FetchableFiles) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(f.item.raw.FetchableFiles)
}

func (f FetchableFiles) IsInherited() (bool, error) {
	switch f.item.FetchableFiles.GetRaw().(type) {
	case map[string]string:
		// Maps are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}

type BootFiles struct {
	item *Item
}

func (c BootFiles) Get() map[string]string {
	return c.item.raw.BootFiles
}

func (c BootFiles) Set(bootfiles map[string]string) {
	c.item.raw.BootFiles = bootfiles
}

type TemplateFiles struct {
	item *Item
}

func (tf TemplateFiles) Get() map[string]string {
	return tf.item.raw.TemplateFiles
}

func (tf TemplateFiles) Set(templatefiles map[string]string) {
	tf.item.raw.TemplateFiles = templatefiles
}

type Owners struct {
	item *Item
}

func (o Owners) Get() []string {
	return o.item.resolved.Owners
}

func (o Owners) GetRaw() interface{} {
	return o.item.raw.Owners
}

func (o Owners) Set(kernelOptions []string) {
	o.item.resolved.Owners = kernelOptions
}

func (o Owners) SetRaw(kernelOptions interface{}) {
	o.item.raw.Owners = kernelOptions
}

func (o Owners) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(o.item.raw.Owners)
}

func (o Owners) IsInherited() (bool, error) {
	switch o.item.Owners.GetRaw().(type) {
	case []string:
		// Maps are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}

type MgmtClasses struct {
	item *Item
}

func (mp MgmtClasses) Get() []string {
	return mp.item.resolved.MgmtClasses
}

func (mp MgmtClasses) GetRaw() interface{} {
	return mp.item.raw.MgmtClasses
}

func (mp MgmtClasses) Set(mgmtclasses []string) {
	mp.item.resolved.MgmtClasses = mgmtclasses
}

func (mp MgmtClasses) SetRaw(mgmtclasses interface{}) {
	mp.item.raw.MgmtClasses = mgmtclasses
}

func (mp MgmtClasses) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(mp.item.raw.MgmtClasses)
}

func (mp MgmtClasses) IsInherited() (bool, error) {
	switch mp.item.MgmtClasses.GetRaw().(type) {
	case map[string]string:
		// Maps are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}

type MgmtParameters struct {
	item *Item
}

func (mp MgmtParameters) Get() map[string]interface{} {
	return mp.item.resolved.MgmtParameters
}

func (mp MgmtParameters) GetRaw() interface{} {
	return mp.item.raw.MgmtParameters
}

func (mp MgmtParameters) Set(kernelOptions map[string]interface{}) {
	mp.item.resolved.MgmtParameters = kernelOptions
}

func (mp MgmtParameters) SetRaw(kernelOptions interface{}) {
	mp.item.raw.MgmtParameters = kernelOptions
}

func (mp MgmtParameters) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(mp.item.raw.MgmtParameters)
}

func (mp MgmtParameters) IsInherited() (bool, error) {
	switch mp.item.MgmtParameters.GetRaw().(type) {
	case map[string]string:
		// Maps are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}
