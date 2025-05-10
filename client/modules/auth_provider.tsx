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
    const userInfo = localStorage.getItem('user_info')

    if (!userInfo) {
      setAuthenticated(false)
      setUser({ name: '', id: '', email: '', picture: '' })
      return
    }

    try {
      const user: UserInfo = JSON.parse(userInfo)
      if (user && user.id) {
        setUser({
          name: user.name,
          id: user.id,
          email: user.email,
          picture: user.picture
        })
        setAuthenticated(true)
      } else {
        localStorage.removeItem('user_info')
        setAuthenticated(false)
        setUser({ name: '', id: '', email: '', picture: '' })
      }
    } catch {
      localStorage.removeItem('user_info')
      setAuthenticated(false)
      setUser({ name: '', id: '', email: '', picture: '' })
    }
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
