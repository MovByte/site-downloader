package resourceAttrsMap

// Where to find external references in HTML
// attribute, selectors for it
var HTMLResourceSelectors = map[string][]string{
	"src": {"script[src]", "img[src]", "video[src]", "audio[src]", "embed[src]", "iframe[src]", "source[src]", "track[src]", "picture[srcset]", "meta[content]"},
	"href": {"a[href]", "link[rel='stylesheet'][href]"},
	"data": {"object[data]"},
}

// TODO: Export