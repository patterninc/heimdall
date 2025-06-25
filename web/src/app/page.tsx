'use client'

import dynamic from 'next/dynamic'
import { Suspense } from 'react'

const Jobs = dynamic(() => import('../modules/Jobs/Jobs'), {
  ssr: false,
})

export default function Index() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Jobs />
    </Suspense>
  )
}
