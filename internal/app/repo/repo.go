package repo

import "context"

type Repo interface {
	Tx(ctx context.Context, repo Repo) (Repo, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type RepoImpl struct {
}

func (rp *RepoImpl) Tx(ctx context.Context, repo Repo) (Repo, error) {
	return nil, nil
}

func (rp *RepoImpl) Commit(ctx context.Context) error {
	return nil
}

func (rp *RepoImpl) Rollback(ctx context.Context) error {
	return nil
}
