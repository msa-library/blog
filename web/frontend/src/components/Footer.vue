<template>
    <div class="footer">
        <v-layout :row="screenWidth > breakPoint"
                  :wrap="screenWidth > breakPoint"
                  :reverse="screenWidth <= breakPoint"
                  :column="screenWidth <= breakPoint"
                  class="footer-inner"
                  align-center
                  :justify-between="screenWidth > breakPoint"
        >
            <v-flex md4 sm12>
                <div>{{year}} © MSA Blog • Все права защищены</div>
            </v-flex>
        </v-layout>
    </div>
</template>

<script>
  export default {
    name: "Footer",
    props: {
      screenWidth: {
        type: Number,
      },
      breakPoint: {
        type: Number,
      }
    },
    data: function () {
      return {
        year: new Date().getFullYear(),
      }
    }
  }
</script>

<style scoped>
    .footer .footer-inner {
        max-width: 1198px;
        margin: 0 auto;
        width: 100%;
        padding: 22px 24px;
        height: 100%;
        position: relative;
        color: #000000;
        font-size: 14px;
    }
    .footer .footer-inner .v-btn {
        font-size: 14px;
    }
    @media screen and (max-width: 959px) {
        .footer .footer-inner {
            padding: 15px 5px 29px;
        }
    }
    @media screen and (max-width: 599px) {
        .footer .footer-inner .v-btn {
            padding: 0 5px;
            min-width: 0;
        }
    }
    @media screen and (max-width: 370px) {
        .footer .footer-inner,
        .footer .footer-inner .v-btn {
            font-size: 12px;
        }
    }
</style>
