import axios from 'axios'
import { Config, EventResponse } from './model'
import { mapEvent } from './mapper'

const production = process.env.NODE_ENV === 'production'

const instance = axios.create({ baseURL: production ? `//${window.location.host}/` : process.env.VUE_APP_API })

export default {
  async event (): Promise<EventResponse> {
    const { data } = await instance.get('events')

    return ({
      events: mapEvent(data.events || []),
      retention: data.retention,
      stats: data.stats,
      outputs: data.outputs
    })
  },
  async config (): Promise<Config> {
    const { data } = await instance.get('config')

    return data
  }
}
