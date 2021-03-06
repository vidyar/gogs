// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/go-martini/martini"

	"github.com/gogits/git"

	"github.com/gogits/gogs/modules/base"
	"github.com/gogits/gogs/modules/middleware"
)

func SingleDownload(ctx *middleware.Context, params martini.Params) {
	// Get tree path
	treename := params["_1"]

	blob, err := ctx.Repo.Commit.GetBlobByPath(treename)
	if err != nil {
		ctx.Handle(404, "repo.SingleDownload(GetBlobByPath)", err)
		return
	}

	data, err := blob.Data()
	if err != nil {
		ctx.Handle(404, "repo.SingleDownload(Data)", err)
		return
	}

	contentType, isTextFile := base.IsTextFile(data)
	_, isImageFile := base.IsImageFile(data)
	ctx.Res.Header().Set("Content-Type", contentType)
	if !isTextFile && !isImageFile {
		ctx.Res.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(treename))
		ctx.Res.Header().Set("Content-Transfer-Encoding", "binary")
	}
	ctx.Res.Write(data)
}

func ZipDownload(ctx *middleware.Context, params martini.Params) {
	commitId := ctx.Repo.CommitId
	archivesPath := filepath.Join(ctx.Repo.GitRepo.Path, "archives/zip")
	if !com.IsDir(archivesPath) {
		if err := os.MkdirAll(archivesPath, 0755); err != nil {
			ctx.Handle(404, "ZipDownload -> os.Mkdir(archivesPath)", err)
			return
		}
	}

	archivePath := filepath.Join(archivesPath, commitId+".zip")

	if com.IsFile(archivePath) {
		ctx.ServeFile(archivePath, ctx.Repo.Repository.Name+".zip")
		return
	}

	err := ctx.Repo.Commit.CreateArchive(archivePath, git.AT_ZIP)
	if err != nil {
		ctx.Handle(404, "ZipDownload -> CreateArchive "+archivePath, err)
		return
	}

	ctx.ServeFile(archivePath, ctx.Repo.Repository.Name+".zip")
}

func TarGzDownload(ctx *middleware.Context, params martini.Params) {
	commitId := ctx.Repo.CommitId
	archivesPath := filepath.Join(ctx.Repo.GitRepo.Path, "archives/targz")
	if !com.IsDir(archivesPath) {
		if err := os.MkdirAll(archivesPath, 0755); err != nil {
			ctx.Handle(404, "TarGzDownload -> os.Mkdir(archivesPath)", err)
			return
		}
	}

	archivePath := filepath.Join(archivesPath, commitId+".tar.gz")

	if com.IsFile(archivePath) {
		ctx.ServeFile(archivePath, ctx.Repo.Repository.Name+".tar.gz")
		return
	}

	err := ctx.Repo.Commit.CreateArchive(archivePath, git.AT_TARGZ)
	if err != nil {
		ctx.Handle(404, "TarGzDownload -> CreateArchive "+archivePath, err)
		return
	}

	ctx.ServeFile(archivePath, ctx.Repo.Repository.Name+".tar.gz")
}
