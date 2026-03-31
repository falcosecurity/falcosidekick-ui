// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import "testing"

func TestConvertToSeconds(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		// Valid units
		{name: "seconds", input: "30s", want: 30},
		{name: "minutes", input: "5m", want: 300},
		{name: "hours", input: "2h", want: 7200},
		{name: "days", input: "1d", want: 86400},
		{name: "weeks", input: "1w", want: 604800},
		{name: "months", input: "1M", want: 2592000},
		{name: "years", input: "1y", want: 31536000},
		{name: "long form minute", input: "5minute", want: 300},
		{name: "long form hours", input: "2hours", want: 7200},
		// Pure numeric (fallback)
		{name: "pure number", input: "60", want: 60},
		{name: "zero", input: "0s", want: 0},
		// Malformed inputs must return 0, not crash
		{name: "unrecognized suffix", input: "1x", want: 0},
		{name: "unrecognized suffix a", input: "1a", want: 0},
		{name: "symbol suffix", input: "1+", want: 0},
		{name: "multi char bad suffix", input: "5abc", want: 0},
		{name: "empty string", input: "", want: 0},
		{name: "only letters", input: "abc", want: 0},
		{name: "no number prefix", input: "s", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToSeconds(tt.input)
			if got != tt.want {
				t.Errorf("ConvertToSeconds(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
