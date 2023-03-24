<template>
  <v-container fluid fill-height>
          <v-layout flex align-center justify-center>
            <v-flex sm4>
              <v-card>
                <v-card-text>
                  <div>
                      <v-form ref="form">
                        <v-text-field
                          label="Login"
                          v-model="username"
                          required
                        ></v-text-field>
                        <v-text-field
                          label="Password"
                          v-model="password"
                          required
                          type="password"
                        ></v-text-field>
                      </v-form>
                      <v-layout justify-space-between>
                        <v-btn @click="authenticate" class="blue darken-2 white--text">Login</v-btn>
                      </v-layout>
                      <v-alert v-if="failedAuth === true"
                        class="mt-5 mb-1"
                        outlined
                        dense
                        dark
                        type="error"
                      >{{failMsg}}</v-alert>
                  </div>
                </v-card-text>
              </v-card>
            </v-flex>
          </v-layout>
       </v-container>
</template>

<script>
import { mapActions } from 'vuex';
import { requests } from '../http';
import router from '../router';

export default {
  data() {
    return {
      username: '',
      password: '',
      failedAuth: false,
      failMsg: '',
    };
  },
  methods: {
    ...mapActions([
      'setCredentials',
    ]),
    authenticate() {
      this.failedAuth = false;
      this.failMsg = '';
      if (this.username === '' || this.password === '') {
        this.failedAuth = true;
        this.failMsg = 'login and password can\'t be empty';
        return;
      }
      requests.authenticate(
        this.username,
        this.password,
      )
        .catch((error) => {
          this.failedAuth = true;
          this.failMsg = error;
        })
        .then((response) => {
          if (response.status === 200) {
            const payload = {
              username: this.username,
              password: this.password,
            };
            this.setCredentials(payload);
            router.push('/dashboard');
          }
        });
    },
    testlogin() {
      requests.authenticate(
        'anonymous',
        'anonymous',
      )
        .then((response) => {
          if (response.status === 200) {
            const payload = {
              username: 'anonymous',
              password: 'anonymous',
            };
            this.setCredentials(payload);
            router.push('/dashboard');
          }
        });
    },
  },
  mounted() {
    this.testlogin();
  },
};
</script>

