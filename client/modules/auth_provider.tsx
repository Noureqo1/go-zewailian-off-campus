import { useState, createContext, useEffect } from 'react'
import { useRouter } from 'next/router'

export interface UserInfo {
  name: string
  id: string
  email?: string
  picture?: string
  username?: string
}

export const AuthContext = createContext<{
  authenticated: boolean
  setAuthenticated: (auth: boolean) => void
  user: UserInfo
  setUser: (user: UserInfo) => void
}>({
  authenticated: false,
  setAuthenticated: () => {},
  user: { name: '', id: '', email: '', picture: '' },
  setUser: () => {},
})

const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false)
  const [user, setUser] = useState<UserInfo>({ name: '', id: '', email: '', picture: '' })

  const router = useRouter()

  useEffect(() => {
    const checkSession = async () => {
      try {
        const response = await fetch('/api/auth/me', {
          credentials: 'include'
        })
        
        if (response.ok) {
          const user = await response.json()
          setUser({
            name: user.name,
            id: user.id,
            email: user.email,
            picture: user.picture
          })
          setAuthenticated(true)
        } else {
          setAuthenticated(false)
          setUser({ name: '', id: '', email: '', picture: '' })
        }
      } catch {
        setAuthenticated(false)
        setUser({ name: '', id: '', email: '', picture: '' })
      }
    }

    checkSession()
  }, [])

  return (
    <AuthContext.Provider
      value={{
        authenticated: authenticated,
        setAuthenticated: setAuthenticated,
        user: user,
        setUser: setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export default AuthContextProvider
