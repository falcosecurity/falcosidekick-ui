<template>
  <div class="login">
    <form class="login-card" @submit.prevent="authenticate">
      <div class="login-brand">
        <img
          src="https://github.com/falcosecurity/falcosidekick/raw/master/imgs/falcosidekick_color.png"
          alt="">
        <div>
          <div class="login-title">Falcosidekick UI</div>
          <div class="login-sub">Sign in to view Falco events</div>
        </div>
      </div>
      <label>
        <span>Username</span>
        <input v-model="username" type="text" autocomplete="username" autofocus>
      </label>
      <label>
        <span>Password</span>
        <input v-model="password" type="password" autocomplete="current-password">
      </label>
      <div v-if="failedAuth === true" class="login-error">
        <v-icon small>mdi-alert-circle-outline</v-icon> {{ failMsg }}
      </div>
      <button type="submit" class="login-btn">Sign in</button>
    </form>
  </div>
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


<style scoped>
.login {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--sidebar-bg);
  padding: 16px;
}

.login-card {
  width: 100%;
  max-width: 360px;
  background: var(--sidebar-bg-2);
  border: 1px solid var(--sidebar-border);
  border-radius: 10px;
  padding: 28px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
}

.login-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
}

.login-brand img {
  width: 36px;
  height: 36px;
}

.login-title {
  color: var(--sidebar-text-strong);
  font-size: 16px;
  font-weight: 600;
}

.login-sub {
  color: var(--sidebar-text);
  font-size: 12px;
}

label {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

label span {
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--sidebar-heading);
}

input {
  background: var(--sidebar-bg);
  border: 1px solid var(--sidebar-border);
  border-radius: 6px;
  padding: 9px 12px;
  color: var(--sidebar-text-strong);
  font-size: 14px;
  outline: none;
}

input:focus {
  border-color: var(--accent);
}

.login-error {
  font-size: 12px;
  color: #ff8a8a;
}

.login-error .v-icon {
  color: inherit !important;
}

.login-btn {
  margin-top: 4px;
  height: 38px;
  border-radius: 6px;
  background: var(--accent);
  color: #fff;
  font-weight: 600;
  font-size: 14px;
}

.login-btn:hover {
  filter: brightness(1.08);
}
</style>
