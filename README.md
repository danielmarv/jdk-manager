# JDK Manager

A cross-platform command-line tool for managing multiple Java Development Kit (JDK) versions, similar to nvm for Node.js.

## 🚀 Features

- **Easy Installation**: Download and install JDKs from Eclipse Adoptium with a single command
- **Version Switching**: Seamlessly switch between different JDK versions
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Environment Management**: Automatically handles JAVA_HOME and PATH configuration
- **Progress Tracking**: Visual progress bars for downloads
- **Version Listing**: List installed and available JDK versions
- **LTS Support**: Easy identification of Long Term Support versions

## 📦 Installation

To install JDK Manager and make the `jdk` command available globally in your terminal, follow the instructions for your operating system.

### **Linux / macOS**

1.  **Navigate to the project directory:**
    ```bash
    cd /path/to/jdk-manager
    ```
2.  **Run the installer script:**
    ```bash
    chmod +x scripts/install.sh
    sudo scripts/install.sh
    ```
    The script will build the `jdk` executable, copy it to `/usr/local/bin`, and set up a shell function in your `~/.bashrc` or `~/.zshrc` for seamless `jdk use` functionality.
3.  **Restart your terminal** or run `source ~/.bashrc` (or `~/.zshrc`) to activate the `jdk` command.

### **Windows (PowerShell)**

1.  **Navigate to the project directory:**
    ```powershell
    cd D:\projects\jdk-manager # Adjust to your actual project path
    ```
2.  **Run the installer script:**
    ```powershell
    .\scripts\install.ps1
    ```
    The script will build the `jdk.exe` executable, copy it to `%USERPROFILE%\bin`, add this directory to your user's PATH, and set up a PowerShell function in your `$PROFILE` for seamless `jdk use` functionality.
3.  **Restart your PowerShell terminal** to activate the `jdk` command.

## 🔧 Usage

Once installed, you can use the `jdk` command directly from any directory.

### List Available JDK Versions

```bash
jdk list-remote
jdk list-remote --lts
jdk list-remote --all
```

### Install a JDK Version

```bash
jdk install 21
jdk install 17.0.8
jdk install 21 --force
```

### List Installed Versions

```bash
jdk list
```

### Switch JDK Version

```bash
jdk use 21
jdk use 17.0.8
```
After running `jdk use`, you will see a confirmation message and the `java -version` output for the newly active JDK.

### Get Help

```bash
jdk --help
jdk install --help
```

## 📁 Directory Structure

JDKs are installed in \`~/.jdks/\` directory:

```
~/.jdks/
├── 21/
│   ├── bin/
│   ├── lib/
│   └── ...
├── 17.0.8/
│   ├── bin/
│   ├── lib/
│   └── ...
└── 11/
    ├── bin/
    ├── lib/
    └── ...
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or later
- Internet connection for downloading JDKs

### Building

```bash
go build -o jdk main.go
```

### Running Tests

```bash
go test ./...
```

### Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Progress Bar](https://github.com/schollz/progressbar) - Download progress visualization
- [Go Homedir](https://github.com/mitchellh/go-homedir) - Cross-platform home directory detection

## 🗺️ Roadmap

### Phase 1: Enhanced Core Features
- [ ] **Project-based JDK switching**: Support \`.jdk-version\` file in project directories
- [ ] **Shell integration**: Automatic JDK switching when entering directories
- [ ] **Aliases**: Create and manage JDK aliases (\`jdk alias default 21\`)
- [ ] **Configuration file**: User preferences and default settings

### Phase 2: Multi-Distribution Support
- [ ] **Amazon Corretto**: Support for Amazon's OpenJDK distribution
- [ ] **Eclipse Temurin**: Enhanced Adoptium/Temurin support
- [ ] **Azul Zulu**: Support for Azul's JDK builds
- [ ] **Oracle JDK**: Support for Oracle's official JDK (where licensing permits)
- [ ] **GraalVM**: Support for GraalVM distributions

### Phase 3: Platform Enhancement
- [ ] **Windows Full Support**: Native PowerShell integration and .bat hooks
- [ ] **Shell completions**: Bash, Zsh, Fish, and PowerShell completions
- [ ] **Homebrew formula**: Easy installation on macOS
- [ ] **Package managers**: Support for apt, yum, chocolatey

### Phase 4: Advanced Features
- [ ] **GUI Mode**: Optional desktop application for non-CLI users
- [ ] **JDK verification**: Checksum verification for downloaded JDKs
- [ ] **Proxy support**: Corporate proxy configuration
- [ ] **Offline mode**: Use local JDK archives

### Phase 5: Team & Enterprise
- [ ] **Team synchronization**: Share JDK versions across development teams
- [ ] **Private mirrors**: Support for internal JDK repositories
- [ ] **Policy enforcement**: Restrict allowed JDK versions
- [ ] **Audit logging**: Track JDK installations and usage

### Phase 6: Ecosystem Integration
- [ ] **Plugin system**: Extensible architecture for third-party plugins
- [ ] **Maven integration**: Automatic JDK switching based on Maven projects
- [ ] **Gradle integration**: Gradle project JDK detection and switching
- [ ] **IDE integration**: Plugins for IntelliJ IDEA, Eclipse, VS Code
- [ ] **Docker support**: Generate Dockerfiles with specific JDK versions

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. Fork the repository
2. Create your feature branch (\`git checkout -b feature/amazing-feature\`)
3. Commit your changes (\`git commit -m 'Add some amazing feature'\`)
4. Push to the branch (\`git push origin feature/amazing-feature\`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Eclipse Adoptium](https://adoptium.net/) for providing free, high-quality JDK builds
- The Go community for excellent tooling and libraries

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/your-username/jdk-manager/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/your-username/jdk-manager/discussions)
- 📖 **Documentation**: [Wiki](https://github.com/your-username/jdk-manager/wiki)

---

**Made with ❤️ for the Java community**
