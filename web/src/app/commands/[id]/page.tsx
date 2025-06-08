'use client'
import dynamic from 'next/dynamic'

import { use } from 'react'
const CommandDetails = dynamic(
  () =>
    import('../../../modules/Commands/CommandDetails/CommandDetails').then(
      (mod) => mod.CommandDetails,
    ),
  {
    ssr: false,
  },
)

const CommandDetailLayout = (props: { params: Promise<{ id: string }> }) => {
  const params = use(props.params)
  const { id } = params

  return <CommandDetails id={id} />
}
export default CommandDetailLayout
