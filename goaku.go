package goaku

/*
#cgo LDFLAGS: -lakumuli -L/usr/lib/x86_64-linux-gnu/ -lapr-1 -laprutil-1 -lboost_coroutine -lboost_context -lboost_system
#cgo CFLAGS: -I/usr/include/apr-1.0/
#include <akumuli.h>
#include <akumuli_def.h>
#include <akumuli_config.h>
#include <akumuli_version.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func Initialize() {
	C.aku_initialize(nil)
}

func akuError(status C.aku_Status) error {
	if status == C.AKU_SUCCESS {
		return nil
	}
	msg := C.GoString(C.aku_error_message(C.int(status)))
	return errors.New(msg)
}

func aprError(status C.apr_status_t) error {
	if status != 0 {
		buf := make([]byte, 1024)
		msg := C.apr_strerror(status, (*C.char)(unsafe.Pointer(&buf[0])),
			C.apr_size_t(len(buf)))
		errmsg := C.GoString(msg)
		return errors.New(errmsg)
	}
	return nil
}

func CreateDatabase(file_name, metadata_path, volumes_path string, num_volumes int32,
	compression_threshold *uint32, window_size *uint64, max_cache *uint32) error {
	var ct C.uint32_t
	var ws C.uint64_t
	var mc C.uint32_t
	var pct *C.uint32_t = nil
	var pws *C.uint64_t = nil
	var pmc *C.uint32_t = nil
	if compression_threshold != nil {
		ct = C.uint32_t(*compression_threshold)
		pct = &ct
	}
	if window_size != nil {
		ws = C.uint64_t(*window_size)
		pws = &ws
	}
	if max_cache != nil {
		mc = C.uint32_t(*max_cache)
		pmc = &mc
	}
	status := C.aku_create_database(C.CString(file_name), C.CString(metadata_path),
		C.CString(volumes_path), C.int32_t(num_volumes),
		pct, pws, pmc, nil)
	// TODO: pass callbacks
	return aprError(status)
}

func RemoveDatabase(path string) error {
	status := C.aku_remove_database(C.CString(path), nil)
	return aprError(status)
}

type Database struct {
	impl *C.aku_Database
}

func OpenDatabase(path string) (Database, error) {
	var params C.aku_FineTuneParams
	var result Database
	result.impl = C.aku_open_database(C.CString(path), params)
	if result.impl == nil {
		status := C.aku_open_status(result.impl)
		return result, akuError(status)
	}
	return result, nil
}

type ParamId uint32
type Timestamp uint64

func (db Database) Write(pid ParamId, ts Timestamp, buffer []byte) error {
    var memrange C.aku_MemRange
	memrange.address = unsafe.Pointer(&buffer[0])
    memrange.length = C.uint32_t(len(buffer))
	status := C.aku_write(db.impl, C.aku_ParamId(pid), C.aku_TimeStamp(ts), memrange)
	return akuError(status)
}

func (db Database) Close() {
    C.aku_close_database(db.impl)
}
