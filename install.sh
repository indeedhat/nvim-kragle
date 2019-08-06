#!/bin/bash

os=`uname`

if [ "Linux" == $os ]; then
    echo "linux is the default, skipping"
elif [ "Darwin" == $os ]; then
    echo "Swapping out for darwin build"
    mv ./kragle ./kragle.linux
    mv ./kragle.darwin kragle
else
    echo "OS not supported"
fi
