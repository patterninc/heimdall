import { useState } from 'react'

// Helper function to get cookie value
const getCookie = (name: string) => {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) {
    return parts.pop()?.split(';').shift() || null
  }
  return null
}

export function useUser() {
  const [user] = useState<string | null>(() => {
    // Initialize state lazily with the cookie value
    const userFromCookie = getCookie('heimdall-user')
    // Decode the URL-encoded value
    return userFromCookie ? decodeURIComponent(userFromCookie) : null
  })

  return user
}
