import dynamic from 'next/dynamic'

const Commands = dynamic(() => import('../../modules/Commands/Commands'), {
  ssr: false,
})

const CommandsLayput = () => {
  return <Commands />
}

export default CommandsLayput
