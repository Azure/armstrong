package coverage

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

const azureRepoURL = "https://github.com/Azure/azure-rest-api-specs.git"
const azureRepoName = "azure-rest-api-specs"

type AzureRepo struct {
	SpecRootDir string
	Repo        *git.Repository
}

func NewAzureRepo(baseDir string) (*AzureRepo, error) {
	repoDir, err := filepath.Abs(filepath.Join(baseDir, azureRepoName))
	if err != nil {
		return nil, err
	}
	repoDir = filepath.Clean(repoDir)

	r, err := cloneAndOpenRepo(repoDir, azureRepoURL, nil)
	if err != nil {
		return nil, err
	}

	return &AzureRepo{
		SpecRootDir: filepath.Join(repoDir, "specification"),
		Repo:        r,
	}, nil
}

func (repo *AzureRepo) CheckoutRef(ref string) (changed bool, err error) {
	return checkoutRef(repo.Repo, ref, nil)
}

func cloneAndOpenRepo(repoDir string, repoUrl string, auth transport.AuthMethod) (*git.Repository, error) {
	var repoExist bool
	if _, err := os.Stat(repoDir); err == nil {
		r, err := git.PlainOpen(repoDir)
		if err != nil {
			return nil, err
		}
		remote, err := r.Remote("origin")
		if err != nil {
			return nil, err
		}
		// URLs must have at least one element, which will be used for fetching.
		if remote.Config().URLs[0] != repoUrl {
			os.RemoveAll(repoDir)
		} else {
			repoExist = true
		}
	}

	if !repoExist {
		if err := os.Mkdir(repoDir, 0755); err != nil {
			if !errors.Is(err, os.ErrExist) {
				return nil, fmt.Errorf("creating repo directory %q: %w", repoDir, err)
			}
		}

		if _, err := git.PlainClone(repoDir, false, &git.CloneOptions{
			URL:      repoUrl,
			Auth:     auth,
			Progress: os.Stderr,
		}); err != nil {
			os.RemoveAll(repoDir)
			return nil, fmt.Errorf("cloning repo %q: %w", repoUrl, err)

		}
	}

	r, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("open repo %q: %w", repoDir, err)
	}

	return r, nil
}

func checkoutRef(repo *git.Repository, ref string, auth transport.AuthMethod) (changed bool, err error) {
	wt, err := repo.Worktree()
	if err != nil {
		return false, err
	}
	ohead, err := repo.Head()
	if err != nil {
		return false, err
	}

	if err := repo.Fetch(&git.FetchOptions{
		Auth:     auth,
		Progress: os.Stderr,
		Force:    true,
	}); err != nil {
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return false, err
		}
	}

	opt := &git.CheckoutOptions{Force: true}
	if _, err := hex.DecodeString(ref); err == nil {
		opt.Hash = plumbing.NewHash(ref)
		if err := wt.Checkout(opt); err != nil {
			return false, err
		}
	} else {
		opt.Branch = plumbing.NewTagReferenceName(ref)
		if err := wt.Checkout(opt); err != nil {
			if !errors.Is(err, plumbing.ErrReferenceNotFound) {
				return false, err
			}

			opt.Branch = plumbing.NewRemoteReferenceName("origin", ref)
			if err := wt.Checkout(opt); err != nil {
				return false, err
			}
		}
	}
	head, _ := repo.Head()
	return head.Hash() == ohead.Hash(), nil
}
