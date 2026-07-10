// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright SUSE LLC

// Run this program against a live Cobbler server to refresh all XML-RPC fixture files under ../fixtures/. To keep
// re-recordings byte-identical, every saved file is normalised: the session token becomes "securetoken99", any
// epoch-timestamp double (ctime/mtime/event timestamps) becomes 0.0, any "YYYY-MM-DD_HHMMSS_" event-id timestamp
// prefix becomes a fixed constant, every 32-char lowercase hex id (uid, system_uid, and similar references) is
// replaced with a deterministic placeholder assigned in order of first appearance, and every <struct>'s <member>
// elements are sorted by name (Go's XML-RPC library builds struct members from map[string]interface{}, whose key
// order Go randomizes on every iteration — see sortStructMembers).
//
// Known residual, non-deterministic fixtures this does NOT fully stabilise (unrelated to the above): get-event-log
// (a raw per-task log with real timestamps/thread names/durations and possible concurrent-task interleaving),
// get-random-mac (the call is intentionally random), and any fixture touching Template.tags (a Python set, whose
// iteration order Cobbler's server randomizes per process — cobbler/items/template.py).
//
// Usage: go run ./cmd/
//
// Prerequisites:
//   - The Cobbler server's collections must be empty before each run (this tool creates all the baseline test data
//     it needs itself, via seedBaselineData and the per-item-type create blocks below):
//     docker exec cobbler-dev bash -c "shopt -s globstar; rm -f /var/lib/cobbler/collections/**/*.json && supervisorctl restart cobblerd"
//   - A real (can be empty) file must exist on the server for the Template create/update flow to point its
//     uri.path at (cobbler/items/template.py requires a template's autoinstall path to already exist as a file —
//     it never creates one):
//     docker exec cobbler-dev touch /var/lib/cobbler/templates/testtemplate.template
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	cobbler "github.com/cobbler/cobblerclient"
)

// RecordingHTTPClient intercepts every XML-RPC POST, pretty-prints and saves the request and response as fixture files,
// and normalises the live session token to "securetoken99" so that fixtures are stable across sessions.
//
// Call SetRealToken after Login() succeeds so that the login fixtures are re-written with the normalised token and all
// subsequent fixtures are normalised on the fly.
type RecordingHTTPClient struct {
	inner     *http.Client
	realToken string
	names     []string
	idx       int
	ids       *idRemapper
}

func newRecordingHTTPClient(names []string, ids *idRemapper) *RecordingHTTPClient {
	return &RecordingHTTPClient{inner: http.DefaultClient, names: names, ids: ids}
}

// idRemapper assigns each distinct 32-char lowercase hex id (uid, system_uid, and similar references, as well as
// the uuid segment of event/task ids) a deterministic placeholder the first time it's seen, and reuses that mapping
// on every subsequent occurrence. Shared across the main recorder and all sub-recorders (see makeSubClient) so that
// a value recorded by one and echoed back through another still maps to the same placeholder.
type idRemapper struct {
	seen  map[string]string
	count int
}

func newIDRemapper() *idRemapper {
	return &idRemapper{seen: make(map[string]string)}
}

func (m *idRemapper) sub(real string) string {
	if placeholder, ok := m.seen[real]; ok {
		return placeholder
	}
	placeholder := fmt.Sprintf("%032x", m.count)
	m.count++
	m.seen[real] = placeholder
	return placeholder
}

// RealToken returns the live session token used for normalisation.
func (r *RecordingHTTPClient) RealToken() string { return r.realToken }

// Post implements cobbler.HTTPClient.
func (r *RecordingHTTPClient) Post(uri, bodyType string, req io.Reader) (*http.Response, error) {
	body, err := io.ReadAll(req)
	if err != nil {
		return nil, fmt.Errorf("recorder: reading request: %w", err)
	}
	if r.idx >= len(r.names) {
		return nil, fmt.Errorf("recorder: fixture sequence exhausted; unexpected extra call #%d", r.idx+1)
	}
	name := r.names[r.idx]
	r.idx++

	saveFixture(name+"-req.xml", prettyXML(normalize(body, r.realToken, r.ids)))

	resp, err := r.inner.Post(uri, bodyType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("recorder: forwarding to server: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("recorder: reading response: %w", err)
	}

	saveFixture(name+"-res.xml", prettyXML(normalize(respBody, r.realToken, r.ids)))

	// Return unmodified response so the Go client can use the real token.
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	return resp, nil
}

// SetRealToken stores the live token and re-normalises any fixtures already written before the token was known
// (typically login-res.xml).
func (r *RecordingHTTPClient) SetRealToken(token string) {
	if token == "" || token == r.realToken {
		return
	}
	r.realToken = token
	for i := 0; i < r.idx; i++ {
		for _, suffix := range []string{"-req.xml", "-res.xml"} {
			path := "./fixtures/" + r.names[i] + suffix
			data, err := os.ReadFile(path)
			if err != nil || !bytes.Contains(data, []byte(token)) {
				continue
			}
			saveFixture(r.names[i]+suffix, prettyXML(normalize(data, token, r.ids)))
		}
	}
}

var (
	// epochDoubleRe matches an XML-RPC <double> holding a Unix epoch timestamp (ctime, mtime, or an event's
	// state-time), recognised by its 10-digit integer part. No other <double> value in the fixture corpus falls in
	// this range.
	epochDoubleRe = regexp.MustCompile(`<double>1\d{9}(\.\d+)?</double>`)
	// eventTimestampPrefixRe matches the "YYYY-MM-DD_HHMMSS_" prefix Cobbler embeds in every event/task id
	// (CobblerEvent.__generate_event_id), wherever it occurs — as an XML-RPC string value or as a struct member name.
	eventTimestampPrefixRe = regexp.MustCompile(`\d{4}-\d{2}-\d{2}_\d{6}_`)
	// hexIDRe matches a uuid4-hex id (uid, system_uid, and similar uid-valued references, plus the uuid segment of
	// event/task ids), wherever it occurs — as a struct member value or a positional method-call parameter.
	hexIDRe = regexp.MustCompile(`[0-9a-f]{32}`)
)

// normalize replaces every source of run-to-run drift in a recorded XML-RPC payload with a deterministic
// placeholder: the live session token, epoch timestamps (ctime/mtime/event state-time), the timestamp prefix of
// event/task ids, and uuid4-hex ids (uid, system_uid, and similar references). ids must be shared across every
// recorder that can see the same underlying value (see makeSubClient) so a value recorded by one and echoed back
// through another still maps to the same placeholder.
func normalize(data []byte, token string, ids *idRemapper) []byte {
	if token != "" {
		data = bytes.ReplaceAll(data, []byte(token), []byte("securetoken99"))
	}
	data = epochDoubleRe.ReplaceAll(data, []byte("<double>0.0</double>"))
	data = eventTimestampPrefixRe.ReplaceAll(data, []byte("2000-01-01_000000_"))
	data = hexIDRe.ReplaceAllFunc(data, func(match []byte) []byte {
		return []byte(ids.sub(string(match)))
	})
	return data
}

// xmlNode is a minimal generic XML tree, used only to canonicalize <struct> member order (see sortStructMembers)
// before re-serializing. Go's XML-RPC client library (kolo/xmlrpc) builds <struct> elements from map[string]interface{}
// values, whose key order Go deliberately randomizes on every iteration — so requests built from a multi-key map
// (e.g. register_new_system) otherwise vary byte-for-byte between recorder runs even though member order is not
// semantically significant (cobblerclient's own test harness already sorts struct members before comparing request
// fixtures against what the client actually sends, for exactly this reason — see sortStruct in testing.go).
type xmlNode struct {
	xml.StartElement
	children []*xmlNode
	text     string
}

func parseXMLTree(dec *xml.Decoder) (*xmlNode, error) {
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if start, ok := tok.(xml.StartElement); ok {
			return buildXMLNode(dec, start.Copy())
		}
	}
}

func buildXMLNode(dec *xml.Decoder, start xml.StartElement) (*xmlNode, error) {
	n := &xmlNode{StartElement: start}
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			child, err := buildXMLNode(dec, t.Copy())
			if err != nil {
				return nil, err
			}
			n.children = append(n.children, child)
		case xml.EndElement:
			return n, nil
		case xml.CharData:
			// Skip whitespace-only text nodes; the encoder adds its own indentation. This also matters when
			// re-normalizing an already pretty-printed fixture (see SetRealToken), whose indentation would
			// otherwise show up here as text content.
			if len(bytes.TrimSpace(t)) > 0 {
				n.text += string(t)
			}
		}
	}
}

// sortStructMembers recursively sorts the <member> children of every <struct> element by their <name> child's text,
// mirroring testing.go's sortStruct so that struct-member order is canonicalized the same way the test harness
// already treats it: insignificant.
func sortStructMembers(n *xmlNode) {
	for _, c := range n.children {
		sortStructMembers(c)
	}
	if n.Name.Local == "struct" {
		sort.SliceStable(n.children, func(i, j int) bool {
			return memberName(n.children[i]) < memberName(n.children[j])
		})
	}
}

func memberName(member *xmlNode) string {
	for _, c := range member.children {
		if c.Name.Local == "name" {
			return c.text
		}
	}
	return ""
}

func encodeXMLNode(e *xml.Encoder, n *xmlNode) error {
	if err := e.EncodeToken(n.StartElement); err != nil {
		return err
	}
	if len(n.children) == 0 && n.text != "" {
		if err := e.EncodeToken(xml.CharData(n.text)); err != nil {
			return err
		}
	}
	for _, c := range n.children {
		if err := encodeXMLNode(e, c); err != nil {
			return err
		}
	}
	return e.EncodeToken(n.StartElement.End())
}

// prettyXML reformats XML with 4-space indentation, preserving the original XML declaration, and canonicalizes
// <struct> member order (see sortStructMembers).
func prettyXML(data []byte) []byte {
	s := strings.TrimSpace(string(data))

	// Extract and preserve the XML declaration verbatim.
	decl, body := "", s
	if strings.HasPrefix(s, "<?xml") {
		if idx := strings.Index(s, "?>"); idx >= 0 {
			decl = strings.TrimSpace(s[:idx+2])
			body = strings.TrimSpace(s[idx+2:])
		}
	}

	root, err := parseXMLTree(xml.NewDecoder(strings.NewReader(body)))
	if err != nil {
		return data // fallback: return original bytes
	}
	sortStructMembers(root)

	var enc bytes.Buffer
	e := xml.NewEncoder(&enc)
	e.Indent("", "    ")
	if err := encodeXMLNode(e, root); err != nil {
		return data // fallback: return original bytes
	}
	if err := e.Flush(); err != nil {
		return data
	}

	var out bytes.Buffer
	if decl != "" {
		out.WriteString(decl)
		out.WriteByte('\n')
	}
	out.Write(enc.Bytes())
	b := out.Bytes()
	if len(b) > 0 && b[len(b)-1] != '\n' {
		out.WriteByte('\n')
	}
	return out.Bytes()
}

func saveFixture(filename string, data []byte) {
	if err := os.WriteFile("./fixtures/"+filename, data, 0644); err != nil {
		fmt.Printf("  ERROR: %s: %v\n", filename, err)
	} else {
		fmt.Printf("  %s\n", filename)
	}
}

var config = cobbler.ClientConfig{
	URL:      "http://localhost:80/cobbler_api",
	Username: "cobbler",
	Password: "cobbler",
}

// Paths that exist on the cobbler-dev container's filesystem, used to satisfy Cobbler's requirement that a distro's
// kernel (and, for realism, initrd) resolve to a real file. Cobbler's own system-tests ship these as empty
// placeholder files precisely for this purpose.
const (
	fakeKernel   = "/code/system-tests/images/fake/vmlinuz"
	fakeInitrd   = "/code/system-tests/images/fake/initramfs"
	fakeVirtPath = "/var/lib/libvirt/images"
	// fakeVirtCpus is required alongside fakeVirtPath: cobbler/settings has no default_virt_cpus fallback (unlike
	// default_virt_ram/type/bridge/disk_driver/file_size), so any profile/image without a parent to inherit from
	// must set virt.cpus to a concrete value or every subsequent read raises AttributeError trying to resolve it.
	fakeVirtCpus = 2
)

// fixtureSequence lists every fixture name in the exact order that the corresponding HTTP calls are made by main().
// Each entry maps 1-to-1 to one XML-RPC POST. Per item type, calls are ordered Create, Update, Save, Copy, Rename,
// Read, Find, Delete; the Create/Update steps themselves are multi-call sub-recorder sequences (see makeSubClient)
// and are therefore not listed here — only the calls made through the main recorder (rec/c) are.
var fixtureSequence = []string{
	// ── LOGIN ──────────────────────────────────────────────────────────────
	"login",
	"extended-version",

	// ── AUTH MISC ──────────────────────────────────────────────────────────
	"check-access-no-fail",
	"check-access",
	// "get-authn-module-name",
	"token-check",
	"get-user-from-token",

	// ── VERSION ────────────────────────────────────────────────────────────
	"version",

	// ── CLIENT MISC ────────────────────────────────────────────────────────
	"get-tftp-file",
	"last-modified-time",
	"ping",
	"auto-add-repos",
	"is-autoinstall-in-use",
	"generate-ipxe",
	"run-install-triggers",
	"get-repos-compatible-with-profile",
	"find-system-by-dns-name",
	"get-random-mac",
	"get-status-normal",
	"get-status-text",
	"sync-dhcp",

	// ── ACTIONS ────────────────────────────────────────────────────────────
	"sync",
	"background-sync",
	"background-sync-systems",
	"check",
	"background-buildiso",
	"background-hardlink",
	"background-validate-autoinstall-files",
	"background-replicate",
	"background-aclsetup",
	"background-import",
	"background-reposync",
	"background-mkloaders",
	"background-power-system",

	// ── SIGNATURES ─────────────────────────────────────────────────────────
	"get-signatures",
	"get-valid-breeds",
	"get-valid-os-verions-for-breed", // intentional typo matches existing fixture name
	"get-valid-os-versions",
	"get-valid-archs",
	"background-signature-update",
	"background-signature-reload",

	// ── DISTROS (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-distro-handle",
	"save-distro",
	"copy-distro",
	"rename-distro",
	"get-distros",
	"get-distro",
	"get-item-names-distro",
	"get-distros-since",
	"get-valid-distro-boot-loaders",
	"get-distro-as-rendered",
	"find-distro",
	"find-distro-names",

	// ── PROFILES (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-profile-handle",
	"save-profile",
	"copy-profile",
	"rename-profile",
	"get-profiles",
	"get-profile",
	"get-item-names-profile",
	"get-profiles-since",
	"get-valid-profile-boot-loaders",
	"get-profile-as-rendered",
	"new-subprofile",
	"generate-autoinstall",
	"generate-boot-cfg",
	"generate-script",
	"find-profile",
	"find-profile-names",
	"register-new-system",

	// ── SYSTEMS (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-system-handle",
	"save-system",
	"copy-system",
	"rename-system",
	"get-systems",
	"get-system",
	"get-item-names-system",
	"get-system-since",
	"get-valid-system-boot-loaders",
	"get-system-as-rendered",
	"get-interfaces-get-system",
	"disable-netboot",
	"upload-log-data",
	"clear-system-logs",
	"find-system",
	"find-system-names", // same FindSystem call, different fixture name
	"dump-vars",
	"get-blended-data",
	"get-config-data",
	"power-system",
	"delete-system",
	"delete-profile",
	"delete-distro",

	// ── TEMPLATES (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-template-handle",
	"save-template",
	"copy-template",
	"rename-template",
	"get-template-file-for-profile",
	"get-template-file-for-system",
	"get-templates",
	"get-template",
	"get-templates-since",
	"get-item-names-template",
	"get-template-content",
	"templates-refresh-content",
	"background-templates-refresh-content",
	"find-template",
	"find-template-names",
	"delete-template",

	// ── REPOS (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-repo-handle",
	"save-repo",
	"copy-repo",
	"rename-repo",
	"get-repos",
	"get-repo",
	"get-repo-explicit-string", // same GetRepo call, fixture name differs
	"get-item-names-repo",
	"get-repo-since",
	"get-repo-config-for-profile",
	"get-repo-config-for-system",
	"get-repo-as-rendered",
	"find-repo",
	"find-repo-names",
	"delete-repo",

	// ── IMAGES (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-image-handle",
	"save-image",
	"copy-image",
	"rename-image",
	"get-images",
	"get-image",
	"get-item-names-image",
	"get-images-since",
	"find-image",
	"find-image-names",
	"delete-image",

	// ── MENUS (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-menu-handle",
	"save-menu",
	"copy-menu",
	"rename-menu",
	"get-menus",
	"get-menu",
	"get-item-names-menu",
	"get-menus-since",
	"get-menu-as-rendered",
	"find-menu",
	"find-menu-names",
	"delete-menu",

	// ── NETWORK INTERFACES (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-network-interface-handle",
	"save-network-interface",
	"copy-network-interface",
	"rename-network-interface",
	"get-network-interface",
	"get-network-interfaces",
	"get-item-names-network-interface",
	"get-network-interfaces-since",
	"find-network-interface",
	"find-network-interface-names",
	"delete-network-interface",

	// ── DISTRO GROUP (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-distro-group-handle",
	"save-distro-group",
	"copy-distro-group",
	"rename-distro-group",
	"get-distro-group",
	"get-distro-groups",
	"get-item-names-distro-group",
	"get-distro-groups-since",
	"get-distro-group-as-rendered",
	"find-distro-group",
	"delete-distro-group",

	// ── PROFILE GROUP (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-profile-group-handle",
	"save-profile-group",
	"copy-profile-group",
	"rename-profile-group",
	"get-profile-group",
	"get-profile-groups",
	"get-item-names-profile-group",
	"get-profile-groups-since",
	"get-profile-group-as-rendered",
	"find-profile-group",
	"delete-profile-group",

	// ── SYSTEM GROUP (Save, Copy, Rename, Read, Find, Delete — Create/Update are sub-recorders) ──
	"get-system-group-handle",
	"save-system-group",
	"copy-system-group",
	"rename-system-group",
	"get-system-group",
	"get-system-groups",
	"get-item-names-system-group",
	"get-system-groups-since",
	"get-system-group-as-rendered",
	"find-system-group",
	"delete-system-group",

	// ── ITEMS (generic, single-step) ───────────────────────────────────────
	// Note: ModifyItemInPlace (4 calls) uses its own sub-recorder.
	"find-items-paged",
	"get-item",
	"get-item-flattened",
	"get-item-resolved",
	"find-items",
	"find-item-names",
	"modify-item",
	"get-item-resolved-value",
	"set-item-resolved-value",
	"has-item",
	"new-item-client",

	// ── SETTINGS ───────────────────────────────────────────────────────────
	"get-settings",
	"modify-settings",

	// ── EVENTS ─────────────────────────────────────────────────────────────
	"get-task-status",
	"get-events",
	"get-event-log",

	// ── TRANSACTIONS ───────────────────────────────────────────────────────
	"transaction-begin",
	"transaction-commit",
	"transaction-abort",

	// ── INPUT HELPERS ──────────────────────────────────────────────────────
	"input-boolean",
	"input-int",
	"input-string-or-list",
	"input-string-or-list-no-inherit",
	"input-string-or-dict",
	"input-string-or-dict-no-inherit",

	// ── LOGOUT ─────────────────────────────────────────────────────────────
	"logout",
}

func warn(err error) {
	if err != nil {
		fmt.Printf("  [WARN] %v\n", err)
	}
}

// silentHTTPClient forwards HTTP calls to the real server without recording any fixtures. Used when we need to set up
// server state (e.g. begin a transaction) without consuming a fixture slot in the main sequence.
type silentHTTPClient struct{}

func (s *silentHTTPClient) Post(uri, bodyType string, req io.Reader) (*http.Response, error) {
	return http.DefaultClient.Post(uri, bodyType, req)
}

// makeSubClient creates an isolated recording client for a multi-step operation. Failures in the operation only corrupt
// that operation's own fixture slots, not the main recorder's sequence.
func makeSubClient(names []string, mainRec *RecordingHTTPClient, token string, version cobbler.CobblerVersion) (cobbler.Client, *RecordingHTTPClient) {
	sub := newRecordingHTTPClient(names, mainRec.ids)
	sub.SetRealToken(mainRec.RealToken())
	c := cobbler.NewClient(sub, config)
	c.Token = token
	c.CachedVersion = version
	return c, sub
}

// newSilentClient returns a client that talks to the real server without recording any fixtures. Used for lookups
// and setup calls that must not consume a slot in the main recorder's fixed sequence — e.g. resolving the real
// object id a Copy* call produced before it can be passed to Rename*, which (like modify_item and friends) requires
// a genuine uid rather than a "type::name" string.
func newSilentClient(token string, version cobbler.CobblerVersion) cobbler.Client {
	cl := cobbler.NewClient(&silentHTTPClient{}, config)
	cl.Token = token
	cl.CachedVersion = version
	return cl
}

// seedBaselineData creates the baseline distros/profiles/systems/repos/menus/images/network interfaces that the
// fixture calls below reference by a hard-coded name (e.g. "test", "testprof", "testsys", "centos7-x86_64") but
// never create themselves — mirroring the server-side test data the *_test.go fixtures assume already exists.
// Uses the silent client so no fixture slots are consumed.
func seedBaselineData(token string, version cobbler.CobblerVersion) {
	silent := newSilentClient(token, version)

	d := cobbler.NewDistro()
	d.Name = "test"
	d.Kernel = fakeKernel
	d.Initrd = fakeInitrd
	_, err := silent.CreateDistro(d)
	warn(err)

	// testprof/centos7-x86_64 get their own distro, separate from "test": the "test" distro is deleted alongside
	// the "test" profile at the end of the systems section below, and Cobbler refuses to delete a distro that
	// still has profile dependents ("removal would orphan profile(s)") — testprof/centos7-x86_64 live for the
	// entire run and must not be attached to a distro that gets torn down partway through.
	longLivedDistro := cobbler.NewDistro()
	longLivedDistro.Name = "testprof-distro"
	longLivedDistro.Kernel = fakeKernel
	longLivedDistro.Initrd = fakeInitrd
	_, err = silent.CreateDistro(longLivedDistro)
	warn(err)

	// Profile.distro/System.profile take the referenced item's real uid, not its name — modify_profile/modify_system
	// pass the value straight through to the Python setter, which resolves it via a strict uid-keyed dict lookup
	// (cobbler/items/profile.py's distro setter calls self.api.distros().find(uid=...), and
	// cobbler_collections/collection.py's fast path is `self.listing.get(uid)`, keyed by uid — not name).
	testDistroUID, err := silent.GetDistroHandle("test")
	warn(err)
	longLivedDistroUID, err := silent.GetDistroHandle("testprof-distro")
	warn(err)

	testProfile := cobbler.NewProfile()
	testProfile.Name = "test"
	testProfile.Distro = testDistroUID
	testProfile.Virt.Path = fakeVirtPath
	testProfile.Virt.Cpus = cobbler.Value[int]{Data: fakeVirtCpus}
	_, err = silent.CreateProfile(testProfile)
	warn(err)

	for _, name := range []string{"testprof", "centos7-x86_64"} {
		p := cobbler.NewProfile()
		p.Name = name
		p.Distro = longLivedDistroUID
		p.Virt.Path = fakeVirtPath
		p.Virt.Cpus = cobbler.Value[int]{Data: fakeVirtCpus}
		// ModifyItemInPlace needs a concrete (non-inherited) dict to merge into — NewProfile() defaults
		// KernelOptions to inherited, which serializes as the bare "<<inherit>>" sentinel string rather than a
		// dict, and ModifyItemInPlace's map[string]interface{} cast rejects that.
		p.KernelOptions = cobbler.Value[map[string]interface{}]{Data: map[string]interface{}{}}
		_, err := silent.CreateProfile(p)
		warn(err)
	}

	testprofUID, err := silent.GetProfileHandle("testprof")
	warn(err)

	for _, name := range []string{"test", "testsys", "server1"} {
		sys := cobbler.NewSystem()
		sys.Name = name
		// GetConfigData looks systems up by hostname, not name (cobbler/remote.py's get_config_data takes a
		// hostname straight through to api.get_config_data) — keep them equal for baseline systems.
		sys.Hostname = name
		sys.Profile = testprofUID
		// A concrete power type avoids an internal power-manager error when PowerSystem is later called against
		// a system whose power.type was left at the default inherited value.
		sys.Power.Type = "ipmilanplus"
		_, err := silent.CreateSystem(sys)
		warn(err)
	}

	for _, name := range []string{"test", "testrepo"} {
		r := cobbler.NewRepo()
		r.Name = name
		_, err := silent.CreateRepo(r)
		warn(err)
	}

	for _, name := range []string{"test", "testmenu"} {
		m := cobbler.NewMenu()
		m.Name = name
		_, err := silent.CreateMenu(m)
		warn(err)
	}

	img := cobbler.NewImage()
	img.Name = "test"
	img.Virt.Path = fakeVirtPath
	img.Virt.Cpus = cobbler.Value[int]{Data: fakeVirtCpus}
	_, err = silent.CreateImage(img)
	warn(err)

	serverHandle, err := silent.GetSystemHandle("server1")
	warn(err)
	ni := cobbler.NewNetworkInterface()
	ni.Name = "eth0-server1"
	// Left blank deliberately: cobbler/items/network_interface.py's mac_address setter only enforces uniqueness
	// for a non-empty value, and this interface gets cloned by CopyNetworkInterface later, which would otherwise
	// collide with the clone on a duplicate MAC.
	_, err = silent.CreateNetworkInterface(serverHandle, ni)
	warn(err)
}

func main() {
	// ids is shared by every recorder (main and sub-clients) for the lifetime of the run, so a uid/system_uid/event-id
	// recorded by one and echoed back through another still maps to the same placeholder.
	ids := newIDRemapper()

	// ── LOGIN-ERR: separate mini-recorder with wrong credentials ──────────
	fmt.Println("=== login-err ===")
	errRec := newRecordingHTTPClient([]string{"login-err"}, ids)
	cErr := cobbler.NewClient(errRec, cobbler.ClientConfig{
		URL:      config.URL,
		Username: "wrong",
		Password: "wrong",
	})
	_, _ = cErr.Login() // expected to fail; saves login-err-{req,res}.xml

	// ── MAIN RECORDER ─────────────────────────────────────────────────────
	rec := newRecordingHTTPClient(fixtureSequence, ids)
	c := cobbler.NewClient(rec, config)

	fmt.Println("=== login / extended-version ===")
	if _, err := c.Login(); err != nil {
		fmt.Printf("FATAL: login failed: %v\n", err)
		return
	}
	rec.SetRealToken(c.Token)

	fmt.Println("=== seeding baseline test data ===")
	seedBaselineData(c.Token, c.CachedVersion)

	// ── AUTH MISC ─────────────────────────────────────────────────────────
	fmt.Println("=== auth misc ===")
	_, err := c.CheckAccessNoFail("", "", "")
	warn(err)
	_, err = c.CheckAccess("", "", "")
	warn(err)
	// Only accessible via shared_secret atm. No refresh needed as no modifications have been done for 4.0.0.
	// _, err = c.GetAuthnModuleName()
	// warn(err)
	_, err = c.TokenCheck("my_fake_token") // TestTokenCheck uses "my_fake_token"
	warn(err)
	_, err = c.GetUserFromToken(c.Token)
	warn(err)

	// ── VERSION ───────────────────────────────────────────────────────────
	fmt.Println("=== version ===")
	_, err = c.Version()
	warn(err)

	// ── CLIENT MISC ───────────────────────────────────────────────────────
	fmt.Println("=== client misc ===")
	_, _, err = c.GetTftpFile("/pxelinux.cfg/default", 0, 0)
	warn(err)
	_, err = c.LastModifiedTime()
	warn(err)
	_, err = c.Ping()
	warn(err)
	_, err = c.AutoAddRepos()
	warn(err)
	_, err = c.IsAutoinstallInUse("built-in-default.ks")
	warn(err)
	err = c.GenerateIPxe("", "", "")
	warn(err)
	_, err = c.RunInstallTriggers("", "", "", "")
	warn(err)
	_, err = c.GetReposCompatibleWithProfile("testprof")
	warn(err)
	_, err = c.FindSystemByDnsName("testname")
	warn(err)
	_, err = c.GetRandomMac()
	warn(err)
	_, err = c.GetStatus(cobbler.StatusNormal)
	warn(err)
	_, err = c.GetStatus(cobbler.StatusText)
	warn(err)
	err = c.SyncDhcp()
	warn(err)

	// ── ACTIONS ───────────────────────────────────────────────────────────
	fmt.Println("=== actions ===")
	err = c.Sync()
	warn(err)
	// TestBackgroundSync uses {false, false, false} — not {true, true, true}
	_, err = c.BackgroundSync(cobbler.BackgroundSyncOptions{Dhcp: false, Dns: false, Verbose: false})
	warn(err)
	_, err = c.BackgroundSyncSystems(cobbler.BackgroundSyncSystemsOptions{Systems: []string{"", ""}, Verbose: false})
	warn(err)
	_, err = c.Check()
	warn(err)
	_, err = c.BackgroundBuildiso(cobbler.BuildisoOptions{})
	warn(err)
	_, err = c.BackgroundHardlink()
	warn(err)
	_, err = c.BackgroundValidateAutoinstallFiles()
	warn(err)
	_, err = c.BackgroundReplicate(cobbler.ReplicateOptions{})
	warn(err)
	_, err = c.BackgroundAclSetup(cobbler.AclSetupOptions{AddUser: "testing"})
	warn(err)
	_, err = c.BackgroundImport(cobbler.BackgroundImportOptions{})
	warn(err)
	_, err = c.BackgroundReposync(cobbler.BackgroundReposyncOptions{})
	warn(err)
	_, err = c.BackgroundMkLoaders()
	warn(err)
	_, err = c.BackgroundPowerSystem(cobbler.BackgroundPowerSystemOptions{Systems: []string{"testsys1"}, Power: "off"})
	warn(err)

	// ── SIGNATURES ────────────────────────────────────────────────────────
	fmt.Println("=== signatures ===")
	_, err = c.GetSignatures()
	warn(err)
	_, err = c.GetValidBreeds()
	warn(err)
	_, err = c.GetValidOsVersionsForBreed("redhat")
	warn(err)
	_, err = c.GetValidOsVersions()
	warn(err)
	_, err = c.GetValidArchs()
	warn(err)
	_, err = c.BackgroundSignatureUpdate()
	warn(err)
	_, err = c.BackgroundSignatureReload()
	warn(err)

	// ── DISTROS ───────────────────────────────────────────────────────────
	fmt.Println("=== distro create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"create-distro-name-check",
				"new-distro",
				"new-distro-modify-name",
				"new-distro-modify-comment",
				"new-distro-modify-kernel-options",
				"new-distro-modify-kernel-options-post",
				"new-distro-modify-autoinstall-meta",
				"new-distro-modify-template-files",
				"new-distro-modify-owners",
				"new-distro-modify-arch",
				"new-distro-modify-boot-loaders",
				"new-distro-modify-breed",
				"new-distro-modify-initrd",
				"new-distro-modify-remote-boot-initrd",
				"new-distro-modify-kernel",
				"new-distro-modify-remote-boot-kernel",
				"new-distro-modify-redhat-management-key",
				"new-distro-modify-redhat-management-org",
				"new-distro-modify-redhat-management-user",
				"new-distro-modify-redhat-management-password",
				"new-distro-modify-os-version",
				"new-distro-save",
				"new-distro-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		d := cobbler.NewDistro()
		d.Name = "Ubuntu-20.04-x86_64"
		d.Kernel = fakeKernel
		d.Initrd = fakeInitrd
		_, err = sub.CreateDistro(d)
		warn(err)
	}

	fmt.Println("=== distro update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-distro-handle",
				"update-distro-modify-name",
				"update-distro-modify-comment",
				"update-distro-modify-kernel-options",
				"update-distro-modify-kernel-options-post",
				"update-distro-modify-autoinstall-meta",
				"update-distro-modify-template-files",
				"update-distro-modify-owners",
				"update-distro-modify-arch",
				"update-distro-modify-boot-loaders",
				"update-distro-modify-breed",
				"update-distro-modify-initrd",
				"update-distro-modify-remote-boot-initrd",
				"update-distro-modify-kernel",
				"update-distro-modify-remote-boot-kernel",
				"update-distro-modify-redhat-management-key",
				"update-distro-modify-redhat-management-org",
				"update-distro-modify-redhat-management-user",
				"update-distro-modify-redhat-management-password",
				"update-distro-modify-os-version",
				"update-distro-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		d := cobbler.NewDistro()
		d.Name = "Ubuntu-20.04-x86_64"
		d.Kernel = fakeKernel
		d.Initrd = fakeInitrd
		warn(sub.UpdateDistro(&d))
	}

	fmt.Println("=== distros ===")
	distroHandle, err := c.GetDistroHandle("test")
	warn(err)
	err = c.SaveDistro(distroHandle, true, true, "bypass")
	warn(err)
	err = c.CopyDistro(distroHandle, "test2")
	warn(err)
	// rename_distro resolves object_id via a uid lookup, so the copy's real handle must be fetched first.
	// silentLookup is reused for the same purpose by every other Rename* call below.
	silentLookup := newSilentClient(c.Token, c.CachedVersion)
	test2Handle, err := silentLookup.GetDistroHandle("test2")
	warn(err)
	err = c.RenameDistro(test2Handle, "test1")
	warn(err)
	_, err = c.GetDistros()
	warn(err)
	_, err = c.GetDistro("Ubuntu-20.04-x86_64", false, false)
	warn(err)
	_, err = c.ListDistroNames()
	warn(err)
	_, err = c.GetDistrosSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetValidDistroBootLoaders("Ubuntu-20.04-x86_64")
	warn(err)
	_, err = c.GetDistroAsRendered("Ubuntu-20.04-x86_64")
	warn(err)
	_, err = c.FindDistro(map[string]interface{}{"name": "test"})
	warn(err)
	_, err = c.FindDistroNames(map[string]interface{}{"name": "test"})
	warn(err)

	// ── PROFILES ──────────────────────────────────────────────────────────
	fmt.Println("=== profile create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"create-profile-name-check",
				"new-profile",
				"new-profile-modify-name",
				"new-profile-modify-parent",
				"new-profile-modify-distro",
				"new-profile-modify-comment",
				"new-profile-modify-kernel-options",
				"new-profile-modify-kernel-options-post",
				"new-profile-modify-autoinstall-meta",
				"new-profile-modify-template-files",
				"new-profile-modify-owners",
				"new-profile-modify-autoinstall",
				"new-profile-modify-boot-loaders",
				"new-profile-modify-dhcp-tag",
				"new-profile-modify-name-servers",
				"new-profile-modify-name-servers-search",
				"new-profile-modify-enable-ipxe",
				"new-profile-modify-enable-menu",
				"new-profile-modify-filename",
				"new-profile-modify-menu",
				"new-profile-modify-proxy",
				"new-profile-modify-redhat-management-key",
				"new-profile-modify-redhat-management-org",
				"new-profile-modify-redhat-management-user",
				"new-profile-modify-redhat-management-password",
				"new-profile-modify-repos",
				"new-profile-modify-server",
				"new-profile-modify-next-server-v4",
				"new-profile-modify-next-server-v6",
				"new-profile-modify-virt-auto-boot",
				"new-profile-modify-virt-cpus",
				"new-profile-modify-virt-disk-driver",
				"new-profile-modify-virt-file-size",
				"new-profile-modify-virt-path",
				"new-profile-modify-virt-pxe-boot",
				"new-profile-modify-virt-ram",
				"new-profile-modify-virt-type",
				"new-profile-modify-virt-bridge",
				"new-profile-save",
				"new-profile-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		p := cobbler.NewProfile()
		p.Name = "Ubuntu-20.04-x86_64"
		ubuntuDistroUID, err := silentLookup.GetDistroHandle("Ubuntu-20.04-x86_64")
		warn(err)
		p.Distro = ubuntuDistroUID
		p.Virt.Path = fakeVirtPath
		p.Virt.Cpus = cobbler.Value[int]{Data: fakeVirtCpus}
		_, err = sub.CreateProfile(p)
		warn(err)
	}

	fmt.Println("=== profile update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-profile-handle",
				"update-profile-modify-name",
				"update-profile-modify-parent",
				"update-profile-modify-distro",
				"update-profile-modify-comment",
				"update-profile-modify-kernel-options",
				"update-profile-modify-kernel-options-post",
				"update-profile-modify-autoinstall-meta",
				"update-profile-modify-template-files",
				"update-profile-modify-owners",
				"update-profile-modify-autoinstall",
				"update-profile-modify-boot-loaders",
				"update-profile-modify-dhcp-tag",
				"update-profile-modify-name-servers",
				"update-profile-modify-name-servers-search",
				"update-profile-modify-enable-ipxe",
				"update-profile-modify-enable-menu",
				"update-profile-modify-filename",
				"update-profile-modify-menu",
				"update-profile-modify-proxy",
				"update-profile-modify-redhat-management-key",
				"update-profile-modify-redhat-management-org",
				"update-profile-modify-redhat-management-user",
				"update-profile-modify-redhat-management-password",
				"update-profile-modify-repos",
				"update-profile-modify-server",
				"update-profile-modify-next-server-v4",
				"update-profile-modify-next-server-v6",
				"update-profile-modify-virt-auto-boot",
				"update-profile-modify-virt-cpus",
				"update-profile-modify-virt-disk-driver",
				"update-profile-modify-virt-file-size",
				"update-profile-modify-virt-path",
				"update-profile-modify-virt-pxe-boot",
				"update-profile-modify-virt-ram",
				"update-profile-modify-virt-type",
				"update-profile-modify-virt-bridge",
				"update-profile-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		p := cobbler.NewProfile()
		p.Name = "Ubuntu-20.04-x86_64"
		ubuntuDistroUID, err := silentLookup.GetDistroHandle("Ubuntu-20.04-x86_64")
		warn(err)
		p.Distro = ubuntuDistroUID
		p.Virt.Path = fakeVirtPath
		p.Virt.Cpus = cobbler.Value[int]{Data: fakeVirtCpus}
		warn(sub.UpdateProfile(&p))
	}

	fmt.Println("=== profiles ===")
	profileHandle, err := c.GetProfileHandle("testprof")
	warn(err)
	err = c.SaveProfile(profileHandle, true, true, "bypass")
	warn(err)
	err = c.CopyProfile(profileHandle, "testprof2")
	warn(err)
	testprof2Handle, err := silentLookup.GetProfileHandle("testprof2")
	warn(err)
	err = c.RenameProfile(testprof2Handle, "testprof1")
	warn(err)
	_, err = c.GetProfiles()
	warn(err)
	_, err = c.GetProfile("Ubuntu-20.04-x86_64", false, false)
	warn(err)
	_, err = c.ListProfileNames()
	warn(err)
	_, err = c.GetProfilesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetValidProfileBootLoaders("Ubuntu-20.04-x86_64")
	warn(err)
	_, err = c.GetProfileAsRendered("Ubuntu-20.04-x86_64")
	warn(err)
	_, err = c.NewSubprofile()
	warn(err)
	_, err = c.GenerateAutoinstall(profileHandle, "profile", "uid", "", "")
	warn(err)
	err = c.GenerateBootCfg("testprof", "")
	warn(err)
	err = c.GenerateScript("testprof", "", "preseed_early_default")
	warn(err)
	_, err = c.FindProfile(map[string]interface{}{"name": "test"})
	warn(err)
	_, err = c.FindProfileNames(map[string]interface{}{"name": "test"})
	warn(err)
	// register_new_system passes "profile" straight to System.profile, which — like Profile.distro — requires the
	// referenced item's real uid, not its name.
	testprofUIDForRegister, err := silentLookup.GetProfileHandle("testprof")
	warn(err)
	_, err = c.RegisterNewSystem(map[string]interface{}{
		"name":    "test-register",
		"profile": testprofUIDForRegister,
		"interfaces": map[string]interface{}{
			"default": map[string]interface{}{
				"mac_address": "AA:BB:CC:DD:EE:FF",
				"ip_address":  "192.168.1.1",
				"netmask":     "255.255.255.0",
			},
		},
	})
	warn(err)

	// ── SYSTEMS ───────────────────────────────────────────────────────────
	fmt.Println("=== system create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"create-system-name-check",
				"new-system",
				"set-system-profile",
				"new-system-modify-image",
				"set-system-name",
				"new-system-modify-comment",
				"new-system-modify-kernel-options",
				"new-system-modify-kernel-options-post",
				"new-system-modify-autoinstall-meta",
				"new-system-modify-template-files",
				"new-system-modify-owners",
				"new-system-modify-autoinstall",
				"new-system-modify-boot-loaders",
				"set-system-nameservers",
				"new-system-modify-name-servers-search",
				"new-system-modify-enable-ipxe",
				"new-system-modify-filename",
				"new-system-modify-gateway",
				"set-system-hostname",
				"new-system-modify-ipv6-default-device",
				"new-system-modify-netboot-enabled",
				"new-system-modify-power-type",
				"new-system-modify-power-identity-file",
				"new-system-modify-power-options",
				"new-system-modify-power-user",
				"new-system-modify-power-password",
				"new-system-modify-power-address",
				"new-system-modify-power-id",
				"new-system-modify-proxy",
				"new-system-modify-redhat-management-key",
				"new-system-modify-redhat-management-org",
				"new-system-modify-redhat-management-user",
				"new-system-modify-redhat-management-password",
				"new-system-modify-serial-baud-rate",
				"new-system-modify-serial-device",
				"new-system-modify-server",
				"new-system-modify-status",
				"new-system-modify-next-server-v4",
				"new-system-modify-next-server-v6",
				"new-system-modify-virt-auto-boot",
				"new-system-modify-virt-cpus",
				"new-system-modify-virt-disk-driver",
				"new-system-modify-virt-file-size",
				"new-system-modify-virt-path",
				"new-system-modify-virt-pxe-boot",
				"new-system-modify-virt-ram",
				"new-system-modify-virt-type",
				"new-system-modify-virt-pxe-boot-flag",
				"new-system-save",
				"new-system-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		sys := cobbler.NewSystem()
		sys.Name = "mytestsystem"
		sys.Hostname = "blahhost"
		sys.DNS.NameServers.Data = []string{"8.8.8.8", "8.8.4.4"}
		centosProfileUID, err := silentLookup.GetProfileHandle("centos7-x86_64")
		warn(err)
		sys.Profile = centosProfileUID
		_, err = sub.CreateSystem(sys)
		warn(err)
	}

	fmt.Println("=== system update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-system-handle",
				"update-system-modify-profile",
				"update-system-modify-image",
				"update-system-modify-name",
				"update-system-modify-comment",
				"update-system-modify-kernel-options",
				"update-system-modify-kernel-options-post",
				"update-system-modify-autoinstall-meta",
				"update-system-modify-template-files",
				"update-system-modify-owners",
				"update-system-modify-autoinstall",
				"update-system-modify-boot-loaders",
				"update-system-modify-name-servers",
				"update-system-modify-name-servers-search",
				"update-system-modify-enable-ipxe",
				"update-system-modify-filename",
				"update-system-modify-gateway",
				"update-system-modify-hostname",
				"update-system-modify-ipv6-default-device",
				"update-system-modify-netboot-enabled",
				"update-system-modify-power-type",
				"update-system-modify-power-identity-file",
				"update-system-modify-power-options",
				"update-system-modify-power-user",
				"update-system-modify-power-password",
				"update-system-modify-power-address",
				"update-system-modify-power-id",
				"update-system-modify-proxy",
				"update-system-modify-redhat-management-key",
				"update-system-modify-redhat-management-org",
				"update-system-modify-redhat-management-user",
				"update-system-modify-redhat-management-password",
				"update-system-modify-serial-baud-rate",
				"update-system-modify-serial-device",
				"update-system-modify-server",
				"update-system-modify-status",
				"update-system-modify-next-server-v4",
				"update-system-modify-next-server-v6",
				"update-system-modify-virt-auto-boot",
				"update-system-modify-virt-cpus",
				"update-system-modify-virt-disk-driver",
				"update-system-modify-virt-file-size",
				"update-system-modify-virt-path",
				"update-system-modify-virt-pxe-boot",
				"update-system-modify-virt-ram",
				"update-system-modify-virt-type",
				"update-system-modify-virt-pxe-boot-flag",
			},
			rec, c.Token, c.CachedVersion,
		)
		sys := cobbler.NewSystem()
		sys.Name = "mytestsystem"
		sys.Hostname = "blahhost"
		sys.DNS.NameServers.Data = []string{"8.8.8.8", "8.8.4.4"}
		centosProfileUID, err := silentLookup.GetProfileHandle("centos7-x86_64")
		warn(err)
		sys.Profile = centosProfileUID
		sys.Power.Type = "ipmilanplus"
		sys.Status = "production"
		warn(sub.UpdateSystem(&sys))
	}

	fmt.Println("=== systems ===")
	systemHandle, err := c.GetSystemHandle("testsys")
	warn(err)
	err = c.SaveSystem(systemHandle, true, true, "bypass")
	warn(err)
	err = c.CopySystem(systemHandle, "testsys2")
	warn(err)
	testsys2Handle, err := silentLookup.GetSystemHandle("testsys2")
	warn(err)
	err = c.RenameSystem(testsys2Handle, "testsys1")
	warn(err)
	_, err = c.GetSystems()
	warn(err)
	_, err = c.GetSystem("test", false, false)
	warn(err)
	_, err = c.ListSystemNames()
	warn(err)
	_, err = c.GetSystemsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetValidSystemBootLoaders("test")
	warn(err)
	_, err = c.GetSystemAsRendered("test")
	warn(err)
	_, err = c.GetSystem("testsys", false, false) // get-interfaces-get-system
	warn(err)
	err = c.DisableNetboot("testsys")
	warn(err)
	_, err = c.UploadLogData("testsys", "/var/log/cobbler/testsys.log", 12, 0, "hello world!")
	warn(err)
	_, err = c.ClearSystemLogs(systemHandle)
	warn(err)
	_, err = c.FindSystem(map[string]interface{}{"name": "test"})
	warn(err)
	_, err = c.FindSystem(map[string]interface{}{"name": "test"}) // find-system-names fixture
	warn(err)
	_, err = c.DumpVars(systemHandle, false, true)
	warn(err)
	_, err = c.GetBlendedData("testprof", "")
	warn(err)
	err = c.GetConfigData("testsys")
	warn(err)
	testsys1Handle, err := silentLookup.GetSystemHandle("testsys1")
	warn(err)
	_, err = c.PowerSystem(testsys1Handle, "status")
	warn(err)
	err = c.DeleteSystem("test")
	warn(err)
	err = c.DeleteProfile("test")
	warn(err)
	err = c.DeleteDistro("test")
	warn(err)

	// ── TEMPLATES ─────────────────────────────────────────────────────────
	fmt.Println("=== template create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"new-template",
				"new-template-modify-name",
				"new-template-modify-comment",
				"new-template-modify-kernel-options",
				"new-template-modify-kernel-options-post",
				"new-template-modify-autoinstall-meta",
				"new-template-modify-template-files",
				"new-template-modify-owners",
				"new-template-modify-template-type",
				"new-template-modify-uri-schema",
				"new-template-modify-uri-path",
				"new-template-modify-tags",
				"new-template-modify-content",
				"new-template-save",
				"new-template-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		tpl := cobbler.NewTemplate()
		tpl.Name = "testtemplate"
		// content's setter requires a real, non-empty URI path to write to (cobbler/items/template.py rejects an
		// empty uri.path before even looking at the content value).
		// Must be a path relative to settings.autoinstall_templates_dir, pointing at a file that already exists on
		// the server (cobbler/validate.py's validate_autoinstall_template_file_path requires Path.is_file() to be
		// true before the uri.path setter allows the value at all — content.setter only writes to an existing
		// path, it never creates one). Pre-create this file in the cobbler-dev container before running this tool:
		//   docker exec cobbler-dev touch /var/lib/cobbler/templates/testtemplate.template
		tpl.URI.Path = "testtemplate.template"
		_, err = sub.CreateTemplate(tpl)
		warn(err)
	}

	fmt.Println("=== template update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-template-handle",
				"update-template-modify-name",
				"update-template-modify-comment",
				"update-template-modify-kernel-options",
				"update-template-modify-kernel-options-post",
				"update-template-modify-autoinstall-meta",
				"update-template-modify-template-files",
				"update-template-modify-owners",
				"update-template-modify-template-type",
				"update-template-modify-uri-schema",
				"update-template-modify-uri-path",
				"update-template-modify-tags",
				"update-template-modify-content",
				"update-template-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		tpl := cobbler.NewTemplate()
		tpl.Name = "testtemplate"
		// Must be a path relative to settings.autoinstall_templates_dir, pointing at a file that already exists on
		// the server (cobbler/validate.py's validate_autoinstall_template_file_path requires Path.is_file() to be
		// true before the uri.path setter allows the value at all — content.setter only writes to an existing
		// path, it never creates one). Pre-create this file in the cobbler-dev container before running this tool:
		//   docker exec cobbler-dev touch /var/lib/cobbler/templates/testtemplate.template
		tpl.URI.Path = "testtemplate.template"
		warn(sub.UpdateTemplate(&tpl))
	}

	fmt.Println("=== templates ===")
	tplHandle, err := c.GetTemplateHandle("testtemplate")
	warn(err)
	err = c.SaveTemplate(tplHandle, true, true, "bypass")
	warn(err)
	err = c.CopyTemplate(tplHandle, "testtemplate-copy")
	warn(err)
	// Renames the copy (not tplHandle's original) so "testtemplate" survives for the Read/Find calls below.
	testtemplateCopyHandle, err := silentLookup.GetTemplateHandle("testtemplate-copy")
	warn(err)
	err = c.RenameTemplate(testtemplateCopyHandle, "testtemplate-new")
	warn(err)
	_, err = c.GetTemplateFileForProfile("testprof", "/etc/motd")
	warn(err)
	_, err = c.GetTemplateFileForSystem("testsys", "/etc/motd")
	warn(err)
	_, err = c.GetTemplates()
	warn(err)
	_, err = c.GetTemplate("testtemplate", false, false)
	warn(err)
	_, err = c.GetTemplatesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.ListTemplateNames()
	warn(err)
	_, err = c.GetTemplateContent(tplHandle)
	warn(err)
	err = c.TemplatesRefreshContent(nil)
	warn(err)
	_, err = c.BackgroundTemplatesRefreshContent(nil)
	warn(err)
	_, err = c.FindTemplate(map[string]interface{}{"name": "testtemplate"})
	warn(err)
	_, err = c.FindTemplateNames(map[string]interface{}{"name": "testtemplate"})
	warn(err)
	err = c.DeleteTemplate("testtemplate")
	warn(err)

	// ── REPOS ─────────────────────────────────────────────────────────────
	fmt.Println("=== repo create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"create-repo-name-check",
				"new-repo",
				"new-repo-modify-name",
				"new-repo-modify-comment",
				"new-repo-modify-kernel-options",
				"new-repo-modify-kernel-options-post",
				"new-repo-modify-autoinstall-meta",
				"new-repo-modify-template-files",
				"new-repo-modify-owners",
				"new-repo-modify-apt-components",
				"new-repo-modify-apt-dists",
				"new-repo-modify-arch",
				"new-repo-modify-breed",
				"new-repo-modify-createrepo-flags",
				"new-repo-modify-environment",
				"new-repo-modify-keep-updated",
				"new-repo-modify-mirror",
				"new-repo-modify-mirror-locally",
				"new-repo-modify-mirror-type",
				"new-repo-modify-priority",
				"new-repo-modify-proxy",
				"new-repo-modify-rsyncopts",
				"new-repo-modify-rpm-list",
				"new-repo-modify-yumopts",
				"new-repo-save",
				"new-repo-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		r := cobbler.NewRepo()
		r.Name = "rhel-7-x86_64"
		_, err = sub.CreateRepo(r)
		warn(err)
	}

	fmt.Println("=== repo update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-repo-handle",
				"update-repo-modify-name",
				"update-repo-modify-comment",
				"update-repo-modify-kernel-options",
				"update-repo-modify-kernel-options-post",
				"update-repo-modify-autoinstall-meta",
				"update-repo-modify-template-files",
				"update-repo-modify-owners",
				"update-repo-modify-apt-components",
				"update-repo-modify-apt-dists",
				"update-repo-modify-arch",
				"update-repo-modify-breed",
				"update-repo-modify-createrepo-flags",
				"update-repo-modify-environment",
				"update-repo-modify-keep-updated",
				"update-repo-modify-mirror",
				"update-repo-modify-mirror-locally",
				"update-repo-modify-mirror-type",
				"update-repo-modify-priority",
				"update-repo-modify-proxy",
				"update-repo-modify-rsyncopts",
				"update-repo-modify-rpm-list",
				"update-repo-modify-yumopts",
				"update-repo-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		r := cobbler.NewRepo()
		r.Name = "rhel-7-x86_64"
		warn(sub.UpdateRepo(&r))
	}

	fmt.Println("=== repos ===")
	repoHandle, err := c.GetRepoHandle("testrepo")
	warn(err)
	err = c.SaveRepo(repoHandle, true, true, "bypass")
	warn(err)
	err = c.CopyRepo(repoHandle, "testrepo2")
	warn(err)
	testrepo2Handle, err := silentLookup.GetRepoHandle("testrepo2")
	warn(err)
	err = c.RenameRepo(testrepo2Handle, "testrepo1")
	warn(err)
	_, err = c.GetRepos()
	warn(err)
	_, err = c.GetRepo("rhel-7-x86_64", false, false) // get-repo fixture
	warn(err)
	_, err = c.GetRepo("rhel-7-x86_64", false, false) // get-repo-explicit-string fixture
	warn(err)
	_, err = c.ListRepoNames()
	warn(err)
	_, err = c.GetReposSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetRepoConfigForProfile("testprof")
	warn(err)
	_, err = c.GetRepoConfigForSystem("testsys")
	warn(err)
	_, err = c.GetRepoAsRendered("rhel-7-x86_64")
	warn(err)
	_, err = c.FindRepo(map[string]interface{}{"name": "test"})
	warn(err)
	_, err = c.FindRepoNames(map[string]interface{}{"name": "test"})
	warn(err)
	err = c.DeleteRepo("test")
	warn(err)

	// ── IMAGES ────────────────────────────────────────────────────────────
	fmt.Println("=== image create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"new-image",
				"new-image-modify-name",
				"new-image-modify-comment",
				"new-image-modify-kernel-options",
				"new-image-modify-kernel-options-post",
				"new-image-modify-autoinstall-meta",
				"new-image-modify-template-files",
				"new-image-modify-owners",
				"new-image-modify-arch",
				"new-image-modify-autoinstall",
				"new-image-modify-breed",
				"new-image-modify-file",
				"new-image-modify-image-type",
				"new-image-modify-network-count",
				"new-image-modify-os-version",
				"new-image-modify-boot-loaders",
				"new-image-modify-menu",
				"new-image-modify-virt-auto-boot",
				"new-image-modify-virt-cpus",
				"new-image-modify-virt-disk-driver",
				"new-image-modify-virt-file-size",
				"new-image-modify-virt-path",
				"new-image-modify-virt-pxe-boot",
				"new-image-modify-virt-ram",
				"new-image-modify-virt-type",
				"new-image-modify-virt-bridge",
				"new-image-save",
				"new-image-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		image := cobbler.NewImage()
		image.Name = "testimage"
		image.Virt.Path = fakeVirtPath
		image.Virt.Cpus = cobbler.Value[int]{Data: fakeVirtCpus}
		_, err = sub.CreateImage(image)
		warn(err)
	}

	fmt.Println("=== image update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-image-handle",
				"update-image-modify-name",
				"update-image-modify-comment",
				"update-image-modify-kernel-options",
				"update-image-modify-kernel-options-post",
				"update-image-modify-autoinstall-meta",
				"update-image-modify-template-files",
				"update-image-modify-owners",
				"update-image-modify-arch",
				"update-image-modify-autoinstall",
				"update-image-modify-breed",
				"update-image-modify-file",
				"update-image-modify-image-type",
				"update-image-modify-network-count",
				"update-image-modify-os-version",
				"update-image-modify-boot-loaders",
				"update-image-modify-menu",
				"update-image-modify-virt-auto-boot",
				"update-image-modify-virt-cpus",
				"update-image-modify-virt-disk-driver",
				"update-image-modify-virt-file-size",
				"update-image-modify-virt-path",
				"update-image-modify-virt-pxe-boot",
				"update-image-modify-virt-ram",
				"update-image-modify-virt-type",
				"update-image-modify-virt-bridge",
				"update-image-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		image := cobbler.NewImage()
		image.Name = "testimage"
		image.Virt.Path = fakeVirtPath
		image.Virt.Cpus = cobbler.Value[int]{Data: fakeVirtCpus}
		warn(sub.UpdateImage(&image))
	}

	fmt.Println("=== images ===")
	imageHandle, err := c.GetImageHandle("testimage")
	warn(err)
	err = c.SaveImage(imageHandle, true, true, "bypass")
	warn(err)
	err = c.CopyImage(imageHandle, "testimage2")
	warn(err)
	testimage2Handle, err := silentLookup.GetImageHandle("testimage2")
	warn(err)
	err = c.RenameImage(testimage2Handle, "testimage1")
	warn(err)
	_, err = c.GetImages()
	warn(err)
	_, err = c.GetImage("testimage", false, false)
	warn(err)
	_, err = c.ListImageNames()
	warn(err)
	_, err = c.GetImagesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.FindImage(map[string]interface{}{"name": "testimage"})
	warn(err)
	_, err = c.FindImageNames(map[string]interface{}{"name": "testimage"})
	warn(err)
	err = c.DeleteImage("test")
	warn(err)

	// ── MENUS ─────────────────────────────────────────────────────────────
	fmt.Println("=== menu create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"create-menu-name-check",
				"new-menu",
				"new-menu-modify-name",
				"new-menu-modify-comment",
				"new-menu-modify-kernel-options",
				"new-menu-modify-kernel-options-post",
				"new-menu-modify-autoinstall-meta",
				"new-menu-modify-template-files",
				"new-menu-modify-owners",
				"new-menu-modify-display-name",
				"new-menu-save",
				"new-menu-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		menu := cobbler.NewMenu()
		menu.Name = "grub-menu"
		_, err = sub.CreateMenu(menu)
		warn(err)
	}

	fmt.Println("=== menu update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-menu-handle",
				"update-menu-modify-name",
				"update-menu-modify-comment",
				"update-menu-modify-kernel-options",
				"update-menu-modify-kernel-options-post",
				"update-menu-modify-autoinstall-meta",
				"update-menu-modify-template-files",
				"update-menu-modify-owners",
				"update-menu-modify-display-name",
				"update-menu-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		menu := cobbler.NewMenu()
		menu.Name = "grub-menu"
		warn(sub.UpdateMenu(&menu))
	}

	fmt.Println("=== menus ===")
	menuHandle, err := c.GetMenuHandle("testmenu")
	warn(err)
	err = c.SaveMenu(menuHandle, true, true, "bypass")
	warn(err)
	err = c.CopyMenu(menuHandle, "testmenu2")
	warn(err)
	testmenu2Handle, err := silentLookup.GetMenuHandle("testmenu2")
	warn(err)
	err = c.RenameMenu(testmenu2Handle, "testmenu1")
	warn(err)
	_, err = c.GetMenus()
	warn(err)
	_, err = c.GetMenu("testmenu", false, false)
	warn(err)
	_, err = c.ListMenuNames()
	warn(err)
	_, err = c.GetMenusSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetMenuAsRendered("testmenu")
	warn(err)
	_, err = c.FindMenu(map[string]interface{}{"name": "testmenu"})
	warn(err)
	_, err = c.FindMenuNames(map[string]interface{}{"name": "testmenu"})
	warn(err)
	err = c.DeleteMenu("test")
	warn(err)

	// ── NETWORK INTERFACES ────────────────────────────────────────────────
	fmt.Println("=== network interface create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"new-network-interface",
				"new-network-interface-modify-name",
				"new-network-interface-modify-comment",
				"new-network-interface-modify-kernel-options",
				"new-network-interface-modify-kernel-options-post",
				"new-network-interface-modify-autoinstall-meta",
				"new-network-interface-modify-template-files",
				"new-network-interface-modify-owners",
				"new-network-interface-modify-mac-address",
				"new-network-interface-modify-interface-type",
				"new-network-interface-modify-interface-master",
				"new-network-interface-modify-bonding-opts",
				"new-network-interface-modify-bridge-opts",
				"new-network-interface-modify-connected-mode",
				"new-network-interface-modify-management",
				"new-network-interface-modify-static",
				"new-network-interface-modify-dhcp-tag",
				"new-network-interface-modify-if-gateway",
				"new-network-interface-modify-mtu",
				"new-network-interface-modify-ipv6-default-gateway",
				"new-network-interface-modify-ipv6-static-routes",
				"new-network-interface-modify-virt-bridge",
				"new-network-interface-modify-ipv4-address",
				"new-network-interface-modify-ipv4-netmask",
				"new-network-interface-modify-ipv4-static-routes",
				"new-network-interface-modify-ipv6-address",
				"new-network-interface-modify-ipv6-prefix",
				"new-network-interface-modify-ipv6-secondaries",
				"new-network-interface-modify-ipv6-mtu",
				"new-network-interface-modify-ipv6-option-static-routes",
				"new-network-interface-modify-dns-name",
				"new-network-interface-modify-dns-cnames",
				"new-network-interface-save",
				"new-network-interface-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		silentLookup := newSilentClient(c.Token, c.CachedVersion)
		sysHandle, err := silentLookup.GetSystemHandle("mytestsystem")
		warn(err)
		ni := cobbler.NewNetworkInterface()
		ni.Name = "eth0-mytestsystem"
		ni.MacAddress = "AA:BB:CC:DD:EE:00"
		_, err = sub.CreateNetworkInterface(sysHandle, ni)
		warn(err)
	}

	fmt.Println("=== network interface update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-network-interface-handle",
				"update-network-interface-modify-name",
				"update-network-interface-modify-comment",
				"update-network-interface-modify-kernel-options",
				"update-network-interface-modify-kernel-options-post",
				"update-network-interface-modify-autoinstall-meta",
				"update-network-interface-modify-template-files",
				"update-network-interface-modify-owners",
				"update-network-interface-modify-mac-address",
				"update-network-interface-modify-interface-type",
				"update-network-interface-modify-interface-master",
				"update-network-interface-modify-bonding-opts",
				"update-network-interface-modify-bridge-opts",
				"update-network-interface-modify-connected-mode",
				"update-network-interface-modify-management",
				"update-network-interface-modify-static",
				"update-network-interface-modify-dhcp-tag",
				"update-network-interface-modify-if-gateway",
				"update-network-interface-modify-mtu",
				"update-network-interface-modify-ipv6-default-gateway",
				"update-network-interface-modify-ipv6-static-routes",
				"update-network-interface-modify-virt-bridge",
				"update-network-interface-modify-ipv4-address",
				"update-network-interface-modify-ipv4-netmask",
				"update-network-interface-modify-ipv4-static-routes",
				"update-network-interface-modify-ipv6-address",
				"update-network-interface-modify-ipv6-prefix",
				"update-network-interface-modify-ipv6-secondaries",
				"update-network-interface-modify-ipv6-mtu",
				"update-network-interface-modify-ipv6-option-static-routes",
				"update-network-interface-modify-dns-name",
				"update-network-interface-modify-dns-cnames",
				"update-network-interface-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		ni := cobbler.NewNetworkInterface()
		ni.Name = "eth0-mytestsystem"
		// Must match the case cobbler/validate.py's mac_address() normalises to (lower) exactly: the mac_address
		// setter's "unchanged" short-circuit (`if self._mac_address == address: return`) compares the raw
		// argument against the already-normalised stored value *before* normalising it, so a same-value update
		// in a different case falls through to the duplicate-MAC check and conflicts with itself.
		ni.MacAddress = "aa:bb:cc:dd:ee:00"
		warn(sub.UpdateNetworkInterface(&ni))
	}

	fmt.Println("=== network interfaces ===")
	niHandle, err := c.GetNetworkInterfaceHandle("eth0-server1")
	warn(err)
	err = c.SaveNetworkInterface(niHandle, true, true, "bypass")
	warn(err)
	err = c.CopyNetworkInterface(niHandle, "eth1-server1")
	warn(err)
	// Renames the copy (not niHandle's original) so "eth0-server1" survives for the Read/Find calls below.
	// rename_item resolves object_id via a uid lookup, so the copy's real handle must be fetched first.
	eth1Handle, err := silentLookup.GetNetworkInterfaceHandle("eth1-server1")
	warn(err)
	err = c.RenameNetworkInterface(eth1Handle, "eth2-server1")
	warn(err)
	_, err = c.GetNetworkInterface("eth0-server1", false, false)
	warn(err)
	_, err = c.GetNetworkInterfaces()
	warn(err)
	_, err = c.ListNetworkInterfaceNames()
	warn(err)
	_, err = c.GetNetworkInterfacesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.FindNetworkInterface(map[string]interface{}{"mac_address": "00:11:22:33:44:55"})
	warn(err)
	_, err = c.FindNetworkInterfaceNames(map[string]interface{}{"mac_address": "00:11:22:33:44:55"})
	warn(err)
	err = c.DeleteNetworkInterface("eth0-server1")
	warn(err)

	// ── DISTRO GROUP ──────────────────────────────────────────────────────
	fmt.Println("=== distro group create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"new-distro-group",
				"new-distro-group-modify-name",
				"new-distro-group-modify-comment",
				"new-distro-group-modify-kernel-options",
				"new-distro-group-modify-kernel-options-post",
				"new-distro-group-modify-autoinstall-meta",
				"new-distro-group-modify-template-files",
				"new-distro-group-modify-owners",
				"new-distro-group-modify-items",
				"new-distro-group-save",
				"new-distro-group-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		dg := cobbler.NewDistroGroup()
		dg.Name = "webservers"
		dg.Members = []string{"member-a"}
		_, err = sub.CreateDistroGroup(dg)
		warn(err)
	}

	fmt.Println("=== distro group update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-distro-group-handle",
				"update-distro-group-modify-name",
				"update-distro-group-modify-comment",
				"update-distro-group-modify-kernel-options",
				"update-distro-group-modify-kernel-options-post",
				"update-distro-group-modify-autoinstall-meta",
				"update-distro-group-modify-template-files",
				"update-distro-group-modify-owners",
				"update-distro-group-modify-items",
				"update-distro-group-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		dg := cobbler.NewDistroGroup()
		dg.Name = "webservers"
		dg.Members = []string{"member-a"}
		warn(sub.UpdateDistroGroup(&dg))
	}

	fmt.Println("=== distro groups ===")
	dgh, err := c.GetDistroGroupHandle("webservers")
	warn(err)
	err = c.SaveDistroGroup(dgh, true, true, "bypass")
	warn(err)
	err = c.CopyDistroGroup(dgh, "webservers2")
	warn(err)
	dgWebservers2Handle, err := silentLookup.GetDistroGroupHandle("webservers2")
	warn(err)
	err = c.RenameDistroGroup(dgWebservers2Handle, "webservers-new")
	warn(err)
	_, err = c.GetDistroGroup("webservers", false, false)
	warn(err)
	_, err = c.GetDistroGroups()
	warn(err)
	_, err = c.ListDistroGroupNames()
	warn(err)
	_, err = c.GetDistroGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetDistroGroupAsRendered("webservers")
	warn(err)
	_, err = c.FindDistroGroup(map[string]interface{}{"name": "webservers"})
	warn(err)
	err = c.DeleteDistroGroup("webservers")
	warn(err)

	// ── PROFILE GROUP ─────────────────────────────────────────────────────
	fmt.Println("=== profile group create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"new-profile-group",
				"new-profile-group-modify-name",
				"new-profile-group-modify-comment",
				"new-profile-group-modify-kernel-options",
				"new-profile-group-modify-kernel-options-post",
				"new-profile-group-modify-autoinstall-meta",
				"new-profile-group-modify-template-files",
				"new-profile-group-modify-owners",
				"new-profile-group-modify-items",
				"new-profile-group-save",
				"new-profile-group-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		pg := cobbler.NewProfileGroup()
		pg.Name = "webservers"
		pg.Members = []string{"member-a"}
		_, err = sub.CreateProfileGroup(pg)
		warn(err)
	}

	fmt.Println("=== profile group update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-profile-group-handle",
				"update-profile-group-modify-name",
				"update-profile-group-modify-comment",
				"update-profile-group-modify-kernel-options",
				"update-profile-group-modify-kernel-options-post",
				"update-profile-group-modify-autoinstall-meta",
				"update-profile-group-modify-template-files",
				"update-profile-group-modify-owners",
				"update-profile-group-modify-items",
				"update-profile-group-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		pg := cobbler.NewProfileGroup()
		pg.Name = "webservers"
		pg.Members = []string{"member-a"}
		warn(sub.UpdateProfileGroup(&pg))
	}

	fmt.Println("=== profile groups ===")
	pgh, err := c.GetProfileGroupHandle("webservers")
	warn(err)
	err = c.SaveProfileGroup(pgh, true, true, "bypass")
	warn(err)
	err = c.CopyProfileGroup(pgh, "webservers2")
	warn(err)
	pgWebservers2Handle, err := silentLookup.GetProfileGroupHandle("webservers2")
	warn(err)
	err = c.RenameProfileGroup(pgWebservers2Handle, "webservers-new")
	warn(err)
	_, err = c.GetProfileGroup("webservers", false, false)
	warn(err)
	_, err = c.GetProfileGroups()
	warn(err)
	_, err = c.ListProfileGroupNames()
	warn(err)
	_, err = c.GetProfileGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetProfileGroupAsRendered("webservers")
	warn(err)
	_, err = c.FindProfileGroup(map[string]interface{}{"name": "webservers"})
	warn(err)
	err = c.DeleteProfileGroup("webservers")
	warn(err)

	// ── SYSTEM GROUP ──────────────────────────────────────────────────────
	fmt.Println("=== system group create ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"new-system-group",
				"new-system-group-modify-name",
				"new-system-group-modify-comment",
				"new-system-group-modify-kernel-options",
				"new-system-group-modify-kernel-options-post",
				"new-system-group-modify-autoinstall-meta",
				"new-system-group-modify-template-files",
				"new-system-group-modify-owners",
				"new-system-group-modify-items",
				"new-system-group-save",
				"new-system-group-get",
			},
			rec, c.Token, c.CachedVersion,
		)
		sg := cobbler.NewSystemGroup()
		sg.Name = "webservers"
		sg.Members = []string{"member-a"}
		_, err = sub.CreateSystemGroup(sg)
		warn(err)
	}

	fmt.Println("=== system group update ===")
	{
		sub, _ := makeSubClient(
			[]string{
				"update-system-group-handle",
				"update-system-group-modify-name",
				"update-system-group-modify-comment",
				"update-system-group-modify-kernel-options",
				"update-system-group-modify-kernel-options-post",
				"update-system-group-modify-autoinstall-meta",
				"update-system-group-modify-template-files",
				"update-system-group-modify-owners",
				"update-system-group-modify-items",
				"update-system-group-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		sg := cobbler.NewSystemGroup()
		sg.Name = "webservers"
		sg.Members = []string{"member-a"}
		warn(sub.UpdateSystemGroup(&sg))
	}

	fmt.Println("=== system groups ===")
	sgh, err := c.GetSystemGroupHandle("webservers")
	warn(err)
	err = c.SaveSystemGroup(sgh, true, true, "bypass")
	warn(err)
	err = c.CopySystemGroup(sgh, "webservers2")
	warn(err)
	sgWebservers2Handle, err := silentLookup.GetSystemGroupHandle("webservers2")
	warn(err)
	err = c.RenameSystemGroup(sgWebservers2Handle, "webservers-new")
	warn(err)
	_, err = c.GetSystemGroup("webservers", false, false)
	warn(err)
	_, err = c.GetSystemGroups()
	warn(err)
	_, err = c.ListSystemGroupNames()
	warn(err)
	_, err = c.GetSystemGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	warn(err)
	_, err = c.GetSystemGroupAsRendered("webservers")
	warn(err)
	_, err = c.FindSystemGroup(map[string]interface{}{"name": "webservers"})
	warn(err)
	err = c.DeleteSystemGroup("webservers")
	warn(err)

	// Recreate the distro_group deleted above, silently (no fixture slot consumed), so that ModifyItemInPlace and
	// GetItemResolvedValue below have a real item UID to work with.
	var itemUUID string
	{
		silent := cobbler.NewClient(&silentHTTPClient{}, config)
		silent.Token = c.Token
		silent.CachedVersion = c.CachedVersion
		dg := cobbler.NewDistroGroup()
		dg.Name = "webservers"
		dg.Members = []string{"member-a"}
		createdDG, createErr := silent.CreateDistroGroup(dg)
		warn(createErr)
		if createdDG != nil {
			itemUUID = createdDG.Uid
		}
	}

	// ── ITEMS (generic) ───────────────────────────────────────────────────
	fmt.Println("=== items generic ===")
	_, err = c.FindItemsPaged("menu", map[string]interface{}{"display_name": ""}, "", 1, 5)
	warn(err)
	_, err = c.GetItem("system", "test", false, false)
	warn(err)
	_, err = c.GetItem("system", "testsys", true, false)
	warn(err)
	_, err = c.GetItem("system", "testsys", false, true)
	warn(err)
	_, err = c.FindItems("profile", map[string]interface{}{"name": "test*"}, "name", false)
	warn(err)
	_, err = c.FindItemNames("profile", map[string]interface{}{"name": "test*"}, "name")
	warn(err)
	// modify_item resolves object_id via a uid lookup, so the recreated "webservers" group's real uid is used here.
	err = c.ModifyItem("distro_group", itemUUID, []string{"comment"}, "hello")
	warn(err)

	// ModifyItemInPlace makes 4 HTTP calls; isolated sub-recorder prevents
	// sequence corruption if any call fails.
	{
		sub, _ := makeSubClient(
			[]string{
				"modify-item-in-place-get",
				"modify-item-in-place-handle",
				"modify-item-in-place-modify",
				"modify-item-in-place-save",
			},
			rec, c.Token, c.CachedVersion,
		)
		// distro_group (and the other group types) don't actually expose kernel_options/autoinstall_meta/etc. as
		// real properties server-side — cobbler/items/abstract/item_group.py never defines them, so modify_item
		// silently no-ops (setattr on an unknown attribute) and to_dict() never reflects it, leaving GetItem's
		// response without a "kernel_options" key at all. Use a profile instead, which has a real one.
		warn(sub.ModifyItemInPlace("profile", "testprof", "kernel_options", map[string]interface{}{"test": "1"}))
	}

	_, err = c.GetItemResolvedValue(itemUUID, []string{"kernel_options"})
	warn(err)
	err = c.SetItemResolvedValue(itemUUID, []string{"comment"}, "hello")
	warn(err)
	_, err = c.HasItem("template", "testtemplate")
	warn(err)
	err = c.NewItem("template", false)
	warn(err)

	// ── SETTINGS ──────────────────────────────────────────────────────────
	fmt.Println("=== settings ===")
	_, err = c.GetSettings()
	warn(err)
	// TestModifySettings uses ("auth_token_expiration", 7200) with int value
	_, err = c.ModifySetting("auth_token_expiration", 7200)
	warn(err)

	// Discover a real event ID so GetTaskStatus / GetEventLog fixtures reflect
	// an actual server event rather than a stale hard-coded ID. GetEvents ranges over a map, so its result order is
	// randomised by the Go runtime; sort by ID (which is timestamp-prefixed) to deterministically pick the
	// earliest-created background task every run instead of a random one.
	var bgTaskID string
	{
		silent := cobbler.NewClient(&silentHTTPClient{}, config)
		silent.Token = c.Token
		silent.CachedVersion = c.CachedVersion
		if evts, evtErr := silent.GetEvents(""); evtErr == nil && len(evts) > 0 {
			sort.Slice(evts, func(i, j int) bool { return evts[i].ID < evts[j].ID })
			bgTaskID = evts[0].ID
		}
	}

	// ── EVENTS ────────────────────────────────────────────────────────────
	fmt.Println("=== events ===")
	_, err = c.GetTaskStatus(bgTaskID)
	warn(err)
	_, err = c.GetEvents("")
	warn(err)
	_, err = c.GetEventLog(bgTaskID)
	warn(err)

	// ── TRANSACTIONS ──────────────────────────────────────────────────────
	fmt.Println("=== transactions ===")
	err = c.TransactionBegin()
	warn(err)
	err = c.TransactionCommit()
	warn(err)
	// Begin a new transaction silently (not recorded) so that the following
	// TransactionAbort call has an active transaction to abort.
	{
		silent := cobbler.NewClient(&silentHTTPClient{}, config)
		silent.Token = c.Token
		silent.CachedVersion = c.CachedVersion
		_ = silent.TransactionBegin()
	}
	err = c.TransactionAbort()
	warn(err)

	// ── INPUT HELPERS ─────────────────────────────────────────────────────
	fmt.Println("=== input helpers ===")
	_, _, err = c.InputBoolean("yes")
	warn(err)
	_, _, err = c.InputInt("42")
	warn(err)
	_, _, err = c.InputStringOrList("a,b,c")
	warn(err)
	_, _, err = c.InputStringOrListNoInherit("a,b,c")
	warn(err)
	_, _, err = c.InputStringOrDict("k1=v1 k2=v2")
	warn(err)
	_, _, err = c.InputStringOrDictNoInherit("k1=v1 k2=v2")
	warn(err)

	// ── LOGOUT ────────────────────────────────────────────────────────────
	fmt.Println("=== logout ===")
	_, err = c.Logout()
	warn(err)

	fmt.Printf("\nDone. Main recorder: %d/%d fixtures written.\n", rec.idx, len(fixtureSequence))
}
