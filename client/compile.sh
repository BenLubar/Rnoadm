#!/bin/sh -e

if cd closure-library 2>/dev/null
then
	git pull
	cd ..
else
	git clone http://code.google.com/p/closure-library/
fi

if test ! -f compiler.jar
then
	wget http://dl.google.com/closure-compiler/compiler-latest.zip
	unzip compiler-latest.zip compiler.jar
fi

python closure-library/closure/bin/build/closurebuilder.py \
	--root=closure-library/ \
	--root=rnoadm/ \
	--namespace="rnoadm.main" \
	--output_mode=compiled \
	--compiler_jar=compiler.jar \
	--compiler_flags="--define=goog.DEBUG" \
	--compiler_flags="--define=goog.asserts.ENABLE_ASSERTS" \
	--compiler_flags="--accept_const_keyword" \
	--compiler_flags="--compilation_level=ADVANCED_OPTIMIZATIONS" \
	--compiler_flags="--language_in=ECMASCRIPT5_STRICT" \
	--compiler_flags="--warning_level=VERBOSE" \
	--compiler_flags="--use_types_for_optimization" \
	> ../resource/client.js

cd ../resource
go run update.go client.js
cd ..
go install
