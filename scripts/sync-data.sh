#!/usr/bin/env bash
set -euo pipefail
SRC="${1:-../changelogs-info/src/content}"
mkdir -p data/releases data/models
cp -f "$SRC"/releases/*.json data/releases/ 2>/dev/null || true
cp -f "$SRC"/models/*.json data/models/ 2>/dev/null || true
cp -f "$SRC"/compatibility.json data/ 2>/dev/null || true
cp -f "$SRC"/deprecations.json data/ 2>/dev/null || true
cp -f "$SRC"/recommendations.json data/ 2>/dev/null || true
echo "Synced data from $SRC"
