<template>
  <v-app id="inspire">
    <v-content>
      <v-container
        class="fill-height"
        fluid
      >
        <v-row
          align="center"
          justify="center"
        >
          <v-col
            cols="12"
            sm="8"
            md="4"
          >
            <v-card class="elevation-12">
              <v-toolbar
                color="primary"
                dark
                flat
              >
                <v-toolbar-title>Вход</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-form ref="formformSignOn">
                  <v-text-field
                    name="login"
                    prepend-icon="person"
                    type="text"
                    v-model="email"
                    :rules="emailRules"
                    label="E-mail"
                    required
                  ></v-text-field>

                  <v-text-field
                    id="password"
                    label="Пароль"
                    name="password"
                    v-model="password"
                    prepend-icon="lock"
                    type="password"
                    :rules="passwordRules"
                    required
                  ></v-text-field>
                </v-form>
                <v-alert type="error" outlined v-model="submitHasError">
                  {{this.submitError}}
                </v-alert>
              </v-card-text>
              <v-card-actions>
                <div class="flex-grow-1"></div>
                <v-btn color="primary" text to="/signup">Регистрация</v-btn>
                <v-btn color="primary" @click="submit">Вход</v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import { constants } from 'crypto';
import { setTimeout } from 'timers';
  export default {
    props: {
      source: String,
    },
    data: () => ({
      valid: true,

      submitError:null,
      
      
      email: 'ptimofeev@yandex.ru',
      emailRules: [
        v => !!v || 'E-mail обязательное поле',
        v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
      ],

      password: '12345678',
      passwordRules: [
        v => !!v || 'Пароль обязательное поле',
      ],
      
    }),
    computed:{
      submitHasError:function(){
        return this.submitError?true:false;
      }
    },
    methods: {
      submit () {
        if (this.$refs.formformSignOn.validate()) {
          //this.snackbar = true
          var data={
            Username:this.email,
            Password:this.password
          }
          this.$store.dispatch('signin',data).then(this.submitSuccess,this.submitFail);
        }
      },
      submitSuccess(response){

        this.$cookie.set('Authorization', response.headers.get('Authorization'), 1);
        if(this.$route.params.next_to){
          this.$router.push(this.$route.params.next_to);
        }else{
          this.$router.push('/');
        }
      },
      submitFail(err){
        this.submitError=err.body.message;
      }
    },
    mounted:function(){
      //var auth=this.$cookie.get('authorization');
      //if(auth){
      //  this.$router.push('/');
      //}
    }
  }
</script>