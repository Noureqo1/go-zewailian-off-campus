import { useState, useContext, useEffect } from 'react'
import { API_URL } from '../../constants'
import { useRouter } from 'next/router'
import { AuthContext, UserInfo } from '../../modules/auth_provider'
import { FcGoogle } from 'react-icons/fc'
import Link from 'next/link'

const LoginPage = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const { authenticated, setAuthenticated, setUser } = useContext(AuthContext)

  const router = useRouter()

  useEffect(() => {
    if (authenticated) {
      router.push('/')
      return
    }

    const urlParams = new URLSearchParams(window.location.search)
    const success = urlParams.get('success')
    if (success === 'true') {
      // Get user info after successful Google OAuth
      fetch(`${API_URL}/auth/me`, {
        credentials: 'include'
      })
        .then(res => res.json())
        .then(data => {
          if (!data.error) {
            setUser(data)
            setAuthenticated(true)
            setUser(data)
            setAuthenticated(true)
            router.push('/')
          }
        })
        .catch(err => console.error('Error fetching user info:', err))
    }
  }, [authenticated, router, setUser, setAuthenticated])

  const submitHandler = async (e: React.SyntheticEvent) => {
    e.preventDefault()

    try {
      const res = await fetch(`${API_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      })

      const data = await res.json()
      if (res.ok) {
        const user: UserInfo = {
          name: data.name,
          id: data.id,
          email: data.email,
          picture: data.picture
        }

        setUser(user)
        setAuthenticated(true)
        return router.push('/')
      }
    } catch (err) {
      console.log(err)
    }
  }

  const handleGoogleLogin = () => {
    const popup = window.open(
      `${API_URL}/auth/google/login`,
      'Google Login',
      'width=500,height=600'
    )

    window.addEventListener('message', (event) => {
      if (event.origin !== 'http://localhost:8080') return
      if (event.data.type === 'oauth_success') {
        const { user } = event.data
        setUser(user)
        setAuthenticated(true)
        localStorage.setItem('user_info', JSON.stringify(user))
        router.push('/')
        popup?.close()
      }
    })
  }

  return (
    <div className='flex items-center justify-center min-w-full min-h-screen'>
      <div className='w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md'>
        <h1 className='text-2xl font-bold text-center text-gray-900'>
          Welcome Back
        </h1>
        <p className='text-center text-gray-600'>Sign in to continue to Chat</p>

        <form className='space-y-4' onSubmit={submitHandler}>
          <div>
            <input
              type='email'
              placeholder='Email'
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className='w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500'
            />
          </div>
          <div>
            <input
              type='password'
              placeholder='Password'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className='w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500'
            />
          </div>
          <button
            type='submit'
            className='w-full px-4 py-2 text-black bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
          >
            Sign in
          </button>
        </form>

        <div className='relative'>
          <div className='absolute inset-0 flex items-center'>
            <div className='w-full border-t border-gray-300'></div>
          </div>
          <div className='relative flex justify-center text-sm'>
            <span className='px-2 text-gray-500 bg-white'>Or continue with</span>
          </div>
        </div>

        <button
          onClick={handleGoogleLogin}
          className='flex items-center justify-center w-full px-4 py-2 space-x-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
        >
          <FcGoogle className='w-6 h-6' aria-hidden='true' />
          <span>Sign in with Google</span>
        </button> 


        <div className='text-center text-sm text-gray-600'>
          Don&apos;t have an account?{' '}
          <Link href='/signup'>
            <span className='font-medium text-black hover:text-gray-700 cursor-pointer'>
              Sign up
            </span>
          </Link>
        </div>
      </div>
    </div>
  )
}

export default LoginPage
