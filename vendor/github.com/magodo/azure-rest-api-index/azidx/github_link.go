package azidx

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/scanner"

	"github.com/go-git/go-git/v5"
	"github.com/go-openapi/jsonreference"
)

func BuildGithubLink(ref jsonreference.Ref, commit, specdir string) (string, error) {
	repo, err := git.PlainOpen(filepath.Dir(specdir))
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return "", err
		}
	} else {
		head, err := repo.Head()
		if err != nil {
			return "", err
		}
		if repoCommit := head.Hash().String(); repoCommit != commit {
			return "", fmt.Errorf("repository commit %q not equals to the commit the index is built %q", repoCommit, commit)
		}
	}

	fpath, err := filepath.Abs(filepath.Join(specdir, ref.GetURL().Path))
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(fpath)
	if err != nil {
		return "", err
	}
	offset, err := JSONPointerOffset(*ref.GetPointer(), string(b))
	if err != nil {
		return "", err
	}
	var sc scanner.Scanner
	sc.Init(bytes.NewBuffer(b))
	fmt.Println(offset)
	for i := 0; i < int(offset); i++ {
		sc.Next()
	}
	pos := sc.Pos()

	specdir, err = filepath.Abs(specdir)
	if err != nil {
		return "", err
	}

	relFile, err := filepath.Rel(specdir, fpath)
	if err != nil {
		return "", err
	}

	return "https://github.com/Azure/azure-rest-api-specs/blob/" + commit + "/specification/" + relFile + "#L" + strconv.Itoa(pos.Line), nil
}
