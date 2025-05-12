import { API_URL } from '../../constants'
import { useState, useContext } from 'react'
import { useRouter } from 'next/router'
import Image from 'next/image'
import Link from 'next/link'

import { AuthContext, UserInfo } from '../../modules/auth_provider'

const SignupPage = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const router = useRouter()
  const { setAuthenticated, setUser } = useContext(AuthContext)

  const handleGoogleSignup = () => {
    window.location.href = `${API_URL}/auth/google/login`
  }

  const submitHandler = async (e: React.SyntheticEvent) => {
    e.preventDefault()

    try {
      const res = await fetch(`${API_URL}/auth/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email, password, name }),
      })

      const data = await res.json()
      if (res.ok) {
        const user: UserInfo = {
          name: data.name,
          id: data.id,
          email: data.email,
          picture: data.picture
        }

        setAuthenticated(true)
        setUser(user)
        router.push('/')
      } else {
        // Handle signup error
        console.error('Signup failed:', data.error)
      }
    } catch (err) {
      console.error('Signup error:', err)
    }
  }

  return (
    <div className='flex items-center justify-center min-w-full min-h-screen'>
      <div className='w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md'>
        <h1 className='text-2xl font-bold text-center text-gray-900'>
          Create Account
        </h1>
        <p className='text-center text-gray-600'>Join our chat community</p>

        <form className='space-y-4' onSubmit={submitHandler}>
          <div>
            <input
              type='text'
              placeholder='Full Name'
              value={name}
              onChange={(e) => setName(e.target.value)}
              className='w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500'
              required
            />
          </div>
          <div>
            <input
              type='email'
              placeholder='Email'
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className='w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500'
              required
            />
          </div>
          <div>
            <input
              type='password'
              placeholder='Password'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className='w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500'
              required
            />
          </div>
          <button
            type='submit'
            className='w-full px-4 py-2 text-black bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
          >
            Sign up
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
          onClick={handleGoogleSignup}
          className='flex items-center justify-center w-full px-4 py-2 space-x-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
        >
          <div className='relative w-6 h-6'>
            <Image src='/google.svg' alt='Google' layout='fill' />
          </div>
          <span>Sign up with Google</span>
        </button>

        <div className='text-center text-sm text-gray-600'>
          Already have an account?{' '}
          <Link href='/login'>
            <span className='font-medium text-indigo-600 hover:text-indigo-500 cursor-pointer'>
              Sign in
            </span>
          </Link>
        </div>
      </div>
    </div>
  )
}

export default SignupPage
