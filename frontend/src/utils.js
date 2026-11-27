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

const PALETTE = [
  '#4C8DF6', '#2BB3A3', '#F2A33A', '#9B7BEA', '#E8618C', '#5BC0EB',
  '#8BC34A', '#FF8A65', '#7986CB', '#4DB6AC', '#BA68C8', '#90A4AE',
];

const PRIORITY_COLORS = {
  Emergency: '#B4232A',
  Alert: '#D6333B',
  Critical: '#EF4444',
  Error: '#F97316',
  Warning: '#F5B53F',
  Notice: '#4C8DF6',
  Informational: '#38BDF8',
  Debug: '#94A3B8',
};

export const PRIORITY_ORDER = Object.keys(PRIORITY_COLORS);

function hash(str) {
  let h = 0;
  for (let i = 0; i < str.length; i += 1) {
    h = (((h << 5) - h) + str.charCodeAt(i)) | 0; // eslint-disable-line no-bitwise
  }
  return Math.abs(h);
}

export const utils = {
  priorityToColor(priority) {
    return PRIORITY_COLORS[priority] || '#7A8699';
  },
  stringToColor(str) {
    return PALETTE[hash(String(str)) % PALETTE.length];
  },
  sortPriorities(list) {
    return [...list].sort((a, b) => {
      const ia = PRIORITY_ORDER.indexOf(a);
      const ib = PRIORITY_ORDER.indexOf(b);
      return (ia === -1 ? 99 : ia) - (ib === -1 ? 99 : ib);
    });
  },
  cssVar(name) {
    return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
  },
};

export default {
};
