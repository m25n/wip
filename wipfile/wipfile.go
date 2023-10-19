package wipfile

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type WIPFile interface {
	AppendLine(line []byte) error
	Lines(func([]byte) error) error
}

func FromEnv() (WIPFile, error) {
	uri, err := uriFromEnvWithDefault()
	if err != nil {
		return nil, err
	}
	switch uri.Scheme {
	case "file":
		filename, err := filepath.Abs(uri.Opaque)
		if err != nil {
			return nil, err
		}
		return &file{filename: filename}, nil
	default:
		return nil, fmt.Errorf(`unknown wipfile scheme "%s"`, uri.Scheme)
	}
}

func uriFromEnvWithDefault() (*url.URL, error) {
	wipfile := os.Getenv("WIPFILE")
	if wipfile == "" {
		return defaultURI()
	}
	return url.Parse(wipfile)
}

func defaultURI() (*url.URL, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return url.Parse(fmt.Sprintf("file://%s", filepath.Join(home, ".wip")))
}
