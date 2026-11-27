// React `useId` values (`_r_7_`) depend on how many roots rendered earlier in the file, so a
// snapshot would change when tests are filtered or reordered. Number them by first appearance.
const REACT_ID = /_r_[0-9a-z]+_/g

/** Returns a copy of `element` with React-generated ids renamed in document order, for snapshots. */
export const withStableIds = <T extends Element>(element: T): T => {
  const clone = element.cloneNode(true) as T
  const ids = new Map<string, string>()
  const rename = (id: string) => {
    if (!ids.has(id)) ids.set(id, `_id${ids.size}_`)
    return ids.get(id)!
  }
  for (const node of [clone, ...clone.querySelectorAll('*')]) {
    for (const attr of [...node.attributes]) {
      const value = attr.value.replace(REACT_ID, rename)
      if (value !== attr.value) node.setAttribute(attr.name, value)
    }
  }
  return clone
}
