'use client'
import dynamic from 'next/dynamic'

const Jobs = dynamic(() => import('../modules/Jobs/Jobs'), {
  ssr: false,
})

export default function Index() {
  return <Jobs />
}
