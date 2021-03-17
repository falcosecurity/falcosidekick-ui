import api from '@/api'
import { EventResponse, FalcoEvent, Priority, SocketMessage } from '@/api/model'
import Vue from 'vue'
import Vuex from 'vuex'
import { mapEvent } from '../api/mapper'
import debounce from 'lodash.debounce'

Vue.use(Vuex)

export const FETCH_EVENTS = 'FETCH_EVENTS'
export const RECEIVE_EVENTS = 'RECEIVE_EVENTS'
export const SET_STATE = 'SET_STATE'

export type State = {
  events: FalcoEvent[];
  outputs: string[];
  retention: number;
  socket: { isConnected: boolean };
}

export default new Vuex.Store<State>({
  state: {
    socket: {
      isConnected: false
    },
    events: [],
    outputs: [],
    retention: 0
  },
  mutations: {
    [SET_STATE]: (state, response: EventResponse) => {
      state.events = response.events || []
      state.outputs = response.outputs || []
      state.retention = response.retention
    }
  },
  actions: {
    [FETCH_EVENTS]: ({ commit }) => {
      api.event().then(response => commit(SET_STATE, response))
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
