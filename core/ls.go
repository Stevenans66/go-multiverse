package core

import (
	"context"
	"sort"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
)

// DirEntry contains info about a file.
type DirEntry struct {
	// Cid is the id of the entry.
	Cid cid.Cid
	// Name is the name of the file.
	Name string
	// IsDir indicates if the file is a directory.
	IsDir bool
	// Size is the size in bytes of the file.
	Size uint64
}

// Ls returns the contents of the directory with the given CID.
func Ls(ctx context.Context, dag ipld.DAGService, id cid.Cid) ([]*DirEntry, error) {
	node, err := dag.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	dir, err := ufsio.NewDirectoryFromNode(dag, node)
	if err != nil {
		return nil, err
	}

	links, err := dir.Links(ctx)
	if err != nil {
		return nil, err
	}

	var list []*DirEntry
	for _, l := range links {
		info, err := stat(ctx, dag, l)
		if err != nil {
			return nil, err
		}

		list = append(list, info)
	}

	// sort files using unix fs rules
	sort.Slice(list, func(i, j int) bool {
		switch {
		case !list[i].IsDir && list[j].IsDir:
			return false
		case list[i].IsDir && !list[j].IsDir:
			return true
		default:
			return list[i].Name < list[j].Name
		}
	})

	return list, nil
}

// stat returns file info for the given link.
func stat(ctx context.Context, dag ipld.DAGService, link *ipld.Link) (*DirEntry, error) {
	node, err := dag.Get(ctx, link.Cid)
	if err != nil {
		return nil, err
	}

	fsnode, err := unixfs.ExtractFSNode(node)
	if err != nil {
		return nil, err
	}

	return &DirEntry{
		Cid:   link.Cid,
		Name:  link.Name,
		IsDir: fsnode.IsDir(),
		Size:  fsnode.FileSize(),
	}, nil
}
