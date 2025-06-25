import { useEffect, useState } from 'react'

export function useUser() {
  const [user, setUser] = useState<string | null>(null)

  useEffect(() => {
    // Read user from cookie
    const getCookie = (name: string) => {
      const value = `; ${document.cookie}`
      const parts = value.split(`; ${name}=`)
      if (parts.length === 2) {
        return parts.pop()?.split(';').shift() || null
      }
      return null
    }

    const userFromCookie = getCookie('heimdall-user')
    setUser(userFromCookie)
  }, [])

  return user
}
