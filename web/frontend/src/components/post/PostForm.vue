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
                <v-toolbar-title>Пост</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-form ref="formPost">
                  <v-text-field
                    name="title"
                    prepend-icon="title"
                    type="text"
                    v-model="title"
                    :rules="titleRules"
                    label="Заголовок"
                    required
                  ></v-text-field>

                  <v-text-field
                    name="subtitle"
                    prepend-icon="text_fields"
                    type="text"
                    v-model="subtitle"
                    :rules="subtitleRules"
                    label="подзаголовок"
                    required
                  ></v-text-field>
                  <v-text-field
                    name="content"
                    prepend-icon="format_align_justify"
                    type="text"
                    v-model="content"
                    :rules="contentRules"
                    label="текст"
                    required
                  ></v-text-field>
                  <v-layout row wrap>
                    <v-flex v-bind:key="index" v-for="(item, index) in categories" xs4>
                      <v-checkbox
                        v-model="item.selected"
                        :label="item.Name"
                      ></v-checkbox>
                    </v-flex>
                  </v-layout>
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
      categories:[],
      submitError:null,
      
      
      title: 'title',
      titleRules: [
        v => !!v || 'Заголовок обязательное поле',
      ],

      subtitle: 'subtitle',
      subtitleRules: [
        v => !!v || 'Подзаголовок обязательное поле',
      ],

      content: 'content',
      contentRules: [
        v => !!v || 'Текст обязательное поле',
      ],
      
    }),
    computed:{
      submitHasError:function(){
        return this.submitError?true:false;
      }
    },
    methods: {
      submit () {
        var categories=[];
        for(var i=0; i<this.categories.length;i++){
          if(this.categories[i] && this.categories[i].selected){
            categories.push(this.categories[i].Slug)
          }
        }
        
        if (this.$refs.formPost.validate()) {
          var data={
            Title:this.title,
            SubTitle:this.subtitle,
            Content:this.content,
            Categories:categories.join(','),
          }
          this.$store.dispatch('createPost',data).then(this.submitSuccess,this.submitFail);
        }
      },
      submitSuccess(response){
        this.$router.push('/');
      },
      submitFail(err){
        this.submitError=err.body.message;
      }
      
    },
    created:function(){
      //Загрузка списка всех категорий
      var self=this;
      this.$store.dispatch('getCategories',{}).then(response=>{
        if(response.body.Categories){
          self.categories=response.body.Categories
        }
      },err=>{
        console.log("Ошибка загрузки категорий: ",err)
      });
    }
    
  }
</script>