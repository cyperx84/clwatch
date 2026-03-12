#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:?Usage: build-release.sh <version>}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST="${REPO_ROOT}/dist"
MODULE="./cmd/clwatch"
LDFLAGS="-s -w -X main.version=${VERSION}"

PLATFORMS=(
  "darwin/arm64"
  "darwin/amd64"
  "linux/arm64"
  "linux/amd64"
  "windows/amd64"
)

rm -rf "${DIST}"
mkdir -p "${DIST}"

cd "${REPO_ROOT}"

for platform in "${PLATFORMS[@]}"; do
  GOOS="${platform%/*}"
  GOARCH="${platform#*/}"
  BINARY="clwatch"
  ARCHIVE_NAME="clwatch-${VERSION}-${GOOS}-${GOARCH}"

  if [ "${GOOS}" = "windows" ]; then
    BINARY="clwatch.exe"
  fi

  echo "Building ${GOOS}/${GOARCH}..."
  GOOS="${GOOS}" GOARCH="${GOARCH}" CGO_ENABLED=0 \
    go build -ldflags "${LDFLAGS}" -o "${DIST}/${ARCHIVE_NAME}/${BINARY}" "${MODULE}"

  # Create tar.gz
  tar -czf "${DIST}/${ARCHIVE_NAME}.tar.gz" -C "${DIST}" "${ARCHIVE_NAME}"
  rm -rf "${DIST:?}/${ARCHIVE_NAME}"
done

echo ""
echo "=== SHA256 Checksums ==="
echo ""
printf "%-50s %s\n" "Archive" "SHA256"
printf "%-50s %s\n" "-------" "------"
for archive in "${DIST}"/*.tar.gz; do
  name="$(basename "${archive}")"
  hash="$(shasum -a 256 "${archive}" | awk '{print $1}')"
  printf "%-50s %s\n" "${name}" "${hash}"
done
