import { useContext } from 'react'
import { useRouter } from 'next/router'
import { AuthContext } from '../modules/auth_provider'
import { API_URL } from '../constants'
import Image from 'next/image'

const Navbar = () => {
  const { user, setAuthenticated, setUser } = useContext(AuthContext)
  const router = useRouter()

  const handleLogout = async () => {
    try {
      await fetch('/api/auth/logout', {
        method: 'POST',
        credentials: 'include'
      })
    } catch (error) {
      console.error('Logout failed:', error)
    }
    setAuthenticated(false)
    setUser({ name: '', id: '', email: '', picture: '' })
    router.push('/login')
  }

  return (
    <nav className="bg-white shadow">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex">
            <div className="flex-shrink-0 flex items-center">
              <h1 className="text-xl font-bold text-indigo-600">Chat App</h1>
            </div>
          </div>
          <div className="flex items-center">
            {user.id && (
              <div className="flex items-center space-x-4">
                {user.picture && (
                  <div className="relative w-8 h-8 rounded-full overflow-hidden">
                    <Image
                      src={user.picture}
                      alt={user.name}
                      layout="fill"
                      objectFit="cover"
                    />
                  </div>
                )}
                <span className="text-gray-700">{user.name}</span>
                <button
                  onClick={handleLogout}
                  className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-black bg-indigo-600"
                >
                  Logout
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}

export default Navbar
