import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <div className="min-h-screen flex items-center justify-center bg-[var(--blueColor)]">
      <h1 className="text-4xl font-bold text-[var(--greenColor)]">Tailwind funcionando</h1>
    </div>
  </StrictMode>,
)
