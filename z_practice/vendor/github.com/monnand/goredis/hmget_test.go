package goredis

import (
	"reflect"
	"testing"
)

func TestContainerToString(t *testing.T) {

	vals := make(map[string]string, 2)
	vals["f1"] = "hello"
	vals["f2"] = "world"

	args := make([]string, 0, 5)
	args = append(args, "myhash")
	args, _ = containerToString(reflect.ValueOf(vals), args)

	correct_args := make([]string, 5)
	correct_args[0] = "myhash"
	correct_args[1] = "f1"
	correct_args[2] = "hello"
	correct_args[3] = "f2"
	correct_args[4] = "world"

	for i, v := range correct_args {
		if args[i] != v {
			t.Fatalf("%dth argument should be %s, but it is %s", i, v, args[i])
		}
	}

}

func TestHmset(t *testing.T) {
	var client Client
	client.Addr = "127.0.0.1:6379"
	client.Db = 13

	vals := make(map[string]string, 2)
	vals["f1"] = "hello"
	vals["f2"] = "world"
	key := "myhash"

	client.Hmset(key, vals)

	for f, v := range vals {
		value, err := client.Hget(key, f)
		if err != nil {
			t.Fatalf("Database error: %v", err)
		}
		str := string(value)
		if str != v {
			t.Fatalf("field %s should be %s but it is %s", f, v, str)
		}
	}

}
