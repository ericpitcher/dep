package gps

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/Masterminds/vcs"
)

type maybeSource interface {
	try(cachedir string, an ProjectAnalyzer) (source, error)
}

type maybeSources []maybeSource

type maybeGitSource struct {
	n   string
	url *url.URL
}

func (m maybeGitSource) try(cachedir string, an ProjectAnalyzer) (source, error) {
	path := filepath.Join(cachedir, "sources", sanitizer.Replace(m.url.String()))
	r, err := vcs.NewGitRepo(m.url.String(), path)
	if err != nil {
		return nil, err
	}

	src := &gitSource{
		baseVCSSource: baseVCSSource{
			an: an,
			dc: newMetaCache(),
			crepo: &repo{
				r:     r,
				rpath: path,
			},
		},
	}

	_, err = src.listVersions()
	if err != nil {
		return nil, err
		//} else if pm.ex.f&existsUpstream == existsUpstream {
		//return pm, nil
	}

	return src, nil
}

type maybeBzrSource struct {
	n   string
	url *url.URL
}

func (m maybeBzrSource) try(cachedir string, an ProjectAnalyzer) (source, error) {
	path := filepath.Join(cachedir, "sources", sanitizer.Replace(m.url.String()))
	r, err := vcs.NewBzrRepo(m.url.String(), path)
	if err != nil {
		return nil, err
	}
	if !r.Ping() {
		return nil, fmt.Errorf("Remote repository at %s does not exist, or is inaccessible", m.url.String())
	}

	return &bzrSource{
		baseVCSSource: baseVCSSource{
			an: an,
			dc: newMetaCache(),
			crepo: &repo{
				r:     r,
				rpath: path,
			},
		},
	}, nil
}

type maybeHgSource struct {
	n   string
	url *url.URL
}

func (m maybeHgSource) try(cachedir string, an ProjectAnalyzer) (source, error) {
	path := filepath.Join(cachedir, "sources", sanitizer.Replace(m.url.String()))
	r, err := vcs.NewHgRepo(m.url.String(), path)
	if err != nil {
		return nil, err
	}
	if !r.Ping() {
		return nil, fmt.Errorf("Remote repository at %s does not exist, or is inaccessible", m.url.String())
	}

	return &hgSource{
		baseVCSSource: baseVCSSource{
			an: an,
			dc: newMetaCache(),
			crepo: &repo{
				r:     r,
				rpath: path,
			},
		},
	}, nil
}
