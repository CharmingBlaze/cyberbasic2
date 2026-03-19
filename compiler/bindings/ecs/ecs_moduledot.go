package ecs

// ecsV2 maps lowercase method names to existing ECS / flat foreigns for the ecs.* v2 module.
var ecsV2 = map[string]string{
	"createworld":         "ECS.CreateWorld",
	"destroyworld":        "ECS.DestroyWorld",
	"createentity":        "ECS.CreateEntity",
	"destroyentity":       "ECS.DestroyEntity",
	"addcomponent":        "ECS.AddComponent",
	"hascomponent":        "ECS.HasComponent",
	"removecomponent":     "ECS.RemoveComponent",
	"gettransformx":       "ECS.GetTransformX",
	"gettransformy":       "ECS.GetTransformY",
	"gettransformz":       "ECS.GetTransformZ",
	"settransform":        "ECS.SetTransform",
	"placeentity":         "ECS.PlaceEntity",
	"getworldpositionx":   "ECS.GetWorldPositionX",
	"getworldpositiony":   "ECS.GetWorldPositionY",
	"getworldpositionz":   "ECS.GetWorldPositionZ",
	"gethealthcurrent":    "ECS.GetHealthCurrent",
	"gethealthmax":        "ECS.GetHealthMax",
	"querycount":          "ECS.QueryCount",
	"queryentity":         "ECS.QueryEntity",
	"createentitydefault": "CreateEntity",
	"addcomponentdefault": "AddComponent",
	"getcomponent":        "GetComponent",
	"removecomponentflat": "RemoveComponent",
	"runsystem":           "RunSystem",
}
