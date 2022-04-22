package io

import (
    "github.com/TylerBrock/colorjson"
)


func PrettyJson(obj interface{}) string {
    str, _ := colorjson.Marshal(obj)
    return string(str)
}
