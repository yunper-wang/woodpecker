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

func (s storage) ApprovalCreate(approval *model.Approval) error {
	approval.CreatedAt = time.Now().Unix()
	_, err := s.engine.Insert(approval)
	return err
}

func (s storage) ApprovalListForPipeline(pipelineID int64) ([]*model.Approval, error) {
	var approvals []*model.Approval
	return approvals, s.engine.Where("pipeline_id = ?", pipelineID).Asc("id").Find(&approvals)
}
