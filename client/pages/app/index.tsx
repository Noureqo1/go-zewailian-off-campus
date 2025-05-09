import React, { useState, useRef, useContext, useEffect } from 'react'
import ChatBody from '../../components/chat_body'
import { WebsocketContext } from '../../modules/websocket_provider'
import { useRouter } from 'next/router'
import { API_URL } from '../../constants'
import autosize from 'autosize'
import { AuthContext } from '../../modules/auth_provider'

export type Message = {
  content: string
  client_id: string
  username: string
  room_id: string
  type: 'recv' | 'self'
}

const Index = () => {
  const [messages, setMessage] = useState<Array<Message>>([])
  const textarea = useRef<HTMLTextAreaElement>(null)
  const { conn } = useContext(WebsocketContext)
  const [users, setUsers] = useState<Array<{ username: string }>>([])
  const { user } = useContext(AuthContext)

  const router = useRouter()

  useEffect(() => {
    if (conn === null) {
      router.push('/')
      return
    }

    const roomId = conn.url.split('/')[5]
    async function getUsers() {
      try {
        const res = await fetch(`${API_URL}/ws/getClients/${roomId}`, {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
        })
        const data = await res.json()

        setUsers(data)
      } catch (e) {
        console.error(e)
      }
    }
    getUsers()
  }, [])

  // Effect for autosize textarea initialization
  useEffect(() => {
    if (textarea.current) {
      autosize(textarea.current)
    }
  }, [textarea])

  // Effect for websocket message handling
  useEffect(() => {
    if (conn === null) {
      router.push('/')
      return
    }

    // Set up message handler
    const handleMessage = (message: MessageEvent) => {
      const m: Message = JSON.parse(message.data)
      
      if (m.content === 'A new user has joined the room') {
        setUsers(prevUsers => [...prevUsers, { username: m.username }])
        return
      }

      if (m.content === 'user left the chat') {
        setUsers(prevUsers => {
          const updatedUsers = prevUsers.filter(u => u.username !== m.username)
          return updatedUsers
        })
        setMessage(prevMessages => [...prevMessages, m])
        return
      }

      // Determine if message is from self or others
      const messageWithType = {...m}
      if (user?.username === m.username) {
        messageWithType.type = 'self'
      } else {
        messageWithType.type = 'recv'
      }
      
      setMessage(prevMessages => [...prevMessages, messageWithType])
    }

    // Add event listener
    conn.addEventListener('message', handleMessage)

    // Cleanup function
    return () => {
      conn.removeEventListener('message', handleMessage)
    }
  }, [conn, router, user?.username])

  const sendMessage = () => {
    if (!textarea.current?.value) return
    if (conn === null) {
      router.push('/')
      return
    }

    conn.send(textarea.current.value)
    textarea.current.value = ''
  }

  return (
    <>
      <div className='flex flex-col w-full'>
        {/* Display online users */}
        <div className='p-2 md:mx-6 bg-dark-secondary rounded-md mb-2'>
          <h3 className='text-white font-bold mb-1'>Online Users ({users.length})</h3>
          <div className='flex flex-wrap gap-2'>
            {users.map((user, index) => (
              <span key={index} className='px-2 py-1 bg-dark-primary text-white rounded-full text-sm'>
                {user.username}
              </span>
            ))}
          </div>
        </div>
        <div className='p-4 md:mx-6 mb-14'>
          <ChatBody data={messages} />
        </div>
        <div className='fixed bottom-0 mt-4 w-full'>
          <div className='flex md:flex-row px-4 py-2 bg-grey md:mx-4 rounded-md'>
            <div className='flex w-full mr-4 rounded-md border border-blue'>
              <textarea
                ref={textarea}
                placeholder='type your message here'
                className='w-full h-10 p-2 rounded-md focus:outline-none'
                style={{ resize: 'none' }}
              />
            </div>
            <div className='flex items-center'>
              <button
                className='p-2 rounded-md bg-blue text-white'
                onClick={sendMessage}
              >
                Send
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export default Index
