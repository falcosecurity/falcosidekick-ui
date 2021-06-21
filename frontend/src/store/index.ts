import api from '@/api'
import { EventResponse, FalcoEvent, Priority, SocketMessage } from '@/api/model'
import Vue from 'vue'
import Vuex from 'vuex'
import { mapEvent } from '../api/mapper'
import debounce from 'lodash.debounce'
import { DisplayMode } from '@/types'

Vue.use(Vuex)

export const FETCH_EVENTS = 'FETCH_EVENTS'
export const FETCH_CONFIG = 'FETCH_CONFIG'
export const RECEIVE_EVENTS = 'RECEIVE_EVENTS'
export const SET_STATE = 'SET_STATE'
export const SET_DISPLAY_MODE = 'SET_DISPLAY_MODE'
export const INIT_DISPLAY_MODE = 'INIT_DISPLAY_MODE'

export type State = {
  events: FalcoEvent[];
  outputs: string[];
  retention: number;
  socket: { isConnected: boolean };
  displayMode: DisplayMode;
}

export default new Vuex.Store<State>({
  state: {
    socket: {
      isConnected: false
    },
    events: [],
    outputs: [],
    retention: 0,
    displayMode: DisplayMode.LIGHT
  },
  mutations: {
    [SET_STATE]: (state, response: EventResponse) => {
      state.events = response.events || []
      state.outputs = response.outputs || []
      state.retention = response.retention
    },
    [SET_DISPLAY_MODE]: (state, displayMode: DisplayMode) => {
      state.displayMode = displayMode

      sessionStorage.setItem('displayMode', displayMode)
    },
    [INIT_DISPLAY_MODE]: (state, displayMode: DisplayMode) => {
      state.displayMode = displayMode
    }
  },
  actions: {
    [FETCH_EVENTS]: ({ commit }) => {
      api.event().then(response => commit(SET_STATE, response))
    },
    [FETCH_CONFIG]: async ({ commit }) => {
      const mode = sessionStorage.getItem('displayMode') as DisplayMode | ''
      if (mode) {
        return commit(INIT_DISPLAY_MODE, mode)
      }

      try {
        const config = await api.config()
        commit(INIT_DISPLAY_MODE, config.displayMode)
      } catch (error) {
        console.warn(error)
      }
    },
    [RECEIVE_EVENTS]: debounce(({ commit }, message: SocketMessage) => {
      const response: EventResponse = {
        events: mapEvent(message.events || []),
        outputs: message.outputs,
        retention: message.retention
      }

      if (message.stats) {
        response.stats = {
          [Priority.EMERGENCY]: message.stats[Priority.EMERGENCY],
          [Priority.ALERT]: message.stats[Priority.ALERT],
          [Priority.CRITICAL]: message.stats[Priority.CRITICAL],
          [Priority.ERROR]: message.stats[Priority.ERROR],
          [Priority.WARNING]: message.stats[Priority.WARNING],
          [Priority.NOTICE]: message.stats[Priority.NOTICE],
          [Priority.INFORMATIONAL]: message.stats[Priority.INFORMATIONAL],
          [Priority.DEBUG]: message.stats[Priority.DEBUG],
          [Priority.NONE]: message.stats[Priority.NONE]
        }
      }

      commit(SET_STATE, response)
    }, 1000)
  }
})
