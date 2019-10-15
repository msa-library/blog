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
            sm="10"
            md="10"
          >
            <v-card class="elevation-12">
              <v-toolbar
                color="primary"
                dark
                flat
              >
                <v-toolbar-title>Категория</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-form ref="formCategory">
                  <v-text-field
                    name="name"
                    prepend-icon="title"
                    type="text"
                    v-model="name"
                    :rules="nameRules"
                    label="Название"
                    required
                  ></v-text-field>
                </v-form>
                <v-alert type="error" outlined v-model="submitHasError">
                  {{this.submitError}}
                </v-alert>
              </v-card-text>
              <v-card-actions>
                <div class="flex-grow-1"></div>
                <v-btn color="primary" text to="/">Отмена</v-btn>
                <v-btn color="primary" @click="submit">Сохранить</v-btn>
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
    data: () => ({
      valid: true,
      submitError:null,
      name: '',
      nameRules: [
        v => !!v || 'Название обязательное поле',
      ],
      
      
    }),
    computed:{
      submitHasError:function(){
        return this.submitError?true:false;
      }
    },
    methods: {
      submit () {
        if (this.$refs.formCategory.validate()) {
          var data={
            Name:this.name,
          }
          this.$store.dispatch('createCategory',data).then(this.submitSuccess,this.submitFail);
        }
      },
      submitSuccess(response){
        console.log(response)
        this.$router.push('/');
      },
      submitFail(err){
        this.submitError=err.body.message;
      }
      
    },
    
  }
</script>