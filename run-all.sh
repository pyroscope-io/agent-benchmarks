#!/bin/bash
pushd runner
go run . ../fibonacci-go-cpu-push/ ../fibonacci-go-mem-push/
popd
