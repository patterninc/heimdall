import { Section, type SectionProps } from '@patterninc/pattern-ui/section'
import React from 'react'

/**
 * Static `Section` whose header keeps its tinted background. Section's header highlights on hover
 * as if it were clickable; ours aren't (not collapsible), so pin the same tint and drop the hover.
 */
const DetailSection = ({
  className,
  ...props
}: Omit<SectionProps, 'variant'>): React.JSX.Element => (
  <Section
    {...props}
    className={`[&>div:first-child]:bg-depth-black-4/30 ${className ?? ''}`}
  />
)

export default DetailSection
