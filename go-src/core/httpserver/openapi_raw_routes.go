package httpserver

// RegAPIHandler 等裸 Handler 不会进入 GoFrame GetOpenApi(); 导出 App/Sofie 文档时在此补全 multipart 接口.

const (
	openAPISchemaLiveRoomDTO   = "xr-game-server.dto.liveroomdto."
	openAPISchemaUserInfoDTO   = "xr-game-server.dto.userinfodto."
	openAPISchemaShortVideoDTO = "xr-game-server.dto.shortvideodto."
)

// MergeOpenApiRawHandlerRoutes 向 OpenAPI 文档补全裸 Handler 路由及对应 components.schemas.
func MergeOpenApiRawHandlerRoutes(doc map[string]any) error {
	if doc == nil {
		return nil
	}
	paths, _ := doc["paths"].(map[string]any)
	if paths == nil {
		paths = map[string]any{}
		doc["paths"] = paths
	}
	schemas := ensureOpenAPIComponentSchemas(doc)
	mergeCreateLiveRoomOpenAPI(paths, schemas)
	mergeUploadAvatarOpenAPI(paths, schemas)
	mergeAppPublishShortVideoOpenAPI(paths, schemas)
	return nil
}

func ensureOpenAPIComponentSchemas(doc map[string]any) map[string]any {
	components, _ := doc["components"].(map[string]any)
	if components == nil {
		components = map[string]any{}
		doc["components"] = components
	}
	schemas, _ := components["schemas"].(map[string]any)
	if schemas == nil {
		schemas = map[string]any{}
		components["schemas"] = schemas
	}
	return schemas
}

func mergeCreateLiveRoomOpenAPI(paths, schemas map[string]any) {
	reqName := openAPISchemaLiveRoomDTO + "CreateLiveRoomReq"
	resName := openAPISchemaLiveRoomDTO + "CreateLiveRoomRes"
	schemas[reqName] = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title":              openAPIStringProp("直播间标题"),
			"cover":              openAPIBinaryProp("封面图片文件"),
			"notice":             openAPIStringProp("公告"),
			"category":           openAPIUintProp("分类(1=hot,2=game,5=语聊房,默认1)", "uint8"),
			"voiceChatMicMode":   openAPIUintProp("语聊上麦方式(1=自由,2=申请,3=房主单麦;category=5时有效,默认1)", "uint8"),
			"tagId":              openAPIUintProp("直播间标签ID", "uint64"),
			"gameCodes":          openAPIStringProp("推荐游戏编码(可多次传 gameCodes 或 gameCodes[];仅 category=2 有效)"),
			"ticket":             openAPINumberProp("直播间来源视频通话门票价格(钻石)", "float64", 0),
			"billing":            openAPINumberProp("视频通话价格(每分钟钻石)", "float64", nil),
			"privateInviteType":  openAPIUintProp("视频通话邀请类型(1=接受所有人,3=拒绝所有人,0或不传时 category=1 默认3 其他默认1)", "uint8"),
		},
	}
	schemas[resName] = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"roomId":  openAPIStringProp("直播间ID"),
			"guildId": openAPIStringProp("所属工会ID"),
		},
	}
	paths["/liveRoom/create"] = map[string]any{
		"post": openAPIMultipartPostOp(
			"创建直播间",
			[]string{"直播间"},
			"服务端 MultipartReader 流式解析(不经过 GoFrame Bind); 实际 Handler 为 handleCreateLiveRoom.",
			reqName,
			resName,
		),
	}
}

func mergeUploadAvatarOpenAPI(paths, schemas map[string]any) {
	reqName := openAPISchemaUserInfoDTO + "UploadAvatarReq"
	resName := openAPISchemaUserInfoDTO + "UploadAvatarRes"
	schemas[reqName] = map[string]any{
		"type": "object",
		"required": []any{"file"},
		"properties": map[string]any{
			"file": openAPIBinaryProp("头像图片文件"),
		},
	}
	schemas[resName] = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"avatar": openAPIStringProp("头像文件名"),
		},
	}
	paths["/userInfo/uploadAvatar"] = map[string]any{
		"post": openAPIMultipartPostOp(
			"上传用户头像",
			[]string{"用户基础信息"},
			"服务端 MultipartReader 流式解析; form 字段名须为 file.",
			reqName,
			resName,
		),
	}
}

func mergeAppPublishShortVideoOpenAPI(paths, schemas map[string]any) {
	reqName := openAPISchemaShortVideoDTO + "AppPublishShortVideoReq"
	resName := openAPISchemaShortVideoDTO + "AppPublishShortVideoRes"
	schemas[reqName] = map[string]any{
		"type": "object",
		"required": []any{"file"},
		"properties": map[string]any{
			"file":              openAPIBinaryProp("短视频文件(必填,仅支持 mp4)"),
			"previewFile":       openAPIBinaryProp("试看视频文件(可选,仅支持 mp4)"),
			"cover":             openAPIBinaryProp("封面图片(可选)"),
			"title":             openAPIStringProp("标题"),
			"isPaid":            openAPIUintProp("是否付费", "uint8"),
			"payDiamond":        openAPINumberProp("付费钻石", "float64", nil),
			"categoryId":        openAPIIntProp("分类ID"),
			"source":            openAPIUintProp("来源", "uint8"),
			"duration":          openAPIUintProp("时长(秒)", "uint32"),
			"freeWatchSeconds":  openAPIUintProp("免费试看秒数", "uint32"),
		},
	}
	if _, ok := schemas[resName]; !ok {
		schemas[resName] = map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":           openAPIStringProp("短视频 ID"),
				"video":        openAPIStringProp("视频完整URL"),
				"previewVideo": openAPIStringProp("试看视频完整URL,未上传时为空"),
				"cover":        openAPIStringProp("封面完整URL"),
			},
		}
	}
	paths["/shortVideo/appPublishShortVideo"] = map[string]any{
		"post": openAPIMultipartPostOp(
			"App流式上传并发布短视频",
			[]string{"短视频"},
			"服务端 MultipartReader 边读边写落盘; 不走 GoFrame Bind.",
			reqName,
			resName,
		),
	}
}

func openAPIMultipartPostOp(summary string, tags []string, description, reqSchema, resSchema string) map[string]any {
	tagSlice := make([]any, len(tags))
	for i, t := range tags {
		tagSlice[i] = t
	}
	op := map[string]any{
		"summary": summary,
		"tags":    tagSlice,
		"requestBody": map[string]any{
			"content": map[string]any{
				"multipart/form-data": map[string]any{
					"schema": map[string]any{
						"$ref":        "#/components/schemas/" + reqSchema,
						"description": "",
					},
				},
			},
		},
		"responses": map[string]any{
			"200": map[string]any{
				"description": "",
				"content": map[string]any{
					"application/json": map[string]any{
						"schema": map[string]any{
							"$ref":        "#/components/schemas/" + resSchema,
							"description": "",
						},
					},
				},
			},
		},
	}
	if description != "" {
		op["description"] = description
	}
	return op
}

func openAPIStringProp(description string) map[string]any {
	return map[string]any{
		"type":        "string",
		"format":      "string",
		"description": description,
	}
}

func openAPIBinaryProp(description string) map[string]any {
	return map[string]any{
		"type":        "string",
		"format":      "binary",
		"description": description,
	}
}

func openAPIUintProp(description, format string) map[string]any {
	return map[string]any{
		"type":        "integer",
		"format":      format,
		"description": description,
	}
}

func openAPIIntProp(description string) map[string]any {
	return map[string]any{
		"type":        "integer",
		"format":      "int",
		"description": description,
	}
}

func openAPINumberProp(description, format string, minimum any) map[string]any {
	prop := map[string]any{
		"type":        "number",
		"format":      format,
		"description": description,
	}
	if minimum != nil {
		prop["minimum"] = minimum
	}
	return prop
}
