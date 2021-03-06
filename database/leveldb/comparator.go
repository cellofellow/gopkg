// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leveldb

// #include <leveldb/c.h>
import "C"

// DestroyComparator deallocates a *C.leveldb_comparator_t.
//
// This is provided as a convienience to advanced users that have implemented
// their own comparators in C in their own code.
func DestroyComparator(cmp *C.leveldb_comparator_t) {
	C.leveldb_comparator_destroy(cmp)
}
