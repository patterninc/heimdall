'use client'
import dynamic from 'next/dynamic'

import { useEffect, useState } from 'react'
const CommandDetails = dynamic(
  () =>
    import('../../../modules/Commands/CommandDetails/CommandDetails').then(
      (mod) => mod.CommandDetails,
    ),
  {
    ssr: false,
  },
)

const CommandDetailLayout = ({
  params,
}: {
  params: Promise<{ id: string }>
}) => {
  const [id, setId] = useState<string>('')

  useEffect(() => {
    params.then((resolvedParams) => {
      setId(resolvedParams.id)
    })
  }, [params])

  return <CommandDetails id={id} />
}
export default CommandDetailLayout
