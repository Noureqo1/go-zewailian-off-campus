import '../styles/globals.css'
import type { AppProps } from 'next/app'
import AuthContextProvider from '../modules/auth_provider'
import WebSocketProvider from '../modules/websocket_provider'
import Layout from '../components/Layout'

export default function App({ Component, pageProps }: AppProps) {
  return (
    <AuthContextProvider>
      <WebSocketProvider>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </WebSocketProvider>
    </AuthContextProvider>
  )
}
