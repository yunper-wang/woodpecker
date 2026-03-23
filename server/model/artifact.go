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

package model

// Artifact represents a file artifact produced by a pipeline workflow
// that can be downloaded by users or consumed by downstream workflows.
type Artifact struct {
	ID         int64  `json:"id"          xorm:"pk autoincr 'id'"`
	RepoID     int64  `json:"repo_id"     xorm:"INDEX 'repo_id'"`
	PipelineID int64  `json:"pipeline_id" xorm:"INDEX 'pipeline_id'"`
	WorkflowID int64  `json:"workflow_id" xorm:"'workflow_id'"`
	Name       string `json:"name"        xorm:"VARCHAR(250) 'name'"`
	FileSize   int64  `json:"file_size"   xorm:"'file_size'"`
	Data       []byte `json:"-"           xorm:"LONGBLOB 'data'"`
	CreatedAt  int64  `json:"created_at"  xorm:"'created_at'"`
} //	@name	Artifact
