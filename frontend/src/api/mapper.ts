import { FalcoEvent, Priority, RawEvent } from './model'

export const mapEvent = (input: RawEvent[]): FalcoEvent[] => input.map<FalcoEvent>((event) => {
  const outputFields = Object.entries(event.output_fields || {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
    acc[key] = new String(value).toString()

    return acc
  }, {})

  const sorted = Object.entries(outputFields)

  sorted.sort(([ak, av], [bk, bv]) => (ak.length + av.length) - (bk.length + bv.length))

  return {
    output: event.output,
    priority: event.priority.toLowerCase() as Priority,
    rule: event.rule,
    time: new Date(event.time),
    outputFields: Object.fromEntries(sorted),
    filterTime: Math.floor((new Date(event.time)).getTime() / 1000) * 1000
  }
}).reverse()

export const initStats = () => ({
  [Priority.EMERGENCY]: 0,
  [Priority.ALERT]: 0,
  [Priority.CRITICAL]: 0,
  [Priority.ERROR]: 0,
  [Priority.WARNING]: 0,
  [Priority.NOTICE]: 0,
  [Priority.INFORMATIONAL]: 0,
  [Priority.DEBUG]: 0,
  [Priority.NONE]: 0
})

const priorityColorMap: { [key in Priority]: string } = {
  [Priority.NONE]: 'grey',
  [Priority.DEBUG]: 'light-blue lighten-1',
  [Priority.INFORMATIONAL]: 'light-blue',
  [Priority.NOTICE]: 'primary',
  [Priority.WARNING]: 'warning',
  [Priority.ERROR]: 'error',
  [Priority.CRITICAL]: 'red darken-1',
  [Priority.ALERT]: 'red darken-2',
  [Priority.EMERGENCY]: 'red darken-3'
}

export const colorCodes: { [key in Priority]: string } & { total: string } = {
  total: '#757575',
  [Priority.EMERGENCY]: '#C62828',
  [Priority.ALERT]: '#D32F2F',
  [Priority.CRITICAL]: '#E53935',
  [Priority.ERROR]: '#FF5252',
  [Priority.WARNING]: '#FB8C00',
  [Priority.NOTICE]: '#1976D2',
  [Priority.INFORMATIONAL]: '#03A9F4',
  [Priority.DEBUG]: '#29B6F6',
  [Priority.NONE]: '#555'
}

export const mapPriorityToColor = (priority: Priority): string => priorityColorMap[priority]
