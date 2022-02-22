#!/bin/bash
pushd runner > /dev/null
go run . ../fibonacci-go-cpu-push/ ../fibonacci-go-mem-push/ ../fibonacci-go-all-push/
popd > /dev/null
