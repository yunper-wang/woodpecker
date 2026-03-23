// Copyright 2022 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package matrix

import (
	"fmt"
	"strings"

	"codeberg.org/6543/xyaml"

	pipeline_errors "go.woodpecker-ci.org/woodpecker/v3/pipeline/errors"
)

const (
	defaultLimitTags = 10
	defaultLimitAxis = 25
)

// Matrix represents the pipeline matrix.
type Matrix map[string][]string

// Axis represents a single permutation of entries from the pipeline matrix.
type Axis map[string]string

// String returns a string representation of an Axis as a comma-separated list
// of environment variables.
func (a Axis) String() string {
	var envs []string
	for k, v := range a {
		envs = append(envs, k+"="+v)
	}
	return strings.Join(envs, " ")
}

// Parse parses the Yaml matrix definition.
// limitAxis and limitTags control the maximum permutations and variable dimensions;
// pass 0 to use the package defaults.
func Parse(data []byte, limitAxis, limitTags int) ([]Axis, error) {
	if limitAxis <= 0 {
		limitAxis = defaultLimitAxis
	}
	if limitTags <= 0 {
		limitTags = defaultLimitTags
	}

	include, excludeFromList, err := parseList(data)
	if err == nil && len(include) != 0 {
		return applyExclude(include, excludeFromList), nil
	}

	matrix, excludeFromMatrix, err := parse(data)
	if err != nil {
		return nil, err
	}

	if len(matrix) == 0 {
		return []Axis{}, nil
	}

	axisList, err := calc(matrix, limitAxis, limitTags)
	if err != nil {
		return axisList, err
	}
	return applyExclude(axisList, excludeFromMatrix), nil
}

// ParseString parses the Yaml string matrix definition.
func ParseString(data string) ([]Axis, error) {
	return Parse([]byte(data), 0, 0)
}

func calc(matrix Matrix, limitAxis, limitTags int) ([]Axis, error) {
	// calculate number of permutations and extract the list of tags
	// (ie go_version, redis_version, etc)
	var perm int
	var tags []string
	for k, v := range matrix {
		perm *= len(v)
		if perm == 0 {
			perm = len(v)
		}
		tags = append(tags, k)
	}

	// structure to hold the transformed result set
	var axisList []Axis

	// for each axis calculate the unique set of values that should be used.
	for p := 0; p < perm; p++ {
		axis := map[string]string{}
		decrease := perm
		for i, tag := range tags {
			elems := matrix[tag]
			decrease /= len(elems)
			elem := p / decrease % len(elems)
			axis[tag] = elems[elem]

			// enforce a maximum number of tags in the pipeline matrix.
			if i > limitTags {
				break
			}
		}

		// append to the list of axis.
		axisList = append(axisList, axis)

		// enforce a maximum number of axis that should be calculated.
		if p > limitAxis {
			return axisList, &pipeline_errors.PipelineError{
				Message: fmt.Sprintf("matrix exceeds maximum of %d axis", limitAxis),
				Type:    pipeline_errors.PipelineErrorTypeCompiler,
			}
		}
	}

	return axisList, nil
}

// applyExclude filters out any axis entries that match any exclude pattern.
func applyExclude(axisList []Axis, excludes []Axis) []Axis {
	if len(excludes) == 0 {
		return axisList
	}
	var result []Axis
	for _, axis := range axisList {
		excluded := false
		for _, ex := range excludes {
			if axisMatchesExclude(axis, ex) {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, axis)
		}
	}
	return result
}

// axisMatchesExclude returns true if all key-value pairs in exclude exist in axis.
func axisMatchesExclude(axis, exclude Axis) bool {
	for k, v := range exclude {
		if axis[k] != v {
			return false
		}
	}
	return true
}

func parse(raw []byte) (Matrix, []Axis, error) {
	// First extract exclude list
	excludeData := struct {
		Matrix struct {
			Exclude []Axis `yaml:"exclude"`
		}
	}{}
	_ = xyaml.Unmarshal(raw, &excludeData)

	// Then parse the full map and remove reserved keys
	fullData := struct {
		Matrix map[string]interface{}
	}{}
	if err := xyaml.Unmarshal(raw, &fullData); err != nil {
		return nil, nil, &pipeline_errors.PipelineError{Message: err.Error(), Type: pipeline_errors.PipelineErrorTypeCompiler}
	}

	matrix := Matrix{}
	for k, v := range fullData.Matrix {
		if k == "include" || k == "exclude" {
			continue
		}
		slice, ok := v.([]interface{})
		if !ok {
			continue
		}
		var vals []string
		for _, item := range slice {
			vals = append(vals, fmt.Sprintf("%v", item))
		}
		matrix[k] = vals
	}
	return matrix, excludeData.Matrix.Exclude, nil
}

func parseList(raw []byte) ([]Axis, []Axis, error) {
	data := struct {
		Matrix struct {
			Include []Axis
			Exclude []Axis
		}
	}{}

	if err := xyaml.Unmarshal(raw, &data); err != nil {
		return nil, nil, &pipeline_errors.PipelineError{Message: err.Error(), Type: pipeline_errors.PipelineErrorTypeCompiler}
	}
	return data.Matrix.Include, data.Matrix.Exclude, nil
}
