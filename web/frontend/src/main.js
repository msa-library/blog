import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';

//VueCookie
import VueCookie from 'vue-cookie';
Vue.use(VueCookie);

//VueRouter
import VueRouter from 'vue-router';
import routes from './routes.js';
Vue.use(VueRouter);
var router = new VueRouter({
  routes:routes,
  mode: 'history',
})

//Vuex & VueResource
import store from './store'


Vue.config.productionTip = false

new Vue({
  vuetify,
  store:store,
  router:router,
  render: h => h(App)
}).$mount('#app')
