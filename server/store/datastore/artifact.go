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

package datastore

import (
	"time"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

func (s storage) ArtifactCreate(artifact *model.Artifact) error {
	artifact.CreatedAt = time.Now().Unix()
	_, err := s.engine.Insert(artifact)
	return err
}

func (s storage) ArtifactFind(id int64) (*model.Artifact, error) {
	artifact := &model.Artifact{}
	return artifact, wrapGet(s.engine.ID(id).Get(artifact))
}

func (s storage) ArtifactListForPipeline(pipelineID int64) ([]*model.Artifact, error) {
	var artifacts []*model.Artifact
	return artifacts, s.engine.Where("pipeline_id = ?", pipelineID).Omit("data").Find(&artifacts)
}

func (s storage) ArtifactDelete(id int64) error {
	_, err := s.engine.ID(id).Delete(new(model.Artifact))
	return err
}
