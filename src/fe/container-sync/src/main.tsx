import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import ButtonUsage from './ButtonUsage.tsx'
import FreeSolo from './FreeSolo.tsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
    {/* <ButtonUsage /> */}
    {/* <FreeSolo /> */}
  </React.StrictMode>,
)
