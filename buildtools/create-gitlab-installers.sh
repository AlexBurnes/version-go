#!/usr/bin/env bash
set -euo pipefail

# Creates installer scripts pointing to GitLab release assets.
# Usage: ./create-gitlab-installers.sh <version> [target_dir]

VERSION="${1:-}"
TARGET_DIR="${2:-installers-gitlab}"

if [[ -z "$VERSION" ]]; then
    echo "Usage: $0 <version> [target_dir]"
    echo "Example: $0 v0.33.1"
    exit 1
fi

CLEAN_VERSION=$(echo "$VERSION" | sed 's/^v//')

mkdir -p "$TARGET_DIR"

PLATFORMS=("linux" "darwin")
ARCHS=("amd64" "arm64")

echo "[INFO] Creating GitLab installers for version $CLEAN_VERSION"

for platform in "${PLATFORMS[@]}"; do
    for arch in "${ARCHS[@]}"; do
        INSTALLER_NAME="version-${platform}-${arch}-install.sh"
        INSTALLER_PATH="${TARGET_DIR}/${INSTALLER_NAME}"

        sed -e "s/{{VERSION}}/$CLEAN_VERSION/g" \
            -e "s/{{PLATFORM}}/$platform/g" \
            -e "s/{{ARCH}}/$arch/g" \
            -e "s/{{BINARY_NAME}}/version/g" \
            "packaging/linux/installer-template-gitlab.sh" > "$INSTALLER_PATH"

        chmod +x "$INSTALLER_PATH"
        echo "[SUCCESS] Created: $INSTALLER_PATH"
    done
done

echo "[SUCCESS] GitLab installers created in $TARGET_DIR/"
