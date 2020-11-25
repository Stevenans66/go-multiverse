package core

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// MergeBase returns the best common ancestor of local and remote.
func MergeBase(ctx context.Context, store *storage.Store, local, remote cid.Cid) (cid.Cid, error) {
	history, err := Walk(ctx, store, local, nil)
	if err != nil {
		return cid.Cid{}, err
	}

	// local is ahead of remote
	if _, ok := history[remote.KeyString()]; ok {
		return remote, nil
	}

	var best cid.Cid
	var err0 error

	// find the least common ancestor by searching
	// for commits that are in both local and remote
	// and that are also independent from each other
	cb := func(id cid.Cid, commit *object.Commit) bool {
		if err0 != nil {
			return false
		}

		if _, ok := history[id.KeyString()]; !ok {
			return true
		}

		var match bool
		if match, err0 = isAncestor(ctx, store, best, id); !match {
			best = id
		}

		return false
	}

	if _, err := Walk(ctx, store, remote, cb); err != nil {
		return cid.Cid{}, err
	}

	return best, err0
}

// isAncestor returns true if child is an ancestor of parent.
func isAncestor(ctx context.Context, store *storage.Store, parent, child cid.Cid) (bool, error) {
	if !(parent.Defined() && child.Defined()) {
		return false, nil
	}

	cb := func(id cid.Cid, commit *object.Commit) bool {
		return id != child
	}

	history, err := Walk(ctx, store, parent, cb)
	if err != nil {
		return false, err
	}

	_, ok := history[child.KeyString()]
	return ok, nil
}