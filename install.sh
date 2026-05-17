#!/bin/sh
# install.sh: download and install the `mm` binary from the latest mumei-md release.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/hir4ta/mumei-md/main/install.sh | sh
#
# Environment overrides:
#   VERSION       pin a release tag, e.g. v0.1.0 (default: latest)
#   INSTALL_DIR   target install directory (default: $HOME/.local/bin)

set -eu

REPO="hir4ta/mumei-md"
BIN="mm"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${VERSION:-}"

info() { printf 'install: %s\n' "$*"; }
err()  { printf 'install: %s\n' "$*" >&2; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || err "required command not found: $1"; }
need curl
need tar
need uname

case "$(uname -s)" in
  Linux)  OS=linux  ;;
  Darwin) OS=darwin ;;
  *) err "unsupported OS: $(uname -s)" ;;
esac

case "$(uname -m)" in
  x86_64|amd64)  ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) err "unsupported architecture: $(uname -m)" ;;
esac

if [ -z "$VERSION" ]; then
  redirect_url=$(curl -fsSI -o /dev/null -w '%{redirect_url}' \
    "https://github.com/${REPO}/releases/latest" || true)
  VERSION="${redirect_url##*/}"
  [ -n "$VERSION" ] || err "could not resolve latest release for ${REPO}"
fi

version_no_v="${VERSION#v}"
asset="mm_${version_no_v}_${OS}_${ARCH}.tar.gz"
url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT INT TERM

info "downloading ${asset}"
curl -fsSL "$url" -o "$tmp/$asset" || err "download failed: $url"

info "extracting"
tar -xzf "$tmp/$asset" -C "$tmp"
[ -f "$tmp/$BIN" ] || err "archive did not contain expected binary: $BIN"

mkdir -p "$INSTALL_DIR"
mv "$tmp/$BIN" "$INSTALL_DIR/$BIN"
chmod +x "$INSTALL_DIR/$BIN"

info "installed $BIN ${VERSION} to $INSTALL_DIR/$BIN"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) info "note: $INSTALL_DIR is not in your PATH; add it to use \`$BIN\` directly" ;;
esac
