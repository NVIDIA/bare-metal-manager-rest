/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package uniquefifo

type Item[K comparable] interface {
	Key() K
}

type UniqueFIFO[K comparable, V Item[K]] struct {
	items    []V
	contains map[K]struct{}
}

func New[K comparable, V Item[K]]() *UniqueFIFO[K, V] {
	return &UniqueFIFO[K, V]{
		items:    make([]V, 0),
		contains: make(map[K]struct{}),
	}
}

func (q *UniqueFIFO[K, V]) Push(item V) bool {
	key := item.Key()
	if _, ok := q.contains[key]; ok {
		return false
	}

	q.items = append(q.items, item)
	q.contains[key] = struct{}{}
	return true
}

func (q *UniqueFIFO[K, V]) Pop() (V, bool) {
	var zero V
	if len(q.items) == 0 {
		return zero, false
	}

	item := q.items[0]
	delete(q.contains, item.Key())
	q.items = q.items[1:]
	return item, true
}

func (q *UniqueFIFO[K, V]) Contains(key K) bool {
	_, ok := q.contains[key]
	return ok
}
