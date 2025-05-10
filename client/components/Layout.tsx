import { useContext } from 'react'
import { AuthContext } from '../modules/auth_provider'
import Navbar from './Navbar'

interface LayoutProps {
  children: React.ReactNode
}

const Layout = ({ children }: LayoutProps) => {
  const { authenticated } = useContext(AuthContext)

  return (
    <div className="min-h-screen bg-gray-50">
      {authenticated && <Navbar />}
      <main>{children}</main>
    </div>
  )
}

export default Layout
