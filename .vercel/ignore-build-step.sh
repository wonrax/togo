#!/bin/bash

echo "VERCEL_GIT_COMMIT_REF: $VERCEL_GIT_COMMIT_REF"

if [[ "$VERCEL_GIT_COMMIT_REF" == "master" ]] ; then
    # Depends on whether the files changed or not
    git diff HEAD^ HEAD --quiet .

else
    # Don't build
    echo "🛑 - Build cancelled"
    exit 0;
fi
