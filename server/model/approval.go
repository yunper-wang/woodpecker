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

// Approval records a single user's approval or decline for a blocked pipeline.
type Approval struct {
	ID         int64  `json:"id"          xorm:"pk autoincr 'id'"`
	PipelineID int64  `json:"pipeline_id" xorm:"INDEX 'pipeline_id'"`
	UserID     int64  `json:"user_id"     xorm:"'user_id'"`
	UserLogin  string `json:"user_login"  xorm:"'user_login'"`
	Action     string `json:"action"      xorm:"varchar(20) 'action'"`  // "approve" or "decline"
	CreatedAt  int64  `json:"created_at"  xorm:"'created_at'"`
} //	@name	Approval
