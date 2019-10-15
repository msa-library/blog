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
                <v-toolbar-title>Регистрация</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-form ref="formSignUp">
                  <v-text-field
                    name="firstName"
                    prepend-icon="person"
                    type="text"
                    v-model="firstName"
                    :rules="firstNameRules"
                    label="Имя"
                    required
                  ></v-text-field>

                  <v-text-field
                    name="lastName"
                    prepend-icon="person"
                    type="text"
                    v-model="lastName"
                    :rules="lastNameRules"
                    label="Фамилия"
                    required
                  ></v-text-field>

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

                  <v-checkbox
                    v-model="agree"
                    :rules="agreeRules"
                    label="Do you agree?"
                    required
                    ></v-checkbox>
                </v-form>
                <v-alert type="error" outlined v-model="submitHasError">
                  {{this.submitError}}
                </v-alert>
              </v-card-text>
              <v-card-actions>
                <div class="flex-grow-1"></div>
                <v-btn color="primary" text to="/signin">Вход</v-btn>
                <v-btn color="primary" @click="submit">Регистрация</v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  export default {
    props: {
      //source: String,
    },
    data: () => ({
      valid: true,

      submitError:null,

      firstName: 'Ali',
      firstNameRules: [
        v => !!v || 'Имя обязательное поле',
      ],

      lastName: 'Connors',
      lastNameRules: [
        v => !!v || 'Фамилия обязательное поле',
      ],
      
      email: 'ali@connors.com',
      emailRules: [
        v => !!v || 'E-mail обязательное поле',
        v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
      ],

      password: '12345678',
      passwordRules: [
        v => !!v || 'Пароль обязательное поле',
      ],

      agree:true,
      agreeRules: [
        v => !!v || 'You must agree to continue!'
      ],
      
    }),
    computed:{
      submitHasError:function(){
        return this.submitError?true:false;
      }
    },
    methods: {
      submit () {
        if (this.$refs.formSignUp.validate()) {
          //this.snackbar = true
          var data={
            Username:this.email,
            Password:this.password,
            FirstName:this.firstName,
            LastName:this.lastName,
          }
          this.$store.dispatch('signup',data).then(this.submitSuccess,this.submitFail);
        }
      },
      submitSuccess(response){
        
        this.$cookie.set('Authorization', response.headers.get('Authorization'), 1);
        this.$router.push('/');
        
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