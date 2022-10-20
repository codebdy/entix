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