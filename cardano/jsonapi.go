// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"encoding/json"
	"errors"

	// R U MAD?
	jscardanowasmbin "github.com/jollheef/go-js-cardano-wasm-bin"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

const MAX_OUTPUT_SIZE = 4096

func newString(instance wasm.Instance, s string) (p wasm.Value, err error) {
	alloc, ok := instance.Exports["alloc"]
	if !ok {
		err = errors.New("function alloc not found")
		return
	}

	p, err = alloc(len(s))
	if err != nil {
		return
	}

	memory := instance.Memory.Data()[p.ToI32():]

	for nth := 0; nth < len(s); nth++ {
		memory[nth] = s[nth]
	}

	memory[len(s)] = 0
	return
}

func newArray(instance wasm.Instance, array []byte) (p wasm.Value, err error) {
	alloc, ok := instance.Exports["alloc"]
	if !ok {
		err = errors.New("function alloc not found")
		return
	}

	p, err = alloc(len(array))
	if err != nil {
		return
	}

	memory := instance.Memory.Data()[p.ToI32():]

	for nth := 0; nth < len(array); nth++ {
		memory[nth] = array[nth]
	}
	return
}

func newArray0(instance wasm.Instance, size int) (p wasm.Value, err error) {
	alloc, ok := instance.Exports["alloc"]
	if !ok {
		err = errors.New("function alloc not found")
		return
	}

	p, err = alloc(size)
	if err != nil {
		return
	}

	memory := instance.Memory.Data()[p.ToI32():]

	for i := 0; i <= size; i++ {
		memory[i] = 0
	}
	return
}

func copyArray(instance wasm.Instance, array, size wasm.Value) (b []byte) {
	memory := instance.Memory.Data()[array.ToI32():]

	for i := int32(0); i < size.ToI32(); i++ {
		b = append(b, memory[i])
	}
	return
}

func jsonapi(function, input string) (output string, err error) {
	instance, err := wasm.NewInstance(jscardanowasmbin.Module)
	if err != nil {
		return
	}
	defer instance.Close()

	p, err := newString(instance, input)
	if err != nil {
		return
	}

	out, err := newArray0(instance, MAX_OUTPUT_SIZE)
	if err != nil {
		return
	}

	f, ok := instance.Exports[function]
	if !ok {
		err = errors.New("function " + function + " not found")
		return
	}

	ret, err := f(p.ToI32(), len(input), out.ToI32())
	if err != nil {
		return
	}

	raw := copyArray(instance, out, ret)

	var result struct {
		Failed bool
		Result json.RawMessage
		Loc    string
		Msg    string
	}
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return
	}
	if result.Failed {
		err = errors.New(result.Msg + " [" + result.Loc + "]")
		return
	}
	output = string(result.Result)
	return
}
