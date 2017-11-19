// Package prune provides node_modules pruning of unnecessary files.
package prune

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// Stats for a prune.
type Stats struct {
	FilesTotal   int64
	FilesRemoved int64
	SizeRemoved  int64
}

// Pruner is a module pruner.
type Pruner struct {
	Dir string
	Log log.Interface
}

// Option function.
type Option func(*Pruner)

// New with the given options.
func New(options ...Option) *Pruner {
	v := &Pruner{Dir: "node_modules", Log: log.Log}
	for _, o := range options {
		o(v)
	}
	return v
}

// WithDir option.
func WithDir(s string) Option {
	return func(v *Pruner) {
		v.Dir = s
	}
}

// Prune performs the pruning.
func (p Pruner) Prune() (*Stats, error) {
	var stats Stats

	err := filepath.Walk(p.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		stats.FilesTotal++

		if !p.prune(path, info) {
			return nil
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		p.Log.WithField("path", path).Debug("prune")
		stats.FilesRemoved++
		stats.SizeRemoved += info.Size()

		if err := os.Remove(path); err != nil {
			return errors.Wrap(err, "removing")
		}

		return nil
	})

	return &stats, err
}

// prune returns true if the file or dir should be pruned.
func (p Pruner) prune(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	ext := filepath.Ext(path)
	return ext == ".ts" || ext == ".md"
}
