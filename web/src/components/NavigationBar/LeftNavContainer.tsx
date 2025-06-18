import dynamic from 'next/dynamic'

const LeftNavBar = dynamic(() => import('./LeftNavBar'), {
  ssr: false,
})

const LeftNavContainer = ({ user }: { user: string | null }) => {
  return <LeftNavBar user={user} />
}

export default LeftNavContainer
