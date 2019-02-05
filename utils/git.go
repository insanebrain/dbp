package utils

import (
    "errors"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
    "io"
)

func GetCommitFiles(repositoryPath string, hash string) ([]string, error) {

    var files []string = nil
    var co *object.Commit

    r, _ := git.PlainOpen(repositoryPath)

    if len(hash) < 40 {
        cl, _ := r.CommitObjects()
        _ = cl.ForEach(func(c *object.Commit) error {
            if c.Hash.String()[:len(hash)] == hash {
                co = c
                cl.Close()
            }
            return nil
        })
    } else {
        co, _ = r.CommitObject(plumbing.NewHash(hash))
    }

    if co == nil {
        return nil, errors.New("commit hash not found")
    }

    parentCommit, err := co.Parents().Next()
    if err != nil {
        if err == io.EOF {
            fl, _ := co.Files()
            _ = fl.ForEach(func(f *object.File) error {
                files = append(files, f.Name)
                return nil
            })
        } else {
            return nil, err
        }
    } else {
        patch, _ := parentCommit.Patch(co)
        for _, fp := range patch.FilePatches() {
            from, to := fp.Files()
            if from == nil {
                files = append(files, to.Path())
            } else if to == nil {
                files = append(files, from.Path())
            } else if from.Path() != to.Path() {
                files = append(files, to.Path())
            } else {
                files = append(files, from.Path())
            }
        }
    }

    return files, nil
}

func GetStatus(repositoryPath string) ([]string, error) {

    var files []string = nil

    r, _ := git.PlainOpen(repositoryPath)

    workTree, err := r.Worktree()
    if err != nil {
        return nil, err
    }
    status, _ := workTree.Status()

    for path, status := range status {
        if status.Staging == git.Unmodified && status.Worktree == git.Unmodified {
            continue
        }

        if status.Staging == git.Renamed {
            path = status.Extra
        }
        files = append(files, path)
    }

    return files, nil

}