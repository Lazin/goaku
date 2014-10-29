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

func initialize() {
    C.aku_initialize(nil);
}
