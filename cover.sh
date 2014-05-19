#!/bin/bash

# Run test coverage on each subdirectories and merge the coverage profile.

echo "mode: count" > profile.cov
fail=0

# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d);
do
if ls $dir/*.go &> /dev/null; then
    go test -covermode=count -coverprofile=$dir/profile.tmp $dir || fail=1
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> profile.cov
        rm $dir/profile.tmp
    fi
fi
done
