"use client"

import { useState, useEffect } from "react"
import "./App.css"
// Corrected import path: go up one level from 'src' to 'frontend', then into 'wailsjs'
import { InstallCLI } from "../wailsjs/go/main/App"

function App() {
  const [installPath, setInstallPath] = useState("")
  const [status, setStatus] = useState("")
  const [error, setError] = useState("")
  const [installing, setInstalling] = useState(false)
  const [osName, setOsName] = useState("")

  useEffect(() => {
    // Set default install path based on OS
    if (window.runtime) {
      setOsName(window.runtime.GOOS)
      if (window.runtime.GOOS === "windows") {
        // Default to USERPROFILE\bin for Windows
        setInstallPath(`${window.runtime.USERPROFILE}\\bin`)
      } else {
        // Default to /usr/local/bin for Linux/macOS
        setInstallPath("/usr/local/bin")
      }
    }
  }, [])

  const handleInstall = async () => {
    setInstalling(true)
    setStatus("Starting installation...")
    setError("")
    try {
      // Call the Go backend function
      const result = await InstallCLI(installPath)
      setStatus(result)
    } catch (err) {
      setError(`Installation failed: ${err}`)
      setStatus("")
    } finally {
      setInstalling(false)
    }
  }

  return (
    <div id="App">
      <div className="container">
        <h1>JDK Manager Installer</h1>
        <p>Install the JDK Manager CLI tool on your system.</p>

        <div className="input-group">
          <label htmlFor="installPath">Installation Directory:</label>
          <input
            id="installPath"
            type="text"
            value={installPath}
            onChange={(e) => setInstallPath(e.target.value)}
            placeholder={osName === "windows" ? "%USERPROFILE%\\bin" : "/usr/local/bin"}
            disabled={installing}
          />
        </div>

        <button onClick={handleInstall} disabled={installing}>
          {installing ? "Installing..." : "Install JDK Manager"}
        </button>

        {status && <pre className="status-message success">{status}</pre>}
        {error && <pre className="status-message error">{error}</pre>}

        <p className="note">
          After installation, please restart your terminal to use the `jdk` command. On Linux/macOS, you might be
          prompted for your sudo password.
        </p>
      </div>
    </div>
  )
}

export default App
