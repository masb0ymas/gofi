#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# =============================================
# Helper Functions
# =============================================

# Get the previous tag or create initial tag if none exists
get_version_tags() {
    echo "🔍 Checking version tags..."

    local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null)
    
    if [ -z "$latest_tag" ]; then
        PREV_TAG=""
        NEW_TAG="v0.0.1"
        echo "⚠️  No previous tags found. Starting with v0.0.1"
    else
        PREV_TAG=$(git describe --tags --abbrev=0)
        NEW_TAG=$(echo $PREV_TAG | awk -F. -v OFS=. '{++$NF} 1')
        echo "✅ Previous tag found: $PREV_TAG"
    fi
}

# Generate or update changelog content
generate_changelog() {
    echo "📝 Generating/updating changelog..."
    
    TEMP_CHANGELOG=$(mktemp)
    
    {
        if [ -f CHANGELOG.md ]; then
            sed '1,/^## /d' CHANGELOG.md > "$TEMP_CHANGELOG"
        fi
        
        echo "# Changelog"
        echo ""
        echo "## $NEW_TAG ($(date +%Y-%m-%d))"
        echo ""
        
        if [ -z "$PREV_TAG" ]; then
            echo "📦 Initial release"
            git log --pretty=format:"* %s" HEAD
        else
            echo "📦 Changes since $PREV_TAG:"
            git log --pretty=format:"* %s" $PREV_TAG..HEAD
        fi
        
        echo ""
        
        if [ -f "$TEMP_CHANGELOG" ]; then
            cat "$TEMP_CHANGELOG"
        fi
    } > CHANGELOG.md
    
    rm -f "$TEMP_CHANGELOG"
    
    echo "✅ Changelog updated successfully"
}

# Create and push new release
create_release() {
    echo "🚀 Creating new release..."
    
    # Stage and commit changelog
    git add .
    git commit -m "chore(release): bump to version $NEW_TAG"
    
    # Create and push tag
    git tag -a "$NEW_TAG" -m "Release version $NEW_TAG"
    git push origin "$NEW_TAG"
    
    # Push changes to main branch
    git push origin main
}

# =============================================
# Main Script
# =============================================

echo "🎉 Starting release process..."
echo "----------------------------------------"

# Step 1: Get version tags
get_version_tags

# Step 2: Display version information
echo -e "\n📋 Release Information"
echo "----------------------------------------"
echo "| Previous Release Tag  | ${PREV_TAG:-None}"
echo "| New Release Tag       | $NEW_TAG"
echo "----------------------------------------"

# Step 3: Generate changelog
echo -e "\n📄 Changelog Generation"
echo "----------------------------------------"
generate_changelog

# Step 4: Create and push release
echo -e "\n📦 Release Creation"
echo "----------------------------------------"
create_release

# Step 5: Completion
echo -e "\n✨ Release completed successfully!"
echo "🏷  New version: $NEW_TAG"
echo "📚 Changelog has been updated"
echo "🌟 All changes have been pushed to remote"
