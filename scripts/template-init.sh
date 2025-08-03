#!/bin/bash

set -e

echo "🚀 Initializing template-goapi..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository. Please run 'git init' first."
    exit 1
fi

# Get git remote URL
REMOTE_URL=$(git remote get-url origin 2>/dev/null || echo "")

if [[ -z "$REMOTE_URL" ]]; then
    echo "❌ Error: No git remote 'origin' found. Please add a remote first:"
    echo "   git remote add origin git@github.com:username/repo-name.git"
    exit 1
fi

echo "📡 Found git remote: $REMOTE_URL"

# Extract repo name from various git URL formats
# Handle: git@github.com:foo/bar.git, https://github.com/foo/bar.git, etc.
REPO_NAME=""
if [[ "$REMOTE_URL" =~ git@[^:]+:([^/]+)/([^/]+)\.git$ ]]; then
    # SSH format: git@github.com:user/repo.git
    REPO_NAME="${BASH_REMATCH[2]}"
elif [[ "$REMOTE_URL" =~ https?://[^/]+/([^/]+)/([^/]+)\.git$ ]]; then
    # HTTPS format: https://github.com/user/repo.git
    REPO_NAME="${BASH_REMATCH[2]}"
elif [[ "$REMOTE_URL" =~ ([^/]+)/([^/]+)$ ]]; then
    # Simple format: user/repo
    REPO_NAME="${BASH_REMATCH[2]}"
else
    echo "❌ Error: Could not parse repository name from remote URL: $REMOTE_URL"
    echo "Please ensure your remote URL follows standard git conventions."
    exit 1
fi

echo "📦 Detected repository name: $REPO_NAME"

# Extract GitHub username/org for full module name
MODULE_NAME=""
if [[ "$REMOTE_URL" =~ git@[^:]+:([^/]+)/([^/]+)\.git$ ]]; then
    # SSH format: git@github.com:user/repo.git
    MODULE_NAME="github.com/${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
elif [[ "$REMOTE_URL" =~ https?://([^/]+)/([^/]+)/([^/]+)\.git$ ]]; then
    # HTTPS format: https://github.com/user/repo.git
    MODULE_NAME="${BASH_REMATCH[1]}/${BASH_REMATCH[2]}/${BASH_REMATCH[3]}"
elif [[ "$REMOTE_URL" =~ ([^/]+)/([^/]+)$ ]]; then
    # Simple format: user/repo - assume GitHub
    MODULE_NAME="github.com/${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
else
    echo "❌ Error: Could not parse module name from remote URL: $REMOTE_URL"
    exit 1
fi

echo "🔧 Module name: $MODULE_NAME"

# Convert repo name to uppercase for PROJECT_NAME
PROJECT_NAME=$(echo "$REPO_NAME" | tr '[:lower:]' '[:upper:]' | tr '-' '_')
echo "📛 Project name: $PROJECT_NAME"

# Create config directory name from repo name (lowercase, with dot prefix)
CONFIG_DIR=".$(echo "$REPO_NAME" | tr '[:upper:]' '[:lower:]')"
echo "📁 Config directory: $CONFIG_DIR"

echo ""
echo "🔄 Replacing template variables in all files..."

# Find all files (excluding .git, node_modules, vendor, etc.)
find . -type f \
    -not -path "./.git/*" \
    -not -path "./node_modules/*" \
    -not -path "./vendor/*" \
    -not -path "./.vscode/*" \
    -not -path "./.idea/*" \
    -not -name "*.log" \
    -not -name "*.tmp" \
    -not -name "template-init.sh" \
    | while read -r file; do
    
    # Check if file contains any template variables
    if grep -l "{{MODULE_NAME}}\|{{PROJECT_NAME}}\|{{CONFIG_DIR}}" "$file" > /dev/null 2>&1; then
        echo "  📝 Processing: $file"
        
        # Use sed to replace template variables
        sed -i.bak \
            -e "s|{{MODULE_NAME}}|$MODULE_NAME|g" \
            -e "s|{{PROJECT_NAME}}|$PROJECT_NAME|g" \
            -e "s|{{CONFIG_DIR}}|$CONFIG_DIR|g" \
            "$file"
        
        # Remove backup file
        rm "${file}.bak"
    fi
done

echo ""
echo "✅ Template initialization complete!"
echo ""
echo "🎯 Summary:"
echo "   Repository: $REPO_NAME"
echo "   Module: $MODULE_NAME" 
echo "   Project: $PROJECT_NAME"
echo "   Config dir: ~/$CONFIG_DIR"
echo ""
echo "🚀 Next steps:"
echo "   1. Run: go mod tidy"
echo "   2. Run: go run cmd/server/main.go"
echo "   3. Test: curl http://localhost:8080/health"
echo ""
echo "🗑️  You can now delete this script: rm scripts/template-init.sh"