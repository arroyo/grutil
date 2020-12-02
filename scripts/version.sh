#!/bin/sh

echo "{"
echo '  "version": "'$(git tag -l | sort -V | tail -n 1)'"'
echo '  "date": "'$(git log -1 --format=%ai $(git tag -l | sort -V | tail -n 1))'"'
echo "}"
