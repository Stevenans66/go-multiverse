package core

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/storage"
)

// WalkFun is called for each commit visited by walk.
type WalkFun func(cid.Cid, *object.Commit) bool

// Walk traverses the commit history starting at the given id.
func Walk(ctx context.Context, store *storage.Store, id cid.Cid, cb WalkFun) (map[string]*object.Commit, error) {
	history := make(map[string]*object.Commit)

	// perform a depth first traversal of parent commits
	getLinks := func(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
		commit, ok := history[id.KeyString()]
		if !ok {
			return nil, errors.New("commit not found")
		}

		return commit.ParentLinks(), nil
	}

	visit := func(id cid.Cid) bool {
		if _, ok := history[id.KeyString()]; ok {
			return false
		}

		node, err := store.Dag.Get(ctx, id)
		if err != nil {
			return false
		}

		commit, err := object.CommitFromCBOR(node.RawData())
		if err != nil {
			return false
		}

		history[id.KeyString()] = commit
		if cb != nil {
			return cb(id, commit)
		}

		return true
	}

	return history, merkledag.Walk(ctx, getLinks, id, visit)
}