#!/bin/sh

# usage: ./install-toolchain.sh $PREFIX # must be absolute path

set -ex

git submodule update --init

PREFIX=$1

mkdir -p $PREFIX
mkdir -p $PREFIX/bin
mkdir -p $PREFIX/clang
mkdir -p $PREFIX/gcc

go build -o $PREFIX/bin ./lfi-as
go build -o $PREFIX/bin ./lfi-compile
go build -o $PREFIX/bin ./specinvoke

cp LICENSE $PREFIX
cp README.md $PREFIX

(cd lfi-leg && make lfi-leg-arm64 && mv lfi-leg-arm64 $PREFIX/bin)

export PATH=$PREFIX/bin:$PATH

(cd lfi-verify && cargo build --release && mv ./target/release/lfi-verify $PREFIX/bin/lfi-verify-verbose)

(cd lfi-veribdd && knit && mv lfi-verify $PREFIX/bin)

(cd lfi-run && knit && mv lfi-run $PREFIX/bin)

# install clang
(cd toolchain/lfi-clang && ./install-toolchain.sh $PREFIX/clang)

# install gcc
(cd toolchain/lfi-gcc && ./install-toolchain.sh $PREFIX/gcc)
