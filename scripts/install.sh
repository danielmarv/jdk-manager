#!/bin/bash

# Cross-platform installer script for JDK Manager (Linux/macOS)

set -e

INSTALL_DIR="/usr/local/bin" # Standard location for executables
JDK_MANAGER_HOME="$(pwd)"    # Current directory of the project

echo "--- JDK Manager Installer (Linux/macOS) ---"

# 1. Build the Go binary
echo "Building JDK Manager binary..."
make build

# Check if build was successful
if [ ! -f "${JDK_MANAGER_HOME}/dist/jdk" ]; then
    echo "Error: Build failed. 'dist/jdk' not found."
    exit 1
fi

# 2. Copy the binary to the installation directory
echo "Copying 'jdk' to ${INSTALL_DIR}..."
sudo cp "${JDK_MANAGER_HOME}/dist/jdk" "${INSTALL_DIR}/jdk"
sudo chmod +x "${INSTALL_DIR}/jdk"
echo "JDK Manager executable installed to ${INSTALL_DIR}/jdk"

# 3. Set up shell function for seamless 'jdk use'
echo "Setting up shell function for seamless 'jdk use'..."

SHELL_PROFILE=""
case "$SHELL" in
    */bash*)
        SHELL_PROFILE="$HOME/.bashrc"
        ;;
    */zsh*)
        SHELL_PROFILE="$HOME/.zshrc"
        ;;
    */fish*)
        SHELL_PROFILE="$HOME/.config/fish/config.fish"
        ;;
    *)
        echo "Warning: Unsupported shell ($SHELL). Please manually add the 'jdk' function to your profile."
        ;;
esac

if [ -n "$SHELL_PROFILE" ]; then
    # Define the shell function
    JDK_FUNCTION_DEFINITION=$(cat <<EOF
# JDK Manager setup
# This function allows 'jdk use <version>' to modify the current shell's environment.
function jdk() {
  local command="\$1"
  shift
  if [ "\$command" = "use" ]; then
    eval "\$("${INSTALL_DIR}/jdk" use "\$@")"
    echo "✓ JDK is now active!"
    java -version # Verify the change
  else
    "${INSTALL_DIR}/jdk" "\$command" "\$@"
  fi
}
EOF
)

    # Check if the function already exists to prevent duplicates
    if ! grep -q "# JDK Manager setup" "$SHELL_PROFILE"; then
        echo "$JDK_FUNCTION_DEFINITION" >> "$SHELL_PROFILE"
        echo "JDK Manager function added to $SHELL_PROFILE"
    else
        echo "JDK Manager function already exists in $SHELL_PROFILE. Skipping."
    fi
fi

echo ""
echo "--- Installation Complete! ---"
echo "Please restart your terminal or run 'source $SHELL_PROFILE' (for Bash/Zsh) to activate 'jdk' command."
echo "You can now use 'jdk install <version>' and 'jdk use <version>' directly."
