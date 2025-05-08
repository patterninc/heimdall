import dynamic from 'next/dynamic'

const CommandDetails = dynamic(
  () =>
    import('../../../modules/Commands/CommandDetails/CommandDetails').then(
      (mod) => mod.CommandDetails,
    ),
  {
    ssr: false,
  },
)

const CommandDetailLayout = ({ params }: { params: { id: string } }) => {
  const { id } = params

  return <CommandDetails id={id} />
}
export default CommandDetailLayout
