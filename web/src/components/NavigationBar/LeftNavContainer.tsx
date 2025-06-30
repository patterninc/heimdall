import dynamic from 'next/dynamic'
import { useUser } from '@/common/hooks/useUser'

const LeftNavBar = dynamic(() => import('./LeftNavBar'), {
  ssr: false,
})

const LeftNavContainer = () => {
  const user = useUser()
  return <LeftNavBar user={user} />
}

export default LeftNavContainer
