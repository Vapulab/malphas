package layout

import "sync"

// Asset versions let Base reference each asset as path?v=<version>. The
// stylesheet and scripts are cached hard — by browsers and by a CDN in front
// (Cloudflare keeps /static/app.css for hours) — so without a version in the
// URL a deploy that changes them stays invisible until those caches expire.
// With the content hash in the URL, a changed file is simply a new URL.
//
// The app sets the versions once at startup (SetAssetVersions); paths with no
// version are referenced as they are.
var (
	assetMu       sync.RWMutex
	assetVersions = map[string]string{}
)

// SetAssetVersions replaces the version map, path → version.
func SetAssetVersions(v map[string]string) {
	m := make(map[string]string, len(v))
	for k, val := range v {
		m[k] = val
	}
	assetMu.Lock()
	assetVersions = m
	assetMu.Unlock()
}

func asset(path string) string {
	assetMu.RLock()
	v := assetVersions[path]
	assetMu.RUnlock()
	if v == "" {
		return path
	}
	return path + "?v=" + v
}
