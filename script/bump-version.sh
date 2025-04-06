#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# =============================================
# Helper Functions
# =============================================

# Get the latest version tag from git
# Returns v0.0.0 if no tags exist
get_latest_version() {
    echo "🔍 Checking latest version tag..."
    
    local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null)
    if [ -z "$latest_tag" ]; then
        echo "⚠️  No version tags found. Starting with v0.0.0"
        echo "v0.0.0"
    else
        echo "✅ Found latest version: $latest_tag"
        echo "$latest_tag"
    fi
}

# Bump version according to semver rules
# Parameters:
#   $1 - Current version (e.g., v1.2.3)
#   $2 - Bump type (major, minor, or patch)
bump_version() {
    local current_version=$1
    local bump_type=$2

    # Split version string into major, minor, patch
    IFS='.' read -r -a version_parts <<< "${current_version#v}"
    local major=${version_parts[0]}
    local minor=${version_parts[1]}
    local patch=${version_parts[2]}

    # Increment version based on bump type
    case $bump_type in
        "major")
            echo "📈 Bumping major version..."
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            echo "📈 Bumping minor version..."
            minor=$((minor + 1))
            patch=0
            ;;
        "patch")
            echo "📈 Bumping patch version..."
            patch=$((patch + 1))
            ;;
        *)
            echo "❌ Error: Invalid bump type '$bump_type'. Must be major, minor, or patch"
            exit 1
            ;;
    esac

    echo "v${major}.${minor}.${patch}"
}

# =============================================
# Main Script
# =============================================

echo "🚀 Starting version bump process..."
echo "----------------------------------------"

# Get current version
latest_version=$(get_latest_version)
echo "----------------------------------------"

# Analyze last commit message to determine bump type
echo "📝 Analyzing last commit message..."
last_commit_message=$(git log -1 --pretty=%B)
echo "Last commit: $last_commit_message"

# Determine version bump type based on conventional commits
# feat!: or fix!: = major bump
# feat: = minor bump
# everything else = patch bump
if [[ "$last_commit_message" =~ ^(feat|fix|refactor|perf|revert)!: ]]; then
    bump_type="major"
    echo "🔥 Breaking change detected! Bumping major version"
elif [[ "$last_commit_message" =~ ^feat: ]]; then
    bump_type="minor"
    echo "✨ New feature detected! Bumping minor version"
else
    bump_type="patch"
    echo "🔧 Patch update! Bumping patch version"
fi

echo "----------------------------------------"

# Calculate new version
new_version=$(bump_version "$latest_version" "$bump_type")
echo "🏷  Version bump summary:"
echo "   From: $latest_version"
echo "   To:   $new_version"
echo "----------------------------------------"

# Uncomment the following line when ready to create the tag
# git tag "$new_version"
echo "✅ Version bump completed!"
echo "💡 To create the tag, uncomment the git tag line in the script"
