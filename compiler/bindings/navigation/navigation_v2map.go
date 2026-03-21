package navigation

// MethodToForeign maps lowercase navigation.* / ai.* method names to RegisterForeign names.
var MethodToForeign = map[string]string{
	"navgridcreate":           "NavGridCreate",
	"navgridsetwalkable":      "NavGridSetWalkable",
	"navgridsetcost":          "NavGridSetCost",
	"navgridfindpath":         "NavGridFindPath",
	"navmeshloadfromfile":     "NavMeshLoadFromFile",
	"navmeshcreatefromterrain": "NavMeshCreateFromTerrain",
	"navmeshaddobstacle":      "NavMeshAddObstacle",
	"navmeshremoveobstacle":   "NavMeshRemoveObstacle",
	"navmeshfindpathraw":      "NavMeshFindPathRaw",
	"navagentcreate":          "NavAgentCreate",
	"navagentsetspeed":        "NavAgentSetSpeed",
	"navagentsetradius":       "NavAgentSetRadius",
	"navagentsetdestination":  "NavAgentSetDestination",
	"navagentgetnextwaypoint": "NavAgentGetNextWaypoint",
	"navagentupdate":          "NavAgentUpdate",
	"navagentsetposition":     "NavAgentSetPosition",
	"navagentgetpositionx":    "NavAgentGetPositionX",
	"navagentgetpositiony":    "NavAgentGetPositionY",
	"navagentgetpositionz":    "NavAgentGetPositionZ",
}
