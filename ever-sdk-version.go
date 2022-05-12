package main
/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lton_client
#include <stdbool.h>
#include <tonclient.h>
*/
import "C"
import (
    "fmt"
    "encoding/json"
    "unsafe"
)

func tc_string(str string) C.tc_string_data_t {
    return C.tc_string_data_t {
        content: C.CString(str),
        len: C.uint(len(str)),
    }
}

func tc_request_sync(context uint, method string, optionsJsonBlob string) *C.tc_string_handle_t {
    return C.tc_request_sync(C.uint(context), tc_string(method), tc_string(optionsJsonBlob));
}

func tc_read_string(context_handler *C.tc_string_handle_t) string {
    handler := C.tc_read_string(context_handler)
    return C.GoStringN((*C.char)(unsafe.Pointer(handler.content)), C.int(handler.len))
}

func tc_result(jsonBlob string) uint {
    type tc_string struct {
        Result uint
    }

    var res tc_string;
    err := json.Unmarshal([]byte(jsonBlob), &res)
    if err != nil {
        fmt.Println("json.Unmarshal[result]error:", err)
    }
    return res.Result
}

func tc_response(jsonBlob string) string {
    type tc_version struct {
        Version string
    }
    type tc_result struct {
        Result tc_version
    }

    var res tc_result;
    err := json.Unmarshal([]byte(jsonBlob), &res)
    if err != nil {
        fmt.Println("json.Unmarshal[response]error:", err)
    }
    return res.Result.Version
}

func main() {
    context_handle := C.tc_create_context(tc_string(""))
    var result string = tc_read_string(context_handle)
    var context uint = tc_result(result)
    fmt.Println(tc_response(tc_read_string(tc_request_sync(context, "client.version", ""))))
}
