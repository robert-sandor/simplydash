#!/bin/bash

if (($# != 1)); then
  echo "Usage: ./release.sh <version>"
  exit 1
fi

git-changelog -o CHANGELOG.md .
git add CHANGELOG.md
git commit -m "release: $1"
git tag "$1"
git push --follow-tags origin
