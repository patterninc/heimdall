import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  // Get the user from the incoming headers - handle case sensitivity
  let user = request.headers.get('X-Heimdall-User')
  const version = request.headers.get('X-Heimdall-Version')

  // If not found with proper casing, try lowercase
  if (!user) {
    user = request.headers.get('x-heimdall-user')
  }

  // Clone the request headers
  const requestHeaders = new Headers(request.headers)

  // Create the response
  const response = NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  })
  // Set the cookie if user exists (backend already provides default if needed)
  if (user) {
    response.cookies.set('heimdall-user', user, {
      httpOnly: false, // Make it readable by JavaScript
      sameSite: 'strict',
      path: '/',
      secure: process.env.NODE_ENV === 'production',
    })
  }

  if (version) {
    response.cookies.set('heimdall-version', version, {
      httpOnly: false,
      sameSite: 'strict',
      path: '/',
      secure: process.env.NODE_ENV === 'production',
    })
  }

  return response
}

export const config = {
  matcher: [
    // Match all routes except static files and API routes
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}
