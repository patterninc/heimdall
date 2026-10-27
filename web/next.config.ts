import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // `pnpm type-check` runs TypeScript 7 in CI instead. `next build` would check with TypeScript 6,
  // because `typescript` is aliased to @typescript/typescript6 until typescript-eslint supports 7.
  typescript: { ignoreBuildErrors: true },

  reactCompiler: true,

  async rewrites() {
    const apiUrl = process.env.API_URL || 'http://localhost:9090'

    return [
      {
        source: '/api/v1/:path*',
        destination: `${apiUrl}/api/v1/:path*`,
      },
    ]
  },

  async headers() {
    return [
      {
        // Apply these headers to all routes
        source: '/:path*',
        headers: [
          { key: 'Access-Control-Allow-Credentials', value: 'true' },
          { key: 'Access-Control-Allow-Origin', value: '*' },
          {
            key: 'Access-Control-Allow-Methods',
            value: 'GET,DELETE,PATCH,POST,PUT,OPTIONS',
          },
          {
            key: 'Access-Control-Allow-Headers',
            value:
              'X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, Authorization, X-Pattern-User, X-Heimdall-User',
          },
        ],
      },
    ]
  },
}

export default nextConfig
