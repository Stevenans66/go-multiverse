package rpc

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/ipfs/go-datastore"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/peer"
)

func TestLog(t *testing.T) {
	ctx := context.Background()

	dstore := datastore.NewMapDatastore()
	store := data.NewStore(dstore)

	mock, err := peer.Mock(ctx, dstore)
	if err != nil {
		t.Fatal("failed to create peer")
	}

	json, err := ioutil.ReadFile("testdata/repository.json")
	if err != nil {
		t.Fatal("failed to read json")
	}

	repo, err := data.RepositoryFromJSON(json)
	if err != nil {
		t.Fatal("failed to parse repository json")
	}

	id, err := data.AddRepository(ctx, mock, repo)
	if err != nil {
		t.Fatal("failed to add repository")
	}

	if err := store.PutCid(repo.Name, id); err != nil {
		t.Fatal("failed to put cid in store")
	}

	json, err = ioutil.ReadFile("testdata/commit.json")
	if err != nil {
		t.Fatal("failed to read json")
	}

	commit, err := data.CommitFromJSON(json)
	if err != nil {
		t.Fatal("failed to parse commit json")
	}

	head, err := data.AddCommit(ctx, mock, commit)
	if err != nil {
		t.Fatal("failed to add commit")
	}

	client, err := connect(mock, store)
	if err != nil {
		t.Fatal("failed to connect to rpc server")
	}

	args := LogArgs{
		Name:   repo.Name,
		Branch: "default",
		Limit:  1,
	}

	var reply LogReply
	if err := client.Call("Service.Log", &args, &reply); err != nil {
		t.Fatal("failed to call rpc")
	}

	if len(reply.IDs) != 1 {
		t.Fatal("unexpected ids")
	}

	if reply.IDs[0] != head {
		t.Error("unexpected log id")
	}

	if len(reply.Commits) != 1 {
		t.Fatal("unexpected commits")
	}
}