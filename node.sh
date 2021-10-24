#!/bin/bash -x

[ ! -d web/lvnfe/node_modules ] && (cd web/lvnfe; npm i)
(cd web/lvnfe; npm start)
