// Package gamedot exposes global "game" as a modfacade over game.RegisterGame foreigns.
package gamedot

import (
	"strings"

	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/vm"
)

// gameNames: all RegisterForeign names from compiler/bindings/game/game.go (kept in sync manually).
var gameNames = []string{
	"CreateParticleSystem", "EmitParticles", "SetParticleColor", "SetParticleLifetime", "SetParticleVelocity",
	"DrawParticles", "ParticleSetLayer",
	"AISetPosition", "GetAIPosition", "AIUpdate", "AIMoveTo", "AISetSpeed", "AIWander", "AIChase", "AIFlee",
	"AIAction", "AICondition", "AIRun", "AISequence", "AISelector",
	"AIBehaviorTreeCreate", "AIBehaviorTreeSetRoot",
	"AnimateValue", "AnimateColor", "AnimatePosition", "AnimateRotation",
	"CoroutineStart", "CoroutineYield", "CoroutineWait", "CoroutineStop",
	"TilemapCreate", "TilemapLoadByPath", "LoadTilemap", "TilemapLoad", "TilemapSave", "TilemapFill",
	"TilemapSetTile", "TilemapGetTile", "DrawTilemap", "TilemapDrawByMapId",
	"SetTileByMapId", "GetTileByMapId", "SetTile", "GetTile", "TilemapCollision",
	"TilemapSetTileset", "TilemapSetTilesetByMapId", "TilemapDeleteByMapId", "TilemapSetLayer", "TilemapSetParallax",
	"WeatherSetType", "WeatherSetIntensity", "WeatherSetWindDirection", "WeatherSetWindSpeed",
	"WeatherSetFogDensity", "WeatherSetLightningFrequency",
	"FireCreate", "FireSetSpreadRate", "FireSetSmokeEmitter", "FireSetLight", "FireSetActive",
	"SmokeSetDissolveRate", "SmokeSetRiseSpeed",
	"EnvironmentSetGlobalWind", "EnvironmentSetTemperature", "EnvironmentSetHumidity",
	"EnvironmentAffectParticles", "EnvironmentAffectWater", "EnvironmentAffectVegetation",
	"TimeSet", "TimeGet", "TimeSetSpeed",
	"SkyboxCreate", "SkyboxSetTexture", "SkyboxSetRotation", "SkyboxSetTint", "DrawSkybox",
	"CloudLayerCreate", "CloudLayerSetTexture", "CloudLayerSetHeight", "DrawCloudLayer",
	"DecalCreate", "DecalSetLifetime", "DecalRemove", "DrawDecals", "DrawFires",
	"PathfindGrid", "PathfindNavmesh", "FollowPath",
	"OnKeyPress", "OnMouseClick", "OnUpdate", "OnDraw", "OnCollision",
	"DebugDrawGrid", "DebugDrawBounds", "DebugLog", "DebugWatch",
	"Noise2D", "Noise3D", "GenerateDungeon", "GenerateTree", "GenerateCity",
	"DialogueLoad", "DialogueStart", "DialogueNext", "DialogueChoice",
	"DialogueShowText", "DialogueShowChoices", "DialogueSetVar", "DialogueGetVar",
	"InventoryCreate", "InventoryAddItem", "InventoryRemoveItem", "InventoryHasItem", "ItemDefine", "ItemSetProperty", "InventoryDraw",
	"CreateHingeJoint", "CreateBallJoint", "CreateSliderJoint", "CreateRagdoll", "RagdollEnable", "RagdollDisable",
	"ShaderGraphCreate", "ShaderGraphConnect", "ShaderNodeAdd", "ShaderNodeTexture", "ShaderNodeColor", "ShaderNodeMultiply", "ShaderNodeTime", "ShaderGraphCompile",
	"NetStartServer", "NetStartClient", "RPC", "ReplicateValue", "ReplicateVariable", "ReplicatePosition", "ReplicateRotation", "ReplicateScale",
	"AnimStateCreate", "AnimStateSetClip", "AnimTransition", "AnimSetParameter", "AnimSetState", "AnimUpdate",
}

func lowerMap(names []string) map[string]string {
	m := make(map[string]string, len(names))
	for _, n := range names {
		m[strings.ToLower(n)] = n
	}
	return m
}

// Register installs global "game" after game.RegisterGame.
func Register(v *vm.VM) {
	v.SetGlobal("game", modfacade.New(v, lowerMap(gameNames)))
}
