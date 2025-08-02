import React from "react"
import ReactDOM from "react-dom/client"
import App from "./App.jsx" // Import your main App component
import "./App.css" // Assuming you might have a global CSS file

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
