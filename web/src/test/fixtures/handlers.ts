import { http, HttpResponse } from 'msw'

import { API_URL } from '@/common/Services'
import { clusterStatuses, clusters } from './clusters'
import { commandStatuses, commands } from './commands'
import { jobStatuses, jobs, jobsPage } from './jobs'

/** Same-origin API base the app fetches from; MSW matches it against jsdom's origin. */
export const api = (path: string) => `${API_URL}${path}`

/** Mirrors the API's `?id=` filter on the commands and clusters list endpoints. */
const filterById = <T extends { id: string }>(rows: T[], request: Request) => {
  const id = new URL(request.url).searchParams.get('id')
  return id ? rows.filter((row) => row.id === id) : rows
}

/** Default happy-path handlers. Override per test with `server.use(...)`. */
export const handlers = [
  http.get(api('/jobs'), () => HttpResponse.json(jobsPage())),
  http.get(api('/job/statuses'), () => HttpResponse.json(jobStatuses)),
  http.get(api('/job/:id'), ({ params }) => {
    const job = jobs.find((row) => row.id === params.id)
    return job
      ? HttpResponse.json(job)
      : HttpResponse.json({ error: 'not found' }, { status: 404 })
  }),
  http.post(api('/job/:id/cancel'), () =>
    HttpResponse.json({ status: 'CANCELING' }),
  ),

  http.get(api('/commands'), ({ request }) =>
    HttpResponse.json({ data: filterById(commands, request) }),
  ),
  http.get(api('/command/statuses'), () => HttpResponse.json(commandStatuses)),

  http.get(api('/clusters'), ({ request }) =>
    HttpResponse.json({ data: filterById(clusters, request) }),
  ),
  http.get(api('/cluster/statuses'), () => HttpResponse.json(clusterStatuses)),
]
