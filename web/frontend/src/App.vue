<template>
  <v-app>
    <Navbar></Navbar>
    <v-content>
      <router-view>
      </router-view>
    </v-content>
    <Footer :screen-width="screenWidth" :break-point="breakPoint"></Footer>
  </v-app>
</template>

<script>
  import Navbar from './components/Navbar.vue'
  import Footer from './components/Footer.vue'
  export default {
    name: 'App',
    components: {
      Navbar,
      Footer,
    },
    data: function () {
      return {
        screenWidth: 0,
        breakPoint: 959,
      }
    },
    methods: {
      onScreenResize: function () {
        this.screenWidth = window.innerWidth;
      },
    },
    mounted: function () {

      this.$store.state.authorization=this.$cookie.get('Authorization');
      
      this.$nextTick(function () {
        window.addEventListener("resize", this.onScreenResize);
        this.onScreenResize();
      });
    }
  };
</script>

<style>
  /* fallback */
  @font-face {
    font-family: 'Material Icons';
    font-style: normal;
    font-weight: 400;
    src: url('./assets/fonts/MaterialIcons.woff2') format('woff2');
  }

  .material-icons {
    font-family: 'Material Icons';
    font-weight: normal;
    font-style: normal;
    font-size: 24px;
    line-height: 1;
    letter-spacing: normal;
    text-transform: none;
    display: inline-block;
    white-space: nowrap;
    word-wrap: normal;
    direction: ltr;
    -webkit-font-feature-settings: 'liga';
    -webkit-font-smoothing: antialiased;
  }

  .application {
    font-family: 'Open Sans', sans-serif;
    font-weight: 400;
    background: #fff;
  }

  /*
  .theme--light.application {
    background: #fff;
  }

  .terms ul{
      list-style-type:none
  }
  */
</style>

