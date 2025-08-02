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

### From Source

```bash
git clone https://github.com/your-username/jdk-manager.git
cd jdk-manager
go build -o jdk main.go
```

### Using Go Install

```bash
go install github.com/your-username/jdk-manager@latest
```

## 🔧 Usage

### List Available JDK Versions

```bash
# List all available versions
jdk list-remote

# List only LTS versions
jdk list-remote --lts

# List all versions including pre-releases
jdk list-remote --all
```

### Install a JDK Version

```bash
# Install latest JDK 21
jdk install 21

# Install specific version
jdk install 17.0.8

# Force reinstall
jdk install 21 --force
```

### List Installed Versions

```bash
jdk list
```

### Switch JDK Version

```bash
# Switch to JDK 21
jdk use 21

# Switch to specific version
jdk use 17.0.8
```

The \`use\` command will output the necessary environment variable commands. Run them to activate the JDK:

**Linux/macOS:**
```bash
export JAVA_HOME="/home/user/.jdks/21"
export PATH="$JAVA_HOME/bin:$PATH"
```

**Windows PowerShell:**
```powershell
$env:JAVA_HOME = "C:\\Users\\user\\.jdks\\21"
$env:PATH = "$env:JAVA_HOME\\bin;$env:PATH"
```

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
- [nvm](https://github.com/nvm-sh/nvm) for inspiration on version management UX
- The Go community for excellent tooling and libraries

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/your-username/jdk-manager/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/your-username/jdk-manager/discussions)
- 📖 **Documentation**: [Wiki](https://github.com/your-username/jdk-manager/wiki)

---

**Made with ❤️ for the Java community**
```

```text file="LICENSE"
MIT License

Copyright (c) 2024 JDK Manager

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
