<template>
    <v-app id="top">
      
      <v-container>
        <BreadcrumbsWidget v-bind:items="breadcrumbs"></BreadcrumbsWidget>
        <v-layout text-xs-center wrap>
          <v-flex xs12 md9 class="pa-2">
            <v-card>
                <v-img
                    :src="post.Src"
                    class="white--text"
                    height="200px"
                    gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
                >
                    <v-card-title
                        class="fill-height align-end"
                        v-text="post.Title"
                    ></v-card-title>
                </v-img>
                <v-card-text>
                    <span>{{getFullName()}}</span><br>
                    <span class="text--primary">
                        <span>{{post.SubTitle}}</span><br>
                        <span>{{post.Content}}</span>
                    </span>
                    <div>
                        Комментарии
                    </div>
                    <div>
                        <v-list>
                            <v-list-item
                                v-for="comment in comments"
                                :key="comment.Slug"   
                            >
                                <v-list-item-content>
                                    <v-list-item-title v-text="comment.Content"></v-list-item-title>
                                </v-list-item-content>
                            </v-list-item>
                        </v-list>
                    </div>
                    <div>
                        <v-form ref="formComment">
                            <v-text-field
                                name="comment"
                                type="text"
                                v-model="comment"
                                :rules="commentRules"
                                label="Комментарий"
                                required
                            ></v-text-field>
                        </v-form>
                        <v-alert type="error" outlined v-model="submitHasError">
                          {{this.submitError}}
                        </v-alert>
                    </div>
                    <div>
                        <v-btn color="primary" @click="submit">Отправить</v-btn>
                    </div>
                </v-card-text>

                
            </v-card>
          </v-flex>
          <v-flex xs12 md3>
            <!-- Автор -->
            <!-- <div class="pa-2">
              <AuthorWidget :item="author"></AuthorWidget>
            </div> -->

            <!-- Категории -->
            <div class="pa-2">
              <CaregoryListWidget v-bind:items="categories"></CaregoryListWidget>
            </div>

            <!-- Авторы -->
            <div class="pa-2">
              <AuthorListWidget v-bind:items="authors"></AuthorListWidget>
            </div>

          </v-flex>
        </v-layout>
      </v-container>
    </v-app>
    
</template>
<script>
    import BreadcrumbsWidget from '../widgets/Breadcrumbs'
    import AuthorWidget from '../widgets/Author'
    import CaregoryListWidget from '../widgets/CaregoryList'
    import AuthorListWidget from '../widgets/AuthorList'
    export default {
        components: {
            BreadcrumbsWidget,
            AuthorWidget,
            CaregoryListWidget,
            AuthorListWidget,
        },
        data: function () {
            return {
                breadcrumbs:[],
                submitError:null,
                
                comment:'',
                commentRules: [
                    v => !!v || 'Комментарий обязательное поле',
                ],

                post:{
                    PostId:'',
                    Title:'',
                    SubTitle:'',
                    Content:'',
                    Src:'',
                },
                author:{
                    SrcCover:'',
                    FirstName:'',
                    LastName:'',
                },
                authors:[],
                categories:[],
                comments:[],
            }
        },
        computed:{
            submitHasError:function(){
                return this.submitError?true:false;
            }
        },
        methods:{
            getFullName:function(){
                return this.author.FirstName+' '+this.author.LastName;
            },
            submit () {
                
                
                if (this.$refs.formComment.validate()) {
                    var data={
                        Content:this.comment,
                        PostId:this.post.Slug,
                    }
                    var self=this;
                    this.$store.dispatch('createComment',data).then(response=>{
                        if(response.body.Comment){
                            self.comments.push(response.body.Comment);
                        }
                        this.$refs.formComment.reset()
                    },err=>{
                        
                    });
                }
            },
        },
        created:function(){

            //Заполнение breadcrumbs
            this.breadcrumbs=[
                {
                    text: 'Главная',
                    disabled: false,
                    to: '/',
                },
                {
                    text: '',
                    disabled: true,
                    to: '/',
                },
            ]
            
            //Запрос данных статьи с сервера
            var data={
                slug:this.$route.params.slug
            }

            var self=this
            this.$store.dispatch('getPost',data).then(response=>{
                if(response.body.Post){
                    
                    self.post=response.body.Post;
                    self.breadcrumbs[1].text=self.post.Title;
                    self.author=response.body.Post.Author;
                    self.categories=response.body.Post.PostCategories;
                    if(response.body.Post.PostComments){
                        self.comments=response.body.Post.PostComments;
                    }
                    
                }
                //console.log(response.body.Post)
            },err=>{
                console.log(err)
            });
            

            //Загрузка авторов
            this.$store.dispatch('getAuthors',{}).then(response=>{
                
                if(response.body.Authors){
                self.authors=response.body.Authors
                
                }
            },err=>{
                console.log("Ошибка загрузки авторов: ",err)
            });
        }
    }
</script>