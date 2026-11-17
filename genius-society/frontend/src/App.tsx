import { BrowserRouter } from 'react-router-dom'
import AppRoutes from './routes'
import { AuthProvider } from './contexts/AuthContext'
import { ToastProvider } from './components/common/Toast'

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <ToastProvider>
          <AppRoutes />
        </ToastProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}
