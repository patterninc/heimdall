import { server } from './fixtures/server'

/** Records the query string of every request MSW sees whose path ends with `pathSuffix`. */
export const recordRequests = (pathSuffix: string) => {
  const seen: URLSearchParams[] = []
  const listener = ({ request }: { request: Request }) => {
    const url = new URL(request.url)
    if (url.pathname.endsWith(pathSuffix)) seen.push(url.searchParams)
  }
  server.events.on('request:start', listener)
  return {
    all: () => seen,
    last: () => seen.at(-1),
    stop: () => server.events.removeListener('request:start', listener),
  }
}
