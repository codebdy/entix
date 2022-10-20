package basic

	//if a.appUuid == consts.SYSTEM_APP_UUID {
		mutationFields[consts.PUBLISH] = &graphql.Field{
			Type: a.modelParser.OutputType(metaEntity.Name()),
			Args: graphql.FieldConfigArgument{
				consts.APPID: &graphql.ArgumentConfig{
					Type: &graphql.NonNull{OfType: graphql.ID},
				},
			},
			Resolve: resolve.PublishMetaResolve,
		}


		func doPublish(r *repository.Repository, appUuid string) error {
			publishedMeta := r.QueryPublishedMeta(appUuid)
			nextMeta := r.QueryNextMeta(appUuid)
			appId := r.QueryAppId(appUuid)
			fmt.Println("Start to publish")
			// fmt.Println("Published Meta ID:", publishedMeta.(utils.Object)["id"])
			// fmt.Println("Next Meta ID:", nextMeta.(utils.Object)["id"])
		
			if nextMeta == nil {
				panic("Can not find unpublished meta")
			}
			publishedModel := model.New(r.Model.AppUuid, r.MergeModel(appUuid, repository.DecodeContent(publishedMeta, appId)))
			nextModel := model.New(r.Model.AppUuid, r.MergeModel(appUuid, repository.DecodeContent(nextMeta, appId)))
			nextModel.Graph.Validate()
			diff := model.CreateDiff(publishedModel, nextModel)
			r.ExcuteDiff(diff)
			fmt.Println("ExcuteDiff success")
			metaObj := nextMeta.(utils.Object)
			metaObj[consts.META_STATUS] = meta.META_STATUS_PUBLISHED
			metaObj[consts.META_PUBLISHEDAT] = time.Now()
			_, err := r.SaveOne(data.NewInstance(metaObj, r.Model.Graph.GetMetaEntity()))
			if err != nil {
				return err
			}
		
			return nil
		}
		
		func PublishMetaResolve(p graphql.ResolveParams) (interface{}, error) {
			defer utils.PrintErrorStack()
			repos := repository.New(model)
			repos.MakeSupperVerifier()
			appUuid := p.Args[consts.APPUUID]
			if appUuid == nil {
				appUuid = consts.SYSTEM_APP_UUID
			}
			err := doPublish(repos, appUuid.(string))
			if err != nil {
				logs.WriteBusinessLog(model, p, logs.PUBLISH_META, logs.FAILURE, err.Error())
			} else {
				logs.WriteBusinessLog(model, p, logs.PUBLISH_META, logs.SUCCESS, "")
			}
			return "success", err
		}
		
		func SyncMetaResolve(p graphql.ResolveParams, model *model.Model) (interface{}, error) {
			object := p.Args[consts.ARG_OBJECT].(map[string]interface{})
			repos := repository.New(model)
			repos.MakeEntityAbilityVerifier(p, meta.META_ENTITY_UUID)
			return repos.InsertOne(data.NewInstance(object, model.Graph.GetMetaEntity()))
		}
		