import dynamic from 'next/dynamic'
import { headers } from 'next/headers'

const LeftNavBar = dynamic(() => import('./LeftNavBar'), {
  ssr: false,
})

const LeftNavContainer = () => {
  const header = headers()
  const user = header.get('X-Heimdall-User')
  return <LeftNavBar user={user} />
}

export default LeftNavContainer
