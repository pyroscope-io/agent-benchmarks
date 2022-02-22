#!/bin/bash
pushd runner > /dev/null
go run . ../fibonacci-go-cpu-push/ ../fibonacci-go-mem-push/
popd > /dev/null
