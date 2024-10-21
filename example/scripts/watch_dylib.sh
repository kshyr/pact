#!/bin/bash

ROOT_DIR="~/projects/opengl_example"
SRC_DIR="$ROOT_DIR/lib"
BUILD_DIR="$ROOT_DIR/build/lib"

compile() {
	clang++ -shared -o $BUILD_DIR/libopengl_example.so \
	$SRC_DIR/*.cpp -Ilib -Ivendor/glfw/include \
	-Ivendor/glad/include -lglfw -fPIC
}

compile

while inotifywait -e close_write --recursive "$SRC_DIR"; do
	compile
done
