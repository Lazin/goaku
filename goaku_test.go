package goaku

import "testing"

func Test_initialization(t *testing.T) {
	Initialize()
}

func Test_create_database_error_handling(t *testing.T) {
    file_name := "/foo"
    err := CreateDatabase(file_name, file_name, file_name, 22, nil, nil, nil)
    if err == nil {
        t.Error("no error found")
    }
}
