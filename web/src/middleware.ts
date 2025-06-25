import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  // Get the user from the incoming headers
  const user = request.headers.get('X-Heimdall-User')

  // Clone the request headers
  const requestHeaders = new Headers(request.headers)

  // Create the response
  const response = NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  })

  // If user exists, you can set it as a cookie or pass it differently
  if (user) {
    response.cookies.set('heimdall-user', user, {
      httpOnly: true,
      sameSite: 'strict',
      path: '/',
    })
  }

  return response
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}
