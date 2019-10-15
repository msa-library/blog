import Vue from 'vue';
import Vuetify from 'vuetify/lib';

//Vue.use(Vuetify);
Vue.use(Vuetify, {
  theme: {
      primary: '#1a76d1',
      secondary: '#ff1760',
      accent: '#805441'
    }
})

export default new Vuetify({
  icons: {
    iconfont: 'mdi',
  },
});


