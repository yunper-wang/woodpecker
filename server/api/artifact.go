// Copyright 2025 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"go.woodpecker-ci.org/woodpecker/v3/server/router/middleware/session"
	"go.woodpecker-ci.org/woodpecker/v3/server/store"
)

// GetArtifacts
//
//	@Summary	List artifacts for a pipeline
//	@Router		/repos/{repo_id}/pipelines/{number}/artifacts [get]
//	@Produce	json
//	@Success	200	{array}		Artifact
//	@Tags		Artifacts
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Param		repo_id			path	int		true	"the repository id"
//	@Param		number			path	int		true	"the number of the pipeline"
func GetArtifacts(c *gin.Context) {
	_store := store.FromContext(c)
	repo := session.Repo(c)

	num, err := strconv.ParseInt(c.Param("number"), 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pl, err := _store.GetPipelineNumber(repo, num)
	if err != nil {
		handleDBError(c, err)
		return
	}

	artifacts, err := _store.ArtifactListForPipeline(pl.ID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, artifacts)
}

// UploadArtifact
//
//	@Summary	Upload an artifact for a pipeline
//	@Router		/repos/{repo_id}/pipelines/{number}/artifacts [post]
//	@Produce	json
//	@Success	200	{object}	Artifact
//	@Tags		Artifacts
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Param		repo_id			path	int		true	"the repository id"
//	@Param		number			path	int		true	"the number of the pipeline"
//	@Param		name			query	string	true	"artifact name"
//	@Param		workflow_id	query	int		false	"workflow id"
func UploadArtifact(c *gin.Context) {
	_store := store.FromContext(c)
	repo := session.Repo(c)

	num, err := strconv.ParseInt(c.Param("number"), 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pl, err := _store.GetPipelineNumber(repo, num)
	if err != nil {
		handleDBError(c, err)
		return
	}

	name := c.Query("name")
	if name == "" {
		c.String(http.StatusBadRequest, "artifact name is required")
		return
	}

	workflowID, _ := strconv.ParseInt(c.Query("workflow_id"), 10, 64)

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	artifact := &model.Artifact{
		RepoID:     repo.ID,
		PipelineID: pl.ID,
		WorkflowID: workflowID,
		Name:       name,
		FileSize:   int64(len(data)),
		Data:       data,
	}

	if err := _store.ArtifactCreate(artifact); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	artifact.Data = nil // don't return raw data in response
	c.JSON(http.StatusOK, artifact)
}

// DownloadArtifact
//
//	@Summary	Download a specific artifact
//	@Router		/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id} [get]
//	@Produce	application/octet-stream
//	@Success	200
//	@Tags		Artifacts
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Param		repo_id			path	int		true	"the repository id"
//	@Param		number			path	int		true	"the number of the pipeline"
//	@Param		artifact_id		path	int		true	"the artifact id"
func DownloadArtifact(c *gin.Context) {
	_store := store.FromContext(c)

	artifactID, err := strconv.ParseInt(c.Param("artifact_id"), 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	artifact, err := _store.ArtifactFind(artifactID)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.DataFromReader(http.StatusOK, artifact.FileSize, "application/octet-stream",
		bytes.NewReader(artifact.Data), map[string]string{
			"Content-Disposition": `attachment; filename="` + artifact.Name + `"`,
		})
}

// DeleteArtifact
//
//	@Summary	Delete an artifact
//	@Router		/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id} [delete]
//	@Success	204
//	@Tags		Artifacts
//	@Param		Authorization	header	string	true	"Insert your personal access token"	default(Bearer <personal access token>)
//	@Param		repo_id			path	int		true	"the repository id"
//	@Param		number			path	int		true	"the number of the pipeline"
//	@Param		artifact_id		path	 int		true	"the artifact id"
func DeleteArtifact(c *gin.Context) {
	_store := store.FromContext(c)

	artifactID, err := strconv.ParseInt(c.Param("artifact_id"), 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := _store.ArtifactDelete(artifactID); err != nil {
		handleDBError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
